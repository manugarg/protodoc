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
	"strings"

	"github.com/cloudprober/cloudprober/logger"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var Files = &protoregistry.Files{}

type Token struct {
	Prefix  string
	Suffix  template.HTML
	Comment string
	Kind    string
	Text    string
	URL     string
	Default string

	MessageHeader bool
	yaml          bool
	NoExtraLine   bool

	// Filed by token processor
	TextHTML  template.HTML
	Sep       string
	ExtraLine string
}

type Formatter struct {
	yaml    bool
	depth   int
	prefix  string
	relPath string
}

func (f Formatter) WithYAML(yaml bool) Formatter {
	f2 := f
	f2.yaml = yaml
	return f2
}

func (f Formatter) WithDepth(depth int) Formatter {
	f2 := f
	f2.depth = depth
	return f2
}

func (f Formatter) WithPrefix(prefix string) Formatter {
	f2 := f
	f2.prefix = prefix
	return f2
}

func (f Formatter) WithRelPath(relPath string) Formatter {
	f2 := f
	f2.relPath = relPath
	return f2
}

func finalToken(fld protoreflect.FieldDescriptor, f Formatter, nocomment bool) *Token {
	var comment string
	if !nocomment {
		comment = formatComment(fld, f)
	}

	kind := fld.Kind().String()
	if fld.Kind() == protoreflect.MessageKind {
		kind = string(fld.Message().FullName())
	}

	tok := &Token{
		yaml:    f.yaml,
		Prefix:  f.prefix,
		Comment: comment,
		Kind:    kind,
		Text:    string(fld.Name()),
	}

	if fld.HasDefault() {
		tok.Default = fld.Default().String()
	}

	if f.yaml {
		tok.Text = fld.JSONName()
	}

	return tok
}

func dumpExtendedMsg(fld protoreflect.FieldDescriptor, f Formatter) ([]*Token, []protoreflect.FullName) {
	var nextMessageName []protoreflect.FullName
	var lines []*Token

	nextMessageName = append(nextMessageName, fld.Message().FullName())
	tok := finalToken(fld, f, false)
	tok.MessageHeader = true
	tok.NoExtraLine = true

	lines = append(lines, tok)

	newPrefix := f.prefix + "  "
	if fld.Cardinality() == protoreflect.Repeated && f.yaml {
		newPrefix = f.prefix + "    "
	}
	toks, next := DumpMessage(fld.Message(), f.WithDepth(f.depth-1).WithPrefix(newPrefix))
	if f.yaml && fld.Cardinality() == protoreflect.Repeated {
		toks[0].Prefix = f.prefix + "  - "
	}
	lines = append(lines, toks...)

	// If it's not a yaml, add a "}" at the end and limit the line break before
	// that to just one (default is 2).
	if !f.yaml {
		lines[len(lines)-1].NoExtraLine = true
		lines = append(lines, &Token{Prefix: f.prefix, Text: "}"})
	}
	nextMessageName = append(nextMessageName, next...)

	return lines, nextMessageName
}

func DumpMessage(md protoreflect.MessageDescriptor, f Formatter) ([]*Token, []protoreflect.FullName) {
	var nextMessageName []protoreflect.FullName

	var lines []*Token

	// We use this to catch duplication for oneof and enum fields.
	done := map[string]bool{}

	for i := 0; i < md.Fields().Len(); i++ {
		fld := md.Fields().Get(i)

		if fld.Kind() == protoreflect.MessageKind && f.depth > 1 {
			toks, next := dumpExtendedMsg(fld, f)
			lines = append(lines, toks...)
			nextMessageName = append(nextMessageName, next...)
		} else {
			if tok := fieldToToken(md.Fields().Get(i), f, &done); tok != nil {
				lines = append(lines, tok)
			}
			if fld.Kind() == protoreflect.MessageKind {
				nextMessageName = append(nextMessageName, fld.Message().FullName())
			}
		}
	}

	return lines, nextMessageName
}

func ArrangeIntoPackages(paths []string, l *logger.Logger) map[string][]string {
	packages := make(map[string][]string)
	for _, path := range paths {
		parts := strings.SplitN(path, ".", 3)
		if len(parts) < 3 {
			l.Warningf("Skipping %s, not enough parts in package", path)
			continue
		}
		if parts[0] != "cloudprober" {
			l.Warningf("Skipping %s, not a cloudprober package", path)
			continue
		}
		packages[parts[1]] = append(packages[parts[1]], path)
	}
	return packages
}
