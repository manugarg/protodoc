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
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestMain(m *testing.M) {
	BuildFileDescRegistry(Files, "testdata", "github.com/manugarg/protodoc", nil)
	m.Run()
}

func TestDumpMessage(t *testing.T) {
	const fldName = "cloudprober.probes.ProbeDef.http_probe"

	tests := []struct {
		name     string
		f        Formatter
		wantToks []*Token
	}{
		{
			name: "default",
			wantToks: []*Token{
				{
					Kind: "cloudprober.probes.http.Header",
					Text: "header",
				},
			},
		},
		{
			name: "depth=2",
			f:    Formatter{}.WithDepth(2),
			wantToks: []*Token{
				{
					Kind:          "cloudprober.probes.http.Header",
					Text:          "header",
					MessageHeader: true,
					NoExtraLine:   true,
				},
				{
					Kind:   "string",
					Text:   "name",
					Prefix: "  ",
				},
				{
					Kind:        "string",
					Text:        "value",
					Prefix:      "  ",
					NoExtraLine: true,
				},
				{
					Kind: "",
					Text: "}",
				},
			},
		},
		{
			name: "yaml,depth=2",
			f:    Formatter{}.WithDepth(2).WithYAML(true, true),
			wantToks: []*Token{
				{
					Kind:          "cloudprober.probes.http.Header",
					Text:          "header",
					MessageHeader: true,
					NoExtraLine:   true,
					yaml:          true,
				},
				{
					Kind:   "string",
					Text:   "name",
					Prefix: "  - ",
					yaml:   true,
				},
				{
					Kind:   "string",
					Text:   "value",
					Prefix: "    ",
					yaml:   true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := Files.FindDescriptorByName(protoreflect.FullName(fldName))
			assert.NoError(t, err)

			md := desc.(protoreflect.FieldDescriptor).Message()
			toks, _ := DumpMessage(md, tt.f)

			assert.Equal(t, tt.wantToks, toks)
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
