package dhcpv4

import "errors"

type Request Packet

var (
	ErrNoRequest = errors.New("dhcpv4: not a request packet")
)

// NewRequestFromBytes converts a byte slice into a request.
func NewRequestFromBytes(b []byte) (*Request, error) {
	p, err := PacketFromBytes(b)
	if err != nil {
		return nil, err
	}

	if OpCode(p.Op()[0]) != BootRequest {
		return nil, ErrNoRequest
	}

	req := Request(p)
	return &req, nil
}

// CreateReply creates a reply, copying relevant fields from the reply into it.
// These fields include the hardware address type and length, the flags, the
// client hardware address, and the relay agent IP address.
func (req *Request) CreateReply() *Reply {
	var rep = Reply(NewPacket(BootReply))

	// Hardware type and address length
	rep.HType()[0] = 1 // Ethernet
	rep.HLen()[0] = 6  // MAC-48 is 6 octets

	// Copy transaction identifier
	copy(rep.XId(), req.XId())

	// Copy fields from request (per RFC2131, section 4.3, table 3)
	copy(rep.Flags(), req.Flags())
	copy(rep.CHAddr(), req.CHAddr())
	copy(rep.GIAddr(), req.GIAddr())

	// The remainder of the fields are set depending on the outcome of the
	// handler. Once the packet has been filled in, it should be validated before
	// sending it out on the wire.
	return &rep
}
