// Copyright 2023 Manu Garg.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protodoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKindToURL(t *testing.T) {
	tests := []struct {
		kind string
		want string
	}{
		{
			kind: "cloudprober.probes.ProbeDef.interval_msec",
			want: "probes.html#cloudprober.probes.ProbeDef.interval_msec",
		},
		{
			kind: "cloudprober.interval_msec",
			want: "",
		},
		{
			kind: "probes.ProbeDef.interval_msec",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			assert.Equal(t, tt.want, kindToURL(tt.kind))
		})
	}
}

func TestArrangeIntoPackages(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  map[string][]string
	}{
		{
			paths: []string{
				"cloudprober.probes.ProbeDef.interval_msec",
				"cloudprober.probes.ProbeDef.timeout_msec",
				"cloudprober.probes.http.ProbeDef.header",
			},
			want: map[string][]string{
				"probes": []string{
					"cloudprober.probes.ProbeDef.interval_msec",
					"cloudprober.probes.ProbeDef.timeout_msec",
					"cloudprober.probes.http.ProbeDef.header",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ArrangeIntoPackages(tt.paths, nil))
		})
	}
}
