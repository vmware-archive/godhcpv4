package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestFromBytes(t *testing.T) {
	{ // Invalid wire-level representation
		_, err := NewRequestFromBytes([]byte{0x1, 0x2, 0x3})
		if assert.Error(t, err) {
			assert.Equal(t, ErrShortPacket, err)
		}
	}

	{ // Given a reply
		b, err := PacketToBytes(NewPacket(BootReply))
		if assert.Nil(t, err) {
			_, err := NewRequestFromBytes(b)
			if assert.Error(t, err) {
				assert.Equal(t, ErrNoRequest, err)
			}
		}
	}

	{ // Given a request
		b, err := PacketToBytes(NewPacket(BootRequest))
		if assert.Nil(t, err) {
			req, err := NewRequestFromBytes(b)
			if assert.Nil(t, err) {
				assert.NotNil(t, req)
			}
		}
	}
}

func TestCreateReply(t *testing.T) {
	req := Request(NewPacket(BootRequest))
	copy(req.XId(), []byte{0, 1, 2, 3})

	rep := req.CreateReply()
	assert.Equal(t, byte(1), rep.HType()[0])
	assert.Equal(t, byte(6), rep.HLen()[0])
	assert.Equal(t, []byte{0, 1, 2, 3}, rep.XId())
}
