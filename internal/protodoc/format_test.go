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
				Comment: "# Interval between two probe runs in milliseconds.\n# Only one of \"interval\" and \"inteval_msec\" should be defined.\n# Default interval is 2s.\n#",
				Text:    "interval_msec",
			},
		},
		{
			name:      "yaml",
			f:         Formatter{}.WithYAML(true),
			nocomment: false,
			want: &Token{
				Kind:    "int32",
				Comment: "# Interval between two probe runs in milliseconds.\n# Only one of \"interval\" and \"inteval_msec\" should be defined.\n# Default interval is 2s.\n#",
				Text:    "intervalMsec",
			},
		},
		{
			name:      "yaml with prefix",
			f:         Formatter{}.WithYAML(true).WithPrefix("  "),
			nocomment: false,
			want: &Token{
				Kind:    "int32",
				Prefix:  "  ",
				Comment: "  # Interval between two probe runs in milliseconds.\n  # Only one of \"interval\" and \"inteval_msec\" should be defined.\n  # Default interval is 2s.\n  #",
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

	tests := []struct {
		name string
		f    Formatter
		want *Token
	}{
		{
			name: "default",
			want: &Token{
				Kind: "oneof",
				TextHTML: `[http_probe &lt;<a href="probes.html#cloudprober.probes.http.ProbeConf">cloudprober.probes.http.ProbeConf</a>&gt; | dns_probe &lt;<a href="probes.html#cloudprober.probes.dns.ProbeConf">cloudprober.probes.dns.ProbeConf</a>&gt; | <br>
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
				TextHTML: `[httpProbe &lt;<a href="probes.html#cloudprober.probes.http.ProbeConf">cloudprober.probes.http.ProbeConf</a>&gt; | dnsProbe &lt;<a href="probes.html#cloudprober.probes.dns.ProbeConf">cloudprober.probes.dns.ProbeConf</a>&gt; | <br>
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
				TextHTML: `[http_probe &lt;<a href="probes.html#cloudprober.probes.http.ProbeConf">cloudprober.probes.http.ProbeConf</a>&gt; | dns_probe &lt;<a href="probes.html#cloudprober.probes.dns.ProbeConf">cloudprober.probes.dns.ProbeConf</a>&gt; | <br>
&nbsp;&nbsp;&nbsp;user_defined_probe &lt;string&gt;]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := Files.FindDescriptorByName(protoreflect.FullName(fldName))
			assert.NoError(t, err)

			oofd := desc.(protoreflect.FieldDescriptor).ContainingOneof()

			assert.Equal(t, tt.want, formatOneOf(oofd, tt.f))
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

			ed := desc.(protoreflect.FieldDescriptor).Enum()

			assert.Equal(t, tt.want, formatEnum(ed, "type", tt.f))
		})
	}
}
