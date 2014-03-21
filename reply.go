package dhcpv4

import "errors"

type Reply Packet

var (
	ErrNoReply = errors.New("dhcpv4: not a reply packet")
)

// NewReplyFromBytes converts a byte slice into a reply.
func NewReplyFromBytes(b []byte) (*Reply, error) {
	p, err := PacketFromBytes(b)
	if err != nil {
		return nil, err
	}

	if OpCode(p.Op()[0]) != BootReply {
		return nil, ErrNoReply
	}

	rep := Reply(p)
	return &rep, nil
}
