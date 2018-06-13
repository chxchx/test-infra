// Copyright 2018 Istio Authors
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

package sisyphus

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/golang/glog"
)

// Storage interface enables additional storage needs for clients besides istio
// and facilitates mocking in tests.
type Storage interface {
	Store(jobName, sha string, newFlakeStat FlakeStat) error
}

// DefaultStorage is empty
type DefaultStorage struct{}

// NewStorage creates a new DefaultStorage
func NewStorage() *DefaultStorage {
	return &DefaultStorage{}
}

// Store records FlakeStat to durable storage
func (s *DefaultStorage) Store(jobName, sha string, newFlakeStat FlakeStat) error {
	log.Printf("newFlakeStat = %v\n", newFlakeStat)
	return nil
}

// SpannerStorage stores flakiness data on cloud spanner
type SpannerStorage struct {
	client spanner.Client
}

// NewSpannerStorage creates a new Storage
func NewSpannerStorage(project, instance, database string) *SpannerStorage {
	db := fmt.Sprintf("projects/%s/instances/%s/database/%s", project, instance, database)
	clnt, err := spanner.NewClient(context.Background(), db)
	if err != nil {
		glog.Fatalf("Unable to connect to spanner db %s", db)
	}
	return &SpannerStorage{
		client: clnt,
	}
}

// Store records FlakeStat to durable storage
func (s *SpannerStorage) Store(jobName, sha string, newFlakeStat FlakeStat) error {
	log.Printf("newFlakeStat = %v\n", newFlakeStat)
	mutation := spanner.InsertMap("flake_stats", map[string]interface{}{
		"test_name":   jobName,
		"sha":         sha,
		"total_rerun": newFlakeStat.TotalRerun,
		"failures":    newFlakeStat.Failures,
	})
	_, err := s.client.Apply(context.Background(), []*spanner.Mutation{mutation})
	return err
}
