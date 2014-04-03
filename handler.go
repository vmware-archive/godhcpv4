package dhcpv4

import (
	"fmt"
	"net"

	"code.google.com/p/go.net/ipv4"
)

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

// Serve reads packets off the network and calls the specified handler.
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

// packetConnFilter wraps net.PacketConn and only reads and writes packet from
// and to the specified network interface.
type packetConnFilter struct {
	net.PacketConn

	ipv4pc *ipv4.PacketConn
	ipv4cm *ipv4.ControlMessage
}

// ReadFrom reads a packet from the connection copying the payload into b. It
// inherits its semantics from ipv4.PacketConn and subsequently net.PacketConn,
// but filters out packets that arrived on an interface other than the one
// specified in the packetConnFilter structure.
func (p *packetConnFilter) ReadFrom(b []byte) (int, net.Addr, error) {
	for {
		n, cm, src, err := p.ipv4pc.ReadFrom(b)
		if err != nil {
			return n, src, err
		}

		// Read another packet if it didn't arrive on the right interface
		if cm.IfIndex != p.ipv4cm.IfIndex {
			continue
		}

		return n, src, err
	}
}

// WriteTo writes a packet with payload b to addr. It inherits its semantics
// from ipv4.PacketConn and subsequently net.PacketConn, but explicitly sends
// the packet over the interface specified in the packetConnFilter structure.
func (p *packetConnFilter) WriteTo(b []byte, addr net.Addr) (int, error) {
	return p.ipv4pc.WriteTo(b, p.ipv4cm, addr)
}

// PacketConnFilter wraps a net.PacketConn and only reads packets from and
// writes packets to the network interface associated with the specified IP
// address. It may return an error if it cannot initialize the underlying
// socket correctly. It panics if it cannot find the network interface
// associated with the specified IP.
func PacketConnFilter(pc net.PacketConn, ip net.IP) (net.PacketConn, error) {
	ipv4pc := ipv4.NewPacketConn(pc)
	if err := ipv4pc.SetControlMessage(ipv4.FlagInterface, true); err != nil {
		return nil, err
	}

	p := packetConnFilter{
		PacketConn: pc,

		ipv4pc: ipv4pc,
		ipv4cm: &ipv4.ControlMessage{
			IfIndex: LookupInterfaceIndexForIP(ip),
		},
	}

	return &p, nil
}

// LookupInterfaceIndexForIP finds the system-wide network interface index that
// is associated with the specified IP address.
func LookupInterfaceIndexForIP(ip net.IP) int {
	is, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range is {
		as, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		for _, a := range as {
			if a.(*net.IPNet).IP.String() == ip.String() {
				return i.Index
			}
		}
	}

	// Not really a recoverable error...
	panic(fmt.Sprintf("dhcpv4: can't find network interface for: %s", ip))
}
