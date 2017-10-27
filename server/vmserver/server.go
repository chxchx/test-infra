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
	"container/list"
	"log"
	"net"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "istio.io/test-infra/server/vmserver/proto"
)

const (
	port = ":50051"
)

var (
	// doubly linked list to implement round robin access to VMs in the pool
	vmPool     = list.New()
	vmPoolLock = &sync.Mutex{}
)

type vm struct {
	name string
	lock *sync.Mutex
	// TODO (chx) all other GCP configs, pending merging of mono repo done
}

// select a VM in round robin fashion
func selectVM() vm {
	vmPoolLock.Lock()
	defer vmPoolLock.Unlock()
	e := vmPool.Front()
	selectedVM, _ := e.Value.(vm)
	vmPool.MoveToBack(e)
	return selectedVM
}

// find the vm struct from the vmPool
func findVM(vmName string) vm {
	vmPoolLock.Lock()
	defer vmPoolLock.Unlock()
	for e := vmPool.Front(); e != nil; e = e.Next() {
		castToVM, _ := e.Value.(vm)
		if castToVM.name == vmName {
			return castToVM
		}
	}
	return nil
}

type server struct{}

func (s *server) Acquire(ctx context.Context, cfg *pb.VMConfig) (*pb.VMInstance, error) {
	selectedVM := selectVM()
	selectedVM.lock.Lock()
	log.Printf("Received Acquire() request from client %s\n", cfg.ClientNamespace)
	return &pb.VMInstance{VmName: "Hello " + cfg.ClientNamespace}, nil
}

func (s *server) Release(ctx context.Context, vmInstance *pb.VMInstance) (*pb.Done, error) {
	selectedVM := findVM(vmInstance.VmName)
	defer selectedVM.lock.Unlock()
	return &pb.Done{}, nil
}

func (s *server) KeepAlive(ctx context.Context, in *pb.Client) (*pb.Done, error) {
	return &pb.Done{}, nil
}

func init() {
	// TODO (chx) populate vmPool, recover states from durable storage
	vmPool.PushFront(vm{
		name: "raw-vm-test-new",
		lock: &sync.Mutex{},
	})
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPreprovisionServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
