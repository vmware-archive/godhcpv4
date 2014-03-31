package dhcpv4

// MessageType is the type for the various DHCP messages defined in RFC2132.
type MessageType byte

const (
	MessageTypeDHCPDiscover = MessageType(1)
	MessageTypeDHCPOffer    = MessageType(2)
	MessageTypeDHCPRequest  = MessageType(3)
	MessageTypeDHCPDecline  = MessageType(4)
	MessageTypeDHCPAck      = MessageType(5)
	MessageTypeDHCPNak      = MessageType(6)
	MessageTypeDHCPRelease  = MessageType(7)
	MessageTypeDHCPInform   = MessageType(8)
)

// Option is the type for DHCP option tags.
type Option byte

// OptionMap maps DHCP option tags to their values.
type OptionMap map[Option][]byte

// GetOption gets the []byte value of an option.
func (om OptionMap) GetOption(o Option) ([]byte, bool) {
	v, ok := om[o]
	return v, ok
}

// SetOption sets the []byte value of an option.
func (om OptionMap) SetOption(o Option, v []byte) {
	om[o] = v
	return
}

// GetMessageType gets the message type from the DHCPMsgType option field.
func (om OptionMap) GetMessageType() MessageType {
	v, ok := om.GetOption(OptionDHCPMsgType)
	if !ok || len(v) != 1 {
		return MessageType(0)
	}

	return MessageType(v[0])
}

// SetMessageType sets the message type in the DHCPMsgType option field.
func (om OptionMap) SetMessageType(m MessageType) {
	om.SetOption(OptionDHCPMsgType, []byte{byte(m)})
	return
}

// From RFC2132: DHCP Options and BOOTP Vendor Extensions
const (
	// RFC2132 Section 3: RFC 1497 Vendor Extensions
	OptionPad           = Option(0)
	OptionEnd           = Option(255)
	OptionSubnetMask    = Option(1)
	OptionTimeOffset    = Option(2)
	OptionRouter        = Option(3)
	OptionTimeServer    = Option(4)
	OptionNameServer    = Option(5)
	OptionDomainServer  = Option(6)
	OptionLogServer     = Option(7)
	OptionQuotesServer  = Option(8)
	OptionLPRServer     = Option(9)
	OptionImpressServer = Option(10)
	OptionRLPServer     = Option(11)
	OptionHostname      = Option(12)
	OptionBootFileSize  = Option(13)
	OptionMeritDumpFile = Option(14)
	OptionDomainName    = Option(15)
	OptionSwapServer    = Option(16)
	OptionRootPath      = Option(17)
	OptionExtensionFile = Option(18)

	// RFC2132 Section 4: IP Layer Parameters per Host
	OptionForwardOnOff  = Option(19)
	OptionSrcRteOnOff   = Option(20)
	OptionPolicyFilter  = Option(21)
	OptionMaxDGAssembly = Option(22)
	OptionDefaultIPTTL  = Option(23)
	OptionMTUTimeout    = Option(24)
	OptionMTUPlateau    = Option(25)

	// RFC2132 Section 5: IP Layer Parameters per Interface
	OptionMTUInterface     = Option(26)
	OptionMTUSubnet        = Option(27)
	OptionBroadcastAddress = Option(28)
	OptionMaskDiscovery    = Option(29)
	OptionMaskSupplier     = Option(30)
	OptionRouterDiscovery  = Option(31)
	OptionRouterRequest    = Option(32)
	OptionStaticRoute      = Option(33)

	// RFC2132 Section 6: Link Layer Parameters per Interface
	OptionTrailers   = Option(34)
	OptionARPTimeout = Option(35)
	OptionEthernet   = Option(36)

	// RFC2132 Section 7: TCP Parameters
	OptionDefaultTCPTTL = Option(37)
	OptionKeepaliveTime = Option(38)
	OptionKeepaliveData = Option(39)

	// RFC2132 Section 8: Application and Service Parameters
	OptionNISDomain        = Option(40)
	OptionNISServers       = Option(41)
	OptionNTPServers       = Option(42)
	OptionVendorSpecific   = Option(43)
	OptionNETBIOSNameSrv   = Option(44)
	OptionNETBIOSDistSrv   = Option(45)
	OptionNETBIOSNodeType  = Option(46)
	OptionNETBIOSScope     = Option(47)
	OptionXWindowFont      = Option(48)
	OptionXWindowManager   = Option(49)
	OptionNISDomainName    = Option(64)
	OptionNISServerAddr    = Option(65)
	OptionHomeAgentAddrs   = Option(68)
	OptionSMTPServer       = Option(69)
	OptionPOP3Server       = Option(70)
	OptionNNTPServer       = Option(71)
	OptionWWWServer        = Option(72)
	OptionFingerServer     = Option(73)
	OptionIRCServer        = Option(74)
	OptionStreetTalkServer = Option(75)
	OptionSTDAServer       = Option(76)

	// RFC2132 Section 9: DHCP Extensions
	OptionAddressRequest = Option(50)
	OptionAddressTime    = Option(51)
	OptionOverload       = Option(52)
	OptionServerName     = Option(66)
	OptionBootfileName   = Option(67)
	OptionDHCPMsgType    = Option(53)
	OptionDHCPServerID   = Option(54)
	OptionParameterList  = Option(55)
	OptionDHCPMessage    = Option(56)
	OptionDHCPMaxMsgSize = Option(57)
	OptionRenewalTime    = Option(58)
	OptionRebindingTime  = Option(59)
	OptionClassID        = Option(60)
	OptionClientID       = Option(61)
)

