package dhcpv4

// DHCPRequest is a client message to servers either (a) requesting offered
// parameters from one server and implicitly declining offers from all others,
// (b) confirming correctness of previously allocated address after, e.g.,
// system reboot, or (c) extending the lease on a particular network address.
type DHCPRequest struct {
	Packet
	ReplyWriter
}

func (req DHCPRequest) CreateDHCPAck() DHCPAck {
	rep := DHCPAck{
		Packet: NewReply(req.Packet),
		req:    req.Packet,
	}

	rep.SetMessageType(MessageTypeDHCPAck)
	return rep
}

func (req DHCPRequest) CreateDHCPNak() DHCPNak {
	rep := DHCPNak{
		Packet: NewReply(req.Packet),
		req:    req.Packet,
	}

	rep.SetMessageType(MessageTypeDHCPNak)
	return rep
}
