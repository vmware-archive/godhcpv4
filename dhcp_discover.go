package dhcpv4

// DHCPDiscover is a client broadcast packet to locate available servers.
type DHCPDiscover struct {
	Packet
	ReplyWriter
}

func (req DHCPDiscover) CreateDHCPOffer() DHCPOffer {
	rep := DHCPOffer{
		Packet: NewReply(req.Packet),
		req:    req.Packet,
	}

	rep.SetMessageType(MessageTypeDHCPOffer)
	return rep
}
