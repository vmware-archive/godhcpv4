package dhcpv4

// Request
type Request interface {
	PacketGetter
	OptionGetter
}

// Reply defines an interface implemented by DHCP replies.
type Reply interface {
	Validate() error
	ToBytes() ([]byte, error)
	Request() Packet

	PacketSetter
	OptionSetter
}

// ReplyWriter defines an interface for the object that writes a reply to the
// network to the intended received, be it via broadcast or unicast.
type ReplyWriter interface {
	WriteReply(r Reply) error
}
