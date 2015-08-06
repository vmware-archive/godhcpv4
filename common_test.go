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

type ValidatingReply interface {
	SetOption(o Option, v []byte)
	Validate() error
}

type replyValidationTestCase struct {
	newReply func() ValidatingReply
	must     []Option
	mustNot  []Option
}

func (r replyValidationTestCase) Test(t *testing.T) {
	var err error

	mustOptions := r.must
	mustNotOptions := r.mustNot

	// Forget each individual one
	for i, o := range mustOptions {
		reply := r.newReply()

		// Add options not to be tested here
		for j, o := range mustOptions {
			if i != j {
				reply.SetOption(o, []byte("foo"))
			}
		}

		// Fail validation without the option
		err = reply.Validate()
		assert.Error(t, err)

		// Pass validation with the option
		reply.SetOption(o, []byte("foo"))
		err = reply.Validate()
		assert.NoError(t, err)
	}

	// Add each individual one
	for _, o := range mustNotOptions {
		reply := r.newReply()

		// Add options not to be tested here
		for _, o := range mustOptions {
			reply.SetOption(o, []byte("foo"))
		}

		// Pass validation without the option
		err = reply.Validate()
		assert.NoError(t, err)

		// Fail validation with the option
		reply.SetOption(o, []byte("foo"))
		err = reply.Validate()
		assert.Error(t, err)
	}
}

type testReplyWriter struct {
	wrote bool
}

func (t *testReplyWriter) WriteReply(r Reply) error {
	t.wrote = true
	return nil
}
