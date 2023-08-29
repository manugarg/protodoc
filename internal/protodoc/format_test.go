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
	BuildFileDescRegistry(Files, "testdata", "github.com/manugarg/protodoc", nil)

	tests := []struct {
		name      string
		fldName   string
		f         Formatter
		nocomment bool
		want      *Token
	}{
		{
			name:      "no comment",
			fldName:   "cloudprober.probes.ProbeDef.interval_msec",
			nocomment: true,
			want: &Token{
				Kind:    "int32",
				Comment: "",
				Text:    "interval_msec",
			},
		},
		{
			name:    "not yaml",
			fldName: "cloudprober.probes.ProbeDef.interval_msec",
			f: Formatter{
				yaml:   false,
				prefix: "",
			},
			nocomment: false,
			want: &Token{
				Kind:    "int32",
				Comment: "# Interval between two probe runs in milliseconds.\n# Only one of \"interval\" and \"inteval_msec\" should be defined.\n# Default interval is 2s.\n#",
				Text:    "interval_msec",
			},
		},
		{
			name:    "yaml",
			fldName: "cloudprober.probes.ProbeDef.interval_msec",
			f: Formatter{
				yaml:   true,
				prefix: "",
			},
			nocomment: false,
			want: &Token{
				Kind:    "int32",
				Comment: "# Interval between two probe runs in milliseconds.\n# Only one of \"interval\" and \"inteval_msec\" should be defined.\n# Default interval is 2s.\n#",
				Text:    "intervalMsec",
			},
		},
		{
			name:    "yaml with prefix",
			fldName: "cloudprober.probes.ProbeDef.interval_msec",
			f: Formatter{
				yaml:   true,
				prefix: "  ",
			},
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
			desc, err := Files.FindDescriptorByName(protoreflect.FullName(tt.fldName))
			assert.NoError(t, err)

			fld := desc.(protoreflect.FieldDescriptor)
			assert.Equal(t, tt.want, finalToken(fld, tt.f, tt.nocomment))
		})
	}
}
