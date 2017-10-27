// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// TODO (chx) set up a server on local that has only one vm

func inspect(err error, fMsg, sMsg string, t *testing.T) {
	if err != nil {
		log.Printf("%s. Error %s\n", fMsg, err)
		t.Error(err)
	} else if sMsg != "" {
		log.Println(sMsg)
	}
}

// TestMutualExclusion tests if Acquire() and Release() ensures serializability
func TestMutualExclusion(t *testing.T) {
	alice, err := NewVMClient("alice")
	inspect(err, "", "", t)
	bob, err := NewVMClient("bob")
	inspect(err, "", "", t)
	sharedBuf := 1
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		vmInstance, err := alice.Acquire()
		inspect(err, "alice failed to acquire", "", t)
		sharedBuf = 2
		err = alice.Release(vmInstance)
		inspect(err, "alice failed to release", "", t)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		vmInstance, err := bob.Acquire()
		inspect(err, "bob failed to acquire", "", t)
		if sharedBuf != 2 {
			inspect(fmt.Errorf("bob acquired vm before alice released it"), "", "", t)
		}
		sharedBuf = 3
		err = alice.Release(vmInstance)
		inspect(err, "bob failed to release", "", t)
	}()
	wg.Wait()
}
