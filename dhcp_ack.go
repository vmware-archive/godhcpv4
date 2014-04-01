package dhcpv4

// DHCPAck is a server to client packet with configuration parameters,
// including committed network address.
type DHCPAck struct {
	Packet

	req Packet
}

// From RFC2131, table 3:
//   Option                    DHCPACK
//   ------                    -------
//   Requested IP address      MUST NOT
//   IP address lease time     MUST (DHCPREQUEST)
//                             MUST NOT (DHCPINFORM)
//   Use 'file'/'sname' fields MAY
//   DHCP message type         DHCPACK
//   Parameter request list    MUST NOT
//   Message                   SHOULD
//   Client identifier         MUST NOT
//   Vendor class identifier   MAY
//   Server identifier         MUST
//   Maximum message size      MUST NOT
//   All others                MAY

var dhcpAckOnRequestValidation = []Validation{
	ValidateMust(OptionAddressTime),
}

var dhcpAckOnInformValidation = []Validation{
	ValidateMustNot(OptionAddressTime),
}

var dhcpAckValidation = []Validation{
	ValidateMustNot(OptionAddressRequest),
	ValidateMustNot(OptionParameterList),
	ValidateMustNot(OptionClientID),
	ValidateMust(OptionDHCPServerID),
	ValidateMustNot(OptionDHCPMaxMsgSize),
}

func (d DHCPAck) Validate() error {
	var err error

	// Validation is subtly different based on type of request
	switch d.req.GetMessageType() {
	case MessageTypeDHCPRequest:
		err = Validate(d.Packet, dhcpAckOnRequestValidation)
	case MessageTypeDHCPInform:
		err = Validate(d.Packet, dhcpAckOnInformValidation)
	}

	if err != nil {
		return err
	}

	return Validate(d.Packet, dhcpAckValidation)
}

func (d DHCPAck) ToBytes() ([]byte, error) {
	return PacketToBytes(d.Packet)
}

func (d DHCPAck) Request() Packet {
	return d.req
}
