package dhcpv4

// DHCPInform is a client to server packet, asking only for local configuration
// parameters; client already has externally configured network address.
type DHCPInform struct {
	Packet
	ReplyWriter
}

func (req DHCPInform) CreateDHCPAck() DHCPAck {
	rep := DHCPAck{
		Packet: NewReply(req.Packet),
		req:    req.Packet,
	}

	rep.SetMessageType(MessageTypeDHCPAck)
	return rep
}
