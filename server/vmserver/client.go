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
	// "log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "istio.io/test-infra/server/vmserver/proto"
)

const (
	address = "localhost:50051"
)

// VMClient exposes managed interface to interact with VM server
type VMClient struct {
	grpcClnt  pb.PreprovisionClient
	namespace string
}

// NewVMClient creates a new VMClient
func NewVMClient(namespace string) (*VMClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to server: %v", err)
	}
	return &VMClient{
		namespace: namespace,
		grpcClnt:  pb.NewPreprovisionClient(conn),
	}, nil
}

// Acquire requests from server exclusive access to a VM
func (c *VMClient) Acquire() (*pb.VMInstance, error) {
	// TODO (chx) track expiration
	return c.grpcClnt.Acquire(context.Background(), &pb.VMConfig{ClientNamespace: c.namespace})
}

// Release a VM back to the pool
func (c *VMClient) Release(vmInstance *pb.VMInstance) error {
	_, err := c.grpcClnt.Release(context.Background(), vmInstance)
	return err
}

// TODO (chx) keep alive
