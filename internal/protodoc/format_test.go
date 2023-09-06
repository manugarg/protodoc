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

func TestFinalToToken(t *testing.T) {
	const fldName = "cloudprober.probes.ProbeDef.interval_msec"

	tests := []struct {
		name      string
		f         Formatter
		nocomment bool
		want      *Token
	}{
		{
			name:      "no comment",
			nocomment: true,
			want: &Token{
				Kind:    "int32",
				Comment: "",
				Text:    "interval_msec",
			},
		},
		{
			name:      "not yaml",
			f:         Formatter{},
			nocomment: false,
			want: &Token{
				Kind:    "int32",
				Comment: "# Interval between two probe runs in milliseconds.\n# Only one of \"interval\" and \"inteval_msec\" should be defined.\n# Default interval is 2s.",
				Text:    "interval_msec",
			},
		},
		{
			name:      "yaml",
			f:         Formatter{}.WithYAML(true),
			nocomment: false,
			want: &Token{
				yaml:    true,
				Kind:    "int32",
				Comment: "# Interval between two probe runs in milliseconds.\n# Only one of \"interval\" and \"inteval_msec\" should be defined.\n# Default interval is 2s.",
				Text:    "intervalMsec",
			},
		},
		{
			name:      "yaml with prefix",
			f:         Formatter{}.WithYAML(true).WithPrefix("  "),
			nocomment: false,
			want: &Token{
				yaml:    true,
				Kind:    "int32",
				Prefix:  "  ",
				Comment: "  # Interval between two probe runs in milliseconds.\n  # Only one of \"interval\" and \"inteval_msec\" should be defined.\n  # Default interval is 2s.",
				Text:    "intervalMsec",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := Files.FindDescriptorByName(protoreflect.FullName(fldName))
			assert.NoError(t, err)

			fld := desc.(protoreflect.FieldDescriptor)
			assert.Equal(t, tt.want, finalToken(fld, tt.f, tt.nocomment))
		})
	}
}

