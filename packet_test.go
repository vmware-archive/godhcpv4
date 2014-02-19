package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPacket struct {
	buf []byte

	offsetOption int
	offsetFile   int
	offsetSName  int
}

func (t *testPacket) writeOptionAtOffset(offset int, o Option, v []byte) int {
	var lv int

	p := t.buf
	switch o {
	case OptionPad, OptionEnd:
		lv = 1
	default:
		lv = 2 + len(v)
	}

	if offset+lv > cap(p) {
		p = make([]byte, offset+lv)
		copy(p, t.buf)
		t.buf = p
	}

	// Write option to offset
	q := p[offset : offset+lv]
	q[0] = byte(o)
	if lv > 1 {
		q[1] = byte(len(v))
		copy(q[2:], v)
	}

	return len(q)
}

func (t *testPacket) appendToOption(o Option, v []byte) int {
	if t.offsetOption == 0 {
		t.offsetOption = 240
	}

	n := t.writeOptionAtOffset(t.offsetOption, o, v)
	t.offsetOption += n
	return n
}

func (t *testPacket) appendToFile(o Option, v []byte) int {
	if t.offsetFile == 0 {
		t.offsetFile = 108
	}

	n := t.writeOptionAtOffset(t.offsetFile, o, v)
	t.offsetFile += n
	if t.offsetFile > 236 {
		panic("overflow in file field")
	}

	return n
}

func (t *testPacket) appendToSName(o Option, v []byte) int {
	if t.offsetSName == 0 {
		t.offsetSName = 44
	}

	n := t.writeOptionAtOffset(t.offsetSName, o, v)
	t.offsetSName += n
	if t.offsetSName > 108 {
		panic("overflow in sname field")
	}

	return n
}

func parseOptions(t *testing.T, p []byte) (OptionMap, error) {
	var o OptionMap
	var err error

	for i := 240; i < len(p); i++ {
		r := RawPacket(p[0 : i+1])
		o, err = r.ParseOptions()
		if len(r) == len(p) {
			assert.Nil(t, err, "expected no error with i=%d", i)
		} else {
			assert.Equal(t, err, ErrShortPacket, "expect short packet error with i=%d", i)
		}
	}

	return o, err
}

func assertOption(t *testing.T, opts OptionMap, o Option, expected []byte) {
	actual, ok := opts[o]
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestPacketParseOptions(t *testing.T) {
	var p *testPacket
	var o OptionMap
	var err error

	// Fabricate packet
	p = new(testPacket)
	p.appendToOption(OptionSubnetMask, []byte{0x12})
	p.appendToOption(OptionEnd, nil)
	if o, err = parseOptions(t, p.buf); assert.Nil(t, err) {
		assert.Equal(t, 1, len(o))
		assertOption(t, o, OptionSubnetMask, []byte{0x12})
	}

	// Fabricate packet with overload into `file` field
	p = new(testPacket)
	p.appendToOption(OptionSubnetMask, []byte{0x12})
	p.appendToOption(OptionOverload, []byte{0x1})
	p.appendToOption(OptionEnd, nil)
	p.appendToFile(OptionTimeOffset, []byte{0x34})
	p.appendToFile(OptionEnd, nil)
	if o, err = parseOptions(t, p.buf); assert.Nil(t, err) {
		assert.Equal(t, 3, len(o))
		assertOption(t, o, OptionSubnetMask, []byte{0x12})
		assertOption(t, o, OptionOverload, []byte{0x1})
		assertOption(t, o, OptionTimeOffset, []byte{0x34})
	}

	// Fabricate packet with overload into `sname` field
	p = new(testPacket)
	p.appendToOption(OptionSubnetMask, []byte{0x12})
	p.appendToOption(OptionOverload, []byte{0x2})
	p.appendToOption(OptionEnd, nil)
	p.appendToSName(OptionRouter, []byte{0x56})
	p.appendToSName(OptionEnd, nil)
	if o, err = parseOptions(t, p.buf); assert.Nil(t, err) {
		assert.Equal(t, 3, len(o))
		assertOption(t, o, OptionSubnetMask, []byte{0x12})
		assertOption(t, o, OptionOverload, []byte{0x2})
		assertOption(t, o, OptionRouter, []byte{0x56})
	}

	// Fabricate packet with overload into `file` AND `sname` fields
	p = new(testPacket)
	p.appendToOption(OptionSubnetMask, []byte{0x12})
	p.appendToOption(OptionOverload, []byte{0x3})
	p.appendToOption(OptionEnd, nil)
	p.appendToFile(OptionTimeOffset, []byte{0x34})
	p.appendToFile(OptionEnd, nil)
	p.appendToSName(OptionRouter, []byte{0x56})
	p.appendToSName(OptionEnd, nil)
	if o, err = parseOptions(t, p.buf); assert.Nil(t, err) {
		assert.Equal(t, 4, len(o))
		assertOption(t, o, OptionSubnetMask, []byte{0x12})
		assertOption(t, o, OptionOverload, []byte{0x3})
		assertOption(t, o, OptionTimeOffset, []byte{0x34})
		assertOption(t, o, OptionRouter, []byte{0x56})
	}
}
