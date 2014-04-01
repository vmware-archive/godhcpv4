package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDHCPDiscoverCreateDHCPOffer(t *testing.T) {
	req := DHCPDiscover{
		Packet: NewPacket(BootRequest),
	}

	rep := req.CreateDHCPOffer()
	assert.Equal(t, MessageTypeDHCPOffer, rep.GetMessageType())
}

// Test dispatch to ReplyWriter
func TestDHCPDiscoverWriteReply(t *testing.T) {
	rw := &testReplyWriter{}

	req := DHCPDiscover{
		Packet:      NewPacket(BootRequest),
		ReplyWriter: rw,
	}

	reps := []Reply{
		req.CreateDHCPOffer(),
	}

	for _, rep := range reps {
		rw.wrote = false
		req.WriteReply(rep)
		assert.True(t, rw.wrote)
	}
}
