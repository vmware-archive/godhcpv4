package dhcpv4

type DHCPNak struct {
	Packet
}

// From RFC2131, table 3:
//   Option                    DHCPNAK
//   ------                    -------
//   Requested IP address      MUST NOT
//   IP address lease time     MUST NOT
//   Use 'file'/'sname' fields MUST NOT
//   DHCP message type         DHCPNAK
//   Parameter request list    MUST NOT
//   Message                   SHOULD
//   Client identifier         MAY
//   Vendor class identifier   MAY
//   Server identifier         MUST
//   Maximum message size      MUST NOT
//   All others                MUST NOT

var dhcpNakAllowedOptions = []Option{
	OptionDHCPMsgType,
	OptionDHCPMessage,
	OptionClientID,
	OptionClassID,
	OptionDHCPServerID,
}

var dhcpNakValidation = []Validation{
	ValidateMust(OptionDHCPServerID),
	ValidateAllowedOptions(dhcpNakAllowedOptions),
}

func (d *DHCPNak) Validate() error {
	return Validate(d.Packet, dhcpNakValidation)
}
