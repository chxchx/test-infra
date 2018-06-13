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
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	u "istio.io/test-infra/toolbox/util"
)

const (
	project  = "istio"
	instance = "sisyphus"
	database = "sisyphus-db"
)

func TestGetLatestRunOnProw(t *testing.T) {
	storage := NewSpannerStorage(project, instance, database)
	expectedJob := "job"
	expectedSHA := 123
	expectedRerun := 3
	expectedFailures := 1
	expectedStat := FlakeStat{
		TestName:   expectedJob,
		SHA:        expectedSHA,
		TotalRerun: expectedRerun,
		Failures:   expectedFailures,
	}
	if storage.Store(expectedJob, expectedSHA, expectedStat) != nil {
		t.Errorf("failed to store: %v", err)
	}
}
