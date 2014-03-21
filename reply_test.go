package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReplyFromBytes(t *testing.T) {
	{ // Invalid wire-level representation
		_, err := NewReplyFromBytes([]byte{0x1, 0x2, 0x3})
		if assert.Error(t, err) {
			assert.Equal(t, ErrShortPacket, err)
		}
	}

	{ // Given a request
		b, err := PacketToBytes(NewPacket(BootRequest))
		if assert.Nil(t, err) {
			_, err := NewReplyFromBytes(b)
			if assert.Error(t, err) {
				assert.Equal(t, ErrNoReply, err)
			}
		}
	}

	{ // Given a reply
		b, err := PacketToBytes(NewPacket(BootReply))
		if assert.Nil(t, err) {
			req, err := NewReplyFromBytes(b)
			if assert.Nil(t, err) {
				assert.NotNil(t, req)
			}
		}
	}
}
