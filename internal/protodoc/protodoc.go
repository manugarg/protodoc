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

	PrefixHTML template.HTML
	TextHTML   template.HTML
}

type Formatter struct {
	yaml   bool
	depth  int
	prefix string
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
	tok.Suffix = " {"
	if f.yaml {
		tok.Suffix = ":"
	}
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
		lines[len(lines)-1].Suffix = "<br>"
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

func ProcessTokensForHTML(toks []*Token) []*Token {
	for _, tok := range toks {
		tok.PrefixHTML = template.HTML(strings.ReplaceAll(tok.Prefix, " ", "&nbsp;"))

		tok.URL = kindToURL(tok.Kind)

		if tok.Suffix == "" {
			tok.Suffix = template.HTML("<br><br>")
			if tok.Default != "" {
				tok.Suffix = template.HTML(" | default: " + tok.Default + "<br><br>")
			}
		} else {
			if !strings.HasSuffix(string(tok.Suffix), "<br>") {
				tok.Suffix = template.HTML(tok.Suffix + "<br>")
			}
		}

		if tok.TextHTML == "" {
			tok.TextHTML = template.HTML(template.HTMLEscapeString(tok.Text))
		}
	}
	return toks
}
