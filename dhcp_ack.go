/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package dhcpv4

import "encoding/binary"

// DHCPAck is a server to client packet with configuration parameters,
// including committed network address.
type DHCPAck struct {
	Packet

	req Request
}

func CreateDHCPAck(req Request) DHCPAck {
	rep := DHCPAck{
		Packet: NewReply(req),
		req:    req,
	}

	rep.SetMessageType(MessageTypeDHCPAck)
	return rep
}

// From RFC2131, table 3:
//   Option                    DHCPACK
//   ------                    -------
//   Requested IP address      MUST NOT
//   IP address lease time     MUST (DHCPREQUEST)
//                             MUST NOT (DHCPINFORM)
//   Use 'file'/'sname' fields MAY
//   DHCP message type         DHCPACK
//   Parameter request list    MUST NOT
//   Message                   SHOULD
//   Client identifier         MUST NOT
//   Vendor class identifier   MAY
//   Server identifier         MUST
//   Maximum message size      MUST NOT
//   All others                MAY

var dhcpAckOnRequestValidation = []Validation{
	ValidateMust(OptionAddressTime),
}

var dhcpAckOnInformValidation = []Validation{
	ValidateMustNot(OptionAddressTime),
}

var dhcpAckValidation = []Validation{
	ValidateMustNot(OptionAddressRequest),
	ValidateMustNot(OptionParameterList),
	ValidateMustNot(OptionClientID),
	ValidateMust(OptionDHCPServerID),
	ValidateMustNot(OptionDHCPMaxMsgSize),
}

func (d DHCPAck) Validate() error {
	var err error

	// Validation is subtly different based on type of request
	switch d.req.GetMessageType() {
	case MessageTypeDHCPRequest:
		err = Validate(d.Packet, dhcpAckOnRequestValidation)
	case MessageTypeDHCPInform:
		err = Validate(d.Packet, dhcpAckOnInformValidation)
	}

	if err != nil {
		return err
	}

	return Validate(d.Packet, dhcpAckValidation)
}

func (d DHCPAck) ToBytes() ([]byte, error) {
	opts := packetToBytesOptions{}

	// Copy MaxMsgSize if set in the request
	if v, ok := d.Request().GetOption(OptionDHCPMaxMsgSize); ok {
		opts.maxLen = binary.BigEndian.Uint16(v)
	}

	return PacketToBytes(d.Packet, &opts)
}

func (d DHCPAck) Request() Request {
	return d.req
}
