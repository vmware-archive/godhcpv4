package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDHCPRequestCreateDHCPAck(t *testing.T) {
	req := DHCPRequest{
		Packet: NewPacket(BootRequest),
	}

	rep := req.CreateDHCPAck()
	assert.Equal(t, MessageTypeDHCPAck, rep.GetMessageType())
}

func TestDHCPRequestCreateDHCPNak(t *testing.T) {
	req := DHCPRequest{
		Packet: NewPacket(BootRequest),
	}

	rep := req.CreateDHCPNak()
	assert.Equal(t, MessageTypeDHCPNak, rep.GetMessageType())
}

// Test dispatch to ReplyWriter
func TestDHCPRequestWriteReply(t *testing.T) {
	rw := &testReplyWriter{}

	req := DHCPRequest{
		Packet:      NewPacket(BootRequest),
		ReplyWriter: rw,
	}

	reps := []Reply{
		req.CreateDHCPAck(),
		req.CreateDHCPNak(),
	}

	for _, rep := range reps {
		rw.wrote = false
		req.WriteReply(rep)
		assert.True(t, rw.wrote)
	}
}
