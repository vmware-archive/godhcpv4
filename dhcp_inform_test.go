package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDHCPInformCreateDHCPAck(t *testing.T) {
	req := DHCPInform{
		Packet: NewPacket(BootRequest),
	}

	rep := req.CreateDHCPAck()
	assert.Equal(t, MessageTypeDHCPAck, rep.GetMessageType())
}

// Test dispatch to ReplyWriter
func TestDHCPInformWriteReply(t *testing.T) {
	rw := &testReplyWriter{}

	req := DHCPInform{
		Packet:      NewPacket(BootRequest),
		ReplyWriter: rw,
	}

	reps := []Reply{
		req.CreateDHCPAck(),
	}

	for _, rep := range reps {
		rw.wrote = false
		req.WriteReply(rep)
		assert.True(t, rw.wrote)
	}
}