func TestFormatOneOf(t *testing.T) {
	const fldName = "cloudprober.probes.ProbeDef.http_probe"
	// Used for duplication test for oneof fileds.
	const fldName2 = "cloudprober.probes.ProbeDef.dns_probe"

	tests := []struct {
		name string
		f    Formatter
		want *Token
	}{
		{
			name: "default",
			want: &Token{
				Kind: "oneof",
				TextHTML: `[http_probe &lt;<a href="probes#cloudprober_probes_http_ProbeConf">cloudprober.probes.http.ProbeConf</a>&gt; | dns_probe &lt;<a href="probes#cloudprober_probes_dns_ProbeConf">cloudprober.probes.dns.ProbeConf</a>&gt; | 
&nbsp;user_defined_probe &lt;string&gt;]`,
			},
		},
		{
			name: "yaml",
			f: Formatter{
				yaml: true,
			},
			want: &Token{
				Kind: "oneof",
				TextHTML: `[httpProbe &lt;<a href="probes#cloudprober_probes_http_ProbeConf">cloudprober.probes.http.ProbeConf</a>&gt; | dnsProbe &lt;<a href="probes#cloudprober_probes_dns_ProbeConf">cloudprober.probes.dns.ProbeConf</a>&gt; | 
&nbsp;userDefinedProbe &lt;string&gt;]`,
			},
		},
		{
			name: "with-prefix",
			f: Formatter{
				prefix: "  ",
			},
			want: &Token{
				Kind:   "oneof",
				Prefix: "  ",
				TextHTML: `[http_probe &lt;<a href="probes#cloudprober_probes_http_ProbeConf">cloudprober.probes.http.ProbeConf</a>&gt; | dns_probe &lt;<a href="probes#cloudprober_probes_dns_ProbeConf">cloudprober.probes.dns.ProbeConf</a>&gt; | 
&nbsp;&nbsp;&nbsp;user_defined_probe &lt;string&gt;]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := Files.FindDescriptorByName(protoreflect.FullName(fldName))
			assert.NoError(t, err)
			fld := desc.(protoreflect.FieldDescriptor)

			done := map[string]bool{}
			assert.Equal(t, tt.want, fieldToToken(fld, tt.f, &done))

			// Make sure we get nil for the second field from same oneof.
			desc, err = Files.FindDescriptorByName(protoreflect.FullName(fldName2))
			assert.NoError(t, err)
			fld = desc.(protoreflect.FieldDescriptor)
			assert.Nil(t, fieldToToken(fld, tt.f, &done))

			assert.Equal(t, tt.want, formatOneOf(fld.ContainingOneof(), tt.f))
		})
	}
}

func TestFormatEnum(t *testing.T) {
	const fldName = "cloudprober.probes.ProbeDef.type"

	tests := []struct {
		name string
		f    Formatter
		want *Token
	}{
		{
			name: "default",
			want: &Token{
				Kind: "enum",
				Text: "type: (HTTP|TCP|EXTENSION|USER_DEFINED)",
			},
		},
		{
			name: "with-prefix",
			f: Formatter{
				prefix: "  ",
			},
			want: &Token{
				Kind:   "enum",
				Prefix: "  ",
				Text:   "type: (HTTP|TCP|EXTENSION|USER_DEFINED)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := Files.FindDescriptorByName(protoreflect.FullName(fldName))
			assert.NoError(t, err)

			fld := desc.(protoreflect.FieldDescriptor)

			done := map[string]bool{}
			assert.Equal(t, tt.want, fieldToToken(fld, tt.f, &done))

			assert.Equal(t, tt.want, formatEnum(fld.Enum(), "type", tt.f))
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
		f    Formatter
		want *Token
	}{
		{
			name: "simple",
			in: &Token{
				Kind:   "string",
				Prefix: "  ",
			},
			want: &Token{
				Kind:      "string",
				Prefix:    "  ",
				ExtraLine: "\n",
			},
		},
		{
			name: "with-url",
			in: &Token{
				Kind: "cloudprober.probes.ProbeDef",
			},
			want: &Token{
				Kind:      "cloudprober.probes.ProbeDef",
				URL:       "probes#cloudprober_probes_ProbeDef",
				ExtraLine: "\n",
			},
		},
		{
			name: "with-url-with-relpath",
			f:    Formatter{}.WithRelPath(".."),
			in: &Token{
				Kind: "cloudprober.probes.ProbeDef",
			},
			want: &Token{
				Kind:      "cloudprober.probes.ProbeDef",
				URL:       "../probes#cloudprober_probes_ProbeDef",
				ExtraLine: "\n",
			},
		},
		{
			name: "with-default",
			in: &Token{
				Kind:    "string",
				Default: "2s",
			},
			want: &Token{
				Kind:      "string",
				Suffix:    " | default: 2s",
				Default:   "2s",
				ExtraLine: "\n",
			},
		},
		{
			name: "header-not-yaml",
			in: &Token{
				Kind:          "string",
				MessageHeader: true,
				NoExtraLine:   true,
			},
			want: &Token{
				Kind:          "string",
				MessageHeader: true,
				Suffix:        " {",
				Sep:           " ",
				NoExtraLine:   true,
				ExtraLine:     "",
			},
		},
		{
			name: "header-yaml",
			in: &Token{
				Kind:          "string",
				MessageHeader: true,
				yaml:          true,
			},
			want: &Token{
				Kind:          "string",
				MessageHeader: true,
				yaml:          true,
				Suffix:        ":",
				Sep:           " ",
				ExtraLine:     "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want.Sep == "" {
				tt.want.Sep = ": "
			}
			assert.Equal(t, []*Token{tt.want}, ProcessTokensForHTML([]*Token{tt.in}, tt.f))
		})
	}
}

func TestKindToURL(t *testing.T) {
	tests := []struct {
		f    Formatter
		kind string
		want string
	}{
		{
			kind: "cloudprober.probes.ProbeDef.interval_msec",
			want: "probes#cloudprober_probes_ProbeDef_interval_msec",
		},
		{
			kind: "cloudprober.probes.ProbeDef.interval_msec",
			want: "../probes#cloudprober_probes_ProbeDef_interval_msec",
			f:    Formatter{}.WithRelPath(".."),
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
			assert.Equal(t, tt.want, kindToURL(tt.kind, tt.f))
		})
	}
}
