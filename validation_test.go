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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMustNot(t *testing.T) {
	var err error

	p := NewPacket(BootReply)
	v := ValidateMustNot(OptionSubnetMask)

	err = Validate(p, []Validation{v})
	assert.NoError(t, err)

	p.SetOption(OptionSubnetMask, []byte("something"))
	err = Validate(p, []Validation{v})
	assert.Error(t, err)
}

func TestValidateMust(t *testing.T) {
	var err error

	p := NewPacket(BootReply)
	v := ValidateMust(OptionSubnetMask)

	err = Validate(p, []Validation{v})
	assert.Error(t, err)

	p.SetOption(OptionSubnetMask, []byte("something"))
	err = Validate(p, []Validation{v})
	assert.NoError(t, err)
}

func TestValidateAllowedOptions(t *testing.T) {
	var err error
	var options = []Option{OptionSubnetMask, OptionTimeOffset}

	p := NewPacket(BootReply)
	v := ValidateAllowedOptions(options)

	// Empty
	err = Validate(p, []Validation{v})
	assert.NoError(t, err)

	// Both allowed options
	for _, o := range options {
		p.SetOption(o, []byte("foo"))
		err = Validate(p, []Validation{v})
		assert.NoError(t, err)
	}

	// Forbidden option
	p.SetOption(OptionRouter, []byte("foo"))
	err = Validate(p, []Validation{v})
	assert.Error(t, err)
}
