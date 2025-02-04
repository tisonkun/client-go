// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	assert := assert.New(t)

	done := make(chan struct{}, 1)
	rl := NewRateLimit(1)
	assert.PanicsWithValue("put a redundant token", rl.PutToken)
	exit := rl.GetToken(done)
	assert.False(exit)
	rl.PutToken()
	assert.PanicsWithValue("put a redundant token", rl.PutToken)

	exit = rl.GetToken(done)
	assert.False(exit)
	done <- struct{}{}
	exit = rl.GetToken(done) // blocked but exit
	assert.True(exit)

	sig := make(chan int, 1)
	go func() {
		exit = rl.GetToken(done) // blocked
		assert.False(exit)
		close(sig)
	}()
	time.Sleep(200 * time.Millisecond)
	rl.PutToken()
	<-sig
}
