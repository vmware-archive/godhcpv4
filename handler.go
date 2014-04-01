package dhcpv4

import "net"

// PacketReader defines the ReadFrom function as defined in net.PacketConn.
type PacketReader interface {
	ReadFrom(b []byte) (n int, addr net.Addr, err error)
}

// PacketWriter defines the ReadFrom function as defined in net.PacketConn.
type PacketWriter interface {
	WriteTo(b []byte, addr net.Addr) (n int, err error)
}

// PacketConn groups PacketReader and PacketWriter to form a subset of net.PacketConn.
type PacketConn interface {
	PacketReader
	PacketWriter
}

type replyWriter struct {
	pw PacketWriter

	// The client address, if any
	addr net.UDPAddr
}

func (rw *replyWriter) WriteReply(r Reply) error {
	var err error

	err = r.Validate()
	if err != nil {
		return err
	}

	bytes, err := r.ToBytes()
	if err != nil {
		return err
	}

	req := r.Request()
	addr := rw.addr
	bcast := req.Flags()[0] & 128

	// Broadcast the reply if the request packet has no address associated with
	// it, or if the client explicitly asks for a broadcast reply.
	if addr.IP.Equal(net.IPv4zero) || bcast > 0 {
		addr.IP = net.IPv4bcast
	}

	_, err = rw.pw.WriteTo(bytes, &addr)
	if err != nil {
		return err
	}

	return nil
}

// Handler defines the interface an object needs to implement to handle DHCP
// packets. The handler should do a type switch on the Request object that is
// passed as argument to determine what kind of packet it is dealing with. It
// can use the WriteReply function on the request to send a reply back to the
// peer responsible for sending the request packet. While the handler may be
// blocking, it is not encouraged. Rather, the handler should return as soon as
// possible to avoid blocking the serve loop. If blocking operations need to be
// executed to determine if the request packet needs a reply, and if so, what
// kind of reply, it is recommended to handle this in separate goroutines. The
// WriteReply function can be called from multiple goroutines without needing
// extra synchronization.
type Handler interface {
	ServeDHCP(req Request)
}

// Serve reads packets off the network and passes them to the specified
// handler. It is up to the handler to packets to per-client serve loops, if
// that is what you want.
func Serve(pc PacketConn, h Handler) error {
	buf := make([]byte, 65536)

	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			return err
		}

		p, err := PacketFromBytes(buf[:n])
		if err != nil {
			continue
		}

		// Filter everything but requests
		if OpCode(p.Op()[0]) != BootRequest {
			continue
		}

		rw := replyWriter{
			pw:   pc,
			addr: *addr.(*net.UDPAddr),
		}

		var req Request

		switch p.GetMessageType() {
		case MessageTypeDHCPDiscover:
			req = DHCPDiscover{p, &rw}
		case MessageTypeDHCPRequest:
			req = DHCPRequest{p, &rw}
		case MessageTypeDHCPDecline:
			req = DHCPDecline{p}
		case MessageTypeDHCPRelease:
			req = DHCPRelease{p}
		case MessageTypeDHCPInform:
			req = DHCPInform{p, &rw}
		}

		if req != nil {
			h.ServeDHCP(req)
		}
	}
}
