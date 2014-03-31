package dhcpv4

type DHCPOffer struct {
	Packet
}

// From RFC2131, table 3:
//   Option                    DHCPOFFER
//   ------                    ---------
//   Requested IP address      MUST NOT
//   IP address lease time     MUST
//   Use 'file'/'sname' fields MAY
//   DHCP message type         DHCPOFFER
//   Parameter request list    MUST NOT
//   Message                   SHOULD
//   Client identifier         MUST NOT
//   Vendor class identifier   MAY
//   Server identifier         MUST
//   Maximum message size      MUST NOT
//   All others                MAY

var dhcpOfferValidation = []Validation{
	ValidateMustNot(OptionAddressRequest),
	ValidateMust(OptionAddressTime),
	ValidateMustNot(OptionParameterList),
	ValidateMustNot(OptionClientID),
	ValidateMust(OptionDHCPServerID),
	ValidateMustNot(OptionDHCPMaxMsgSize),
}

func (d *DHCPOffer) Validate() error {
	return Validate(d.Packet, dhcpOfferValidation)
}
