package dhcpv4

import "errors"

var (
	ErrShortPacket = errors.New("dhcpv4: short packet")
)

type OpCode byte

// Message op codes defined in RFC2132.
const (
	BootRequest = OpCode(1)
	BootReply   = OpCode(2)
)

type MessageType byte

// Message types defined in RFC2132.
const (
	DhcpDiscover = MessageType(1)
	DhcpOffer    = MessageType(2)
	DhcpRequest  = MessageType(3)
	DhcpDecline  = MessageType(4)
	DhcpAck      = MessageType(5)
	DhcpNak      = MessageType(6)
	DhcpRelease  = MessageType(7)
	DhcpInform   = MessageType(8)
)

type RawPacket []byte

func (p RawPacket) Op() []byte     { return p[0:1] }
func (p RawPacket) HType() []byte  { return p[1:2] }
func (p RawPacket) HLen() []byte   { return p[2:3] }
func (p RawPacket) Hops() []byte   { return p[3:4] }
func (p RawPacket) XId() []byte    { return p[4:8] }
func (p RawPacket) Secs() []byte   { return p[8:10] }
func (p RawPacket) Flags() []byte  { return p[10:12] }
func (p RawPacket) CIAddr() []byte { return p[12:16] }
func (p RawPacket) YIAddr() []byte { return p[16:20] }
func (p RawPacket) SIAddr() []byte { return p[20:24] }
func (p RawPacket) GIAddr() []byte { return p[24:28] }
func (p RawPacket) CHAddr() []byte { return p[28:44] }

// SName returns the `sname` portion of the packet.
// This field can be used as extra space to extend the DHCP options, if
// necessary. To enable this, the "Option Overload" option needs to be set in
// the regular options. Also see RFC2132, section 9.3.
func (p RawPacket) SName() []byte {
	return p[44:108]
}

// File returns the `file` portion of the packet.
// This field can be used as extra space to extend the DHCP options, if
// necessary. To enable this, the "Option Overload" option needs to be set in
// the regular options. Also see RFC2132, section 9.3.
func (p RawPacket) File() []byte {
	return p[108:236]
}

// Cookie returns the fixed-value prefix to the `options` portion of the packet.
// According to the RFC, this should equal the 4-octet { 99, 130, 83, 99 }.
func (p RawPacket) Cookie() []byte {
	return p[236:240]
}

// Options returns the variable-sized `options` portion of the packet.
func (p RawPacket) Options() []byte {
	return p[240:]
}

func parseOptionBuffer(x []byte, opts OptionMap) error {
	for {
		if len(x) == 0 {
			return ErrShortPacket
		}

		tag := Option(x[0])
		x = x[1:]
		if tag == OptionEnd {
			break
		}

		// Padding tag
		if tag == OptionPad {
			continue
		}

		// Read length octet
		if len(x) == 0 {
			return ErrShortPacket
		}

		length := int(x[0])
		x = x[1:]
		if len(x) < length {
			return ErrShortPacket
		}

		_, ok := opts[tag]
		if ok {
			// We've got a bad client here; duplicate options are not allowed.
			// Let it slide instead of throwing a fit, for the sake of robustness.
		}

		// Capture option and move to the next one
		opts[tag] = x[0:length]
		x = x[length:]
	}

	return nil
}

func (p RawPacket) ParseOptions() (OptionMap, error) {
	var b []byte
	var err error

	// Facilitate up to 255 option tags
	opts := make(OptionMap, 255)

	// Parse initial set of options
	b = p.Options()
	if err = parseOptionBuffer(b, opts); err != nil {
		return nil, err
	}

	// Parse options from `file` field if necessary
	if x := opts[OptionOverload]; len(x) > 0 && x[0]&0x1 != 0 {
		b = p.File()
		if err = parseOptionBuffer(b, opts); err != nil {
			return nil, err
		}
	}

	// Parse options from `sname` field if necessary
	if x := opts[OptionOverload]; len(x) > 0 && x[0]&0x2 != 0 {
		b = p.SName()
		if err = parseOptionBuffer(b, opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}
