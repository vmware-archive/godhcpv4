package dhcpv4

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptionMapOption(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b []byte

	om := make(OptionMap)

	_, ok = om.GetOption(o)
	assert.False(t, ok)

	a = []byte("foo")
	om.SetOption(o, a)

	b, ok = om.GetOption(o)
	assert.True(t, ok)
	assert.Equal(t, a, b)
}

func TestOptionMapMessageType(t *testing.T) {
	var a, b MessageType

	om := make(OptionMap)

	b = om.GetMessageType()
	assert.Equal(t, MessageType(0), b)

	a = MessageType(1)
	om.SetMessageType(a)

	b = om.GetMessageType()
	assert.Equal(t, a, b)
}

func TestOptionMapUint8(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b uint8

	om := make(OptionMap)

	_, ok = om.GetUint8(o)
	assert.False(t, ok)

	a = uint8(37)
	om.SetUint8(o, a)

	b, ok = om.GetUint8(o)
	assert.True(t, ok)
	assert.Equal(t, a, b)
}

func TestOptionMapUint16(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b uint16

	om := make(OptionMap)

	_, ok = om.GetUint16(o)
	assert.False(t, ok)

	a = uint16(37000)
	om.SetUint16(o, a)

	b, ok = om.GetUint16(o)
	assert.True(t, ok)
	assert.Equal(t, a, b)
}

func TestOptionMapUint32(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b uint32

	om := make(OptionMap)

	_, ok = om.GetUint32(o)
	assert.False(t, ok)

	a = uint32(37000000)
	om.SetUint32(o, a)

	b, ok = om.GetUint32(o)
	assert.True(t, ok)
	assert.Equal(t, a, b)
}

func TestOptionMapString(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b string

	om := make(OptionMap)

	_, ok = om.GetString(o)
	assert.False(t, ok)

	a = "hello world!"
	om.SetString(o, a)

	b, ok = om.GetString(o)
	assert.True(t, ok)
	assert.Equal(t, a, b)
}

func TestOptionMapIP(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b net.IP

	om := make(OptionMap)

	_, ok = om.GetIP(o)
	assert.False(t, ok)

	a = net.IPv4(1, 2, 3, 4)
	om.SetIP(o, a)

	b, ok = om.GetIP(o)
	assert.True(t, ok)
	assert.Equal(t, a, b)
}

func TestOptionMapDuration(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b time.Duration

	om := make(OptionMap)

	_, ok = om.GetDuration(o)
	assert.False(t, ok)

	a = 100 * time.Second
	om.SetDuration(o, a)

	b, ok = om.GetDuration(o)
	assert.True(t, ok)
	assert.Equal(t, a, b)
}

func TestOptionMapDurationTruncateSubSecond(t *testing.T) {
	var o = Option(1)
	var ok bool
	var a, b time.Duration

	om := make(OptionMap)

	a = 100*time.Second + 100*time.Millisecond
	om.SetDuration(o, a)

	b, ok = om.GetDuration(o)
	assert.True(t, ok)
	assert.Equal(t, 100*time.Second, b)
}

func TestOptionMapDecodeWithoutPtr(t *testing.T) {
	om := make(OptionMap)

	var s struct {
		U8  uint8  `code:"1"`
		U16 uint16 `code:"2"`
		U32 uint32 `code:"3"`
		S   string `code:"4"`
	}

	om.SetUint8(Option(1), 37)
	om.SetUint16(Option(2), 37000)
	om.SetUint32(Option(3), 37000000)
	om.SetString(Option(4), "thirtyseven")

	om.Decode(&s)

	assert.Equal(t, uint8(37), s.U8)
	assert.Equal(t, uint16(37000), s.U16)
	assert.Equal(t, uint32(37000000), s.U32)
	assert.Equal(t, "thirtyseven", s.S)
}

func TestOptionMapDecodeWithPtr(t *testing.T) {
	om := make(OptionMap)

	var s struct {
		U8  *uint8  `code:"1"`
		U16 *uint16 `code:"2"`
		U32 *uint32 `code:"3"`
		S   *string `code:"4"`
	}

	om.SetUint8(Option(1), 37)
	om.SetUint16(Option(2), 37000)
	om.SetUint32(Option(3), 37000000)
	om.SetString(Option(4), "thirtyseven")

	om.Decode(&s)

	if assert.NotNil(t, s.U8) {
		assert.Equal(t, uint8(37), *s.U8)
	}
	if assert.NotNil(t, s.U16) {
		assert.Equal(t, uint16(37000), *s.U16)
	}
	if assert.NotNil(t, s.U32) {
		assert.Equal(t, uint32(37000000), *s.U32)
	}
	if assert.NotNil(t, s.S) {
		assert.Equal(t, "thirtyseven", *s.S)
	}
}