// From RFC2241: DHCP Options for Novell Directory Services
const (
	OptionNDSServers  = Option(85)
	OptionNDSTreeName = Option(86)
	OptionNDSContext  = Option(87)
)

// From RFC2242: NetWare/IP Domain Name and Information
const (
	OptionNetWareIPDomain = Option(62)
	OptionNetWareIPOption = Option(63)
)

// From RFC2485: DHCP Option for The Open Group\x27s User Authentication Protocol
const (
	OptionUserAuth = Option(98)
)

// From RFC2563: DHCP Option to Disable Stateless Auto-Configuration in IPv4 Clients
const (
	OptionAutoConfig = Option(116)
)

// From RFC2610: DHCP Options for Service Location Protocol
const (
	OptionDirectoryAgent = Option(78)
	OptionServiceScope   = Option(79)
)

// From RFC2937: The Name Service Search Option for DHCP
const (
	OptionNameServiceSearch = Option(117)
)

// From RFC3004: The User Class Option for DHCP
const (
	OptionUserClass = Option(77)
)

// From RFC3011: The IPv4 Subnet Selection Option for DHCP
const (
	OptionSubnetSelectionOption = Option(118)
)

// From RFC3046: DHCP Relay Agent Information Option
const (
	OptionRelayAgentInformation = Option(82)
)

// From RFC3118: Authentication for DHCP Messages
const (
	OptionAuthentication = Option(90)
)

// From RFC3361: Dynamic Host Configuration Protocol (DHCP-for-IPv4) Option for Session Initiation Protocol (SIP) Servers
const (
	OptionSIPServersDHCPOption = Option(120)
)

// From RFC3397: Dynamic Host Configuration Protocol (DHCP) Domain Search Option
const (
	OptionDomainSearch = Option(119)
)

// From RFC3442: The Classless Static Route Option for Dynamic Host Configuration Protocol (DHCP) version 4
const (
	OptionClasslessStaticRouteOption = Option(121)
)

// From RFC3495: Dynamic Host Configuration Protocol (DHCP) Option for CableLabs Client Configuration
const (
	OptionCCC = Option(122)
)

// From RFC3679: Unused Dynamic Host Configuration Protocol (DHCP) Option Codes
const (
	OptionLDAP           = Option(95)
	OptionNetinfoAddress = Option(112)
	OptionNetinfoTag     = Option(113)
	OptionURL            = Option(114)
)

// From RFC3925: Vendor-Identifying Vendor Options for Dynamic Host Configuration Protocol version 4 (DHCPv4)
const (
	OptionVIVendorClass               = Option(124)
	OptionVIVendorSpecificInformation = Option(125)
)

// From RFC4039: Rapid Commit Option for the Dynamic Host Configuration Protocol version 4 (DHCPv4)
const (
	OptionRapidCommit = Option(80)
)

// From RFC4174: The IPv4 Dynamic Host Configuration Protocol (DHCP) Option for the Internet Storage Name Service
const (
	OptioniSNS = Option(83)
)

// From RFC4280: Dynamic Host Configuration Protocol (DHCP) Options for Broadcast and Multicast Control Servers
const (
	OptionBCMCSControllerDomainNameList    = Option(88)
	OptionBCMCSControllerIPv4AddressOption = Option(89)
)

// From RFC4388: Dynamic Host Configuration Protocol (DHCP) Leasequery
const (
	OptionClientLastTransactionTimeOption = Option(91)
	OptionAssociatedIPOption              = Option(92)
)

// From RFC4578: Dynamic Host Configuration Protocol (DHCP) Options for the Intel Preboot eXecution Environment (PXE)
const (
	OptionClientSystem    = Option(93)
	OptionClientNDI       = Option(94)
	OptionUUIDGUID        = Option(97)
	OptionPXEUndefined128 = Option(128)
	OptionPXEUndefined129 = Option(129)
	OptionPXEUndefined130 = Option(130)
	OptionPXEUndefined131 = Option(131)
	OptionPXEUndefined132 = Option(132)
	OptionPXEUndefined133 = Option(133)
	OptionPXEUndefined134 = Option(134)
	OptionPXEUndefined135 = Option(135)
)

// From RFC4702: The Dynamic Host Configuration Protocol (DHCP) Client Fully Qualified Domain Name (FQDN) Option
const (
	OptionClientFQDN = Option(81)
)

// From RFC4776: Dynamic Host Configuration Protocol (DHCPv4 and DHCPv6) Option for Civic Addresses Configuration Information
const (
	OptionGeoConfCivic = Option(99)
)

// From RFC4833: Timezone Options for DHCP
const (
	OptionPCode = Option(100)
	OptionTCode = Option(101)
)

// From RFC6225: Dynamic Host Configuration Protocol Options for Coordinate-Based Location Configuration Information
const (
	OptionGeoConfOption = Option(123)
	OptionGeoLoc        = Option(144)
)
