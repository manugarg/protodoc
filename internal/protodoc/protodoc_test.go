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
	"html/template"
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
					Kind:   "cloudprober.probes.http.Header",
					Text:   "header",
					Suffix: template.HTML(" {"),
				},
				{
					Kind:   "string",
					Text:   "name",
					Prefix: "  ",
				},
				{
					Kind:   "string",
					Text:   "value",
					Prefix: "  ",
					Suffix: template.HTML("<br>"),
				},
				{
					Kind: "",
					Text: "}",
				},
			},
		},
		{
			name: "yaml,depth=2",
			f:    Formatter{}.WithDepth(2).WithYAML(true),
			wantToks: []*Token{
				{
					Kind:   "cloudprober.probes.http.Header",
					Text:   "header",
					Suffix: template.HTML(":"),
				},
				{
					Kind:   "string",
					Text:   "name",
					Prefix: "  - ",
				},
				{
					Kind:   "string",
					Text:   "value",
					Prefix: "    ",
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

func TestProcessTokensForHTML(t *testing.T) {
	type args struct {
		toks []*Token
	}
	tests := []struct {
		name string
		in   *Token
		want *Token
	}{
		{
			name: "simple",
			in: &Token{
				Kind:   "string",
				Prefix: "  ",
			},
			want: &Token{
				Kind:       "string",
				Prefix:     "  ",
				PrefixHTML: "&nbsp;&nbsp;",
				Suffix:     "<br><br>",
			},
		},
		{
			name: "with-url",
			in: &Token{
				Kind:   "cloudprober.probes.ProbeDef",
				Prefix: "  ",
			},
			want: &Token{
				Kind:       "cloudprober.probes.ProbeDef",
				URL:        "probes.html#cloudprober.probes.ProbeDef",
				Prefix:     "  ",
				PrefixHTML: "&nbsp;&nbsp;",
				Suffix:     "<br><br>",
			},
		},
		{
			name: "with-default",
			in: &Token{
				Kind:    "string",
				Prefix:  "  ",
				Default: "2s",
			},
			want: &Token{
				Kind:       "string",
				Prefix:     "  ",
				PrefixHTML: "&nbsp;&nbsp;",
				Suffix:     " | default: 2s<br><br>",
				Default:    "2s",
			},
		},
		{
			name: "existing-br-suffix",
			in: &Token{
				Kind:   "string",
				Prefix: "  ",
				Suffix: "<br>",
			},
			want: &Token{
				Kind:       "string",
				Prefix:     "  ",
				PrefixHTML: "&nbsp;&nbsp;",
				Suffix:     "<br>",
			},
		},
		{
			name: "existing-br-suffix",
			in: &Token{
				Kind:   "string",
				Prefix: "  ",
				Suffix: "}",
			},
			want: &Token{
				Kind:       "string",
				Prefix:     "  ",
				PrefixHTML: "&nbsp;&nbsp;",
				Suffix:     "}<br>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, []*Token{tt.want}, ProcessTokensForHTML([]*Token{tt.in}))
		})
	}
}
