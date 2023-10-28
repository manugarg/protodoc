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
	"flag"
	"fmt"
	"html/template"
	"path"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var homeURL = flag.String("home_url", "", "Home URL for the documentation.")

func formatComment(fld protoreflect.Descriptor, f Formatter) string {
	ff, err := Files.FindDescriptorByName(fld.FullName())
	if err != nil {
		panic(err)
	}
	wf, err := desc.WrapDescriptor(ff)
	if err != nil {
		panic(err)
	}

	comment := wf.GetSourceInfo().GetLeadingComments()
	if comment != "" && strings.TrimSpace(comment) != "" {
		var temp []string
		lines := strings.Split(comment, "\n")
		for i, line := range lines {
			if i == len(lines)-1 && strings.TrimSpace(line) == "" {
				continue
			}
			temp = append(temp, f.prefix+"#"+line)
		}
		comment = strings.Join(temp, "\n")
	}
	return comment
}

func formatOneOf(ood protoreflect.OneofDescriptor, f Formatter) *Token {
	oof := ood.Fields()
	oneofFields := []string{}

	for i := 0; i < oof.Len(); i++ {
		fld := oof.Get(i)

		if fldEnum := fld.Enum(); fldEnum != nil {
			// We get token text from finalToken and kind string from enum
			tok := finalToken(fld, f, true)
			oneofFields = append(oneofFields, formatEnum(fldEnum, tok.Text, f).Text)
			continue
		}

		tok := finalToken(oof.Get(i), f, true)
		s := fmt.Sprintf("%s &lt;%s&gt;", tok.Text, tok.Kind)
		if strings.HasPrefix(tok.Kind, "cloudprober.") {
			s = fmt.Sprintf("%s &lt;<a href=\"%s\">%s</a>&gt;", tok.Text, kindToURL(tok.Kind, f), tok.Kind)
		}
		oneofFields = append(oneofFields, s)
	}

	text := "["
	for i, tok := range oneofFields {
		if i != 0 && i%2 == 0 {
			text += "\n" + strings.ReplaceAll(f.prefix+" ", " ", "&nbsp;")
		}
		if i == len(oneofFields)-1 {
			text += tok + "]"
			break
		}
		text += tok + " | "
	}
	return &Token{
		Comment:  formatComment(ood, f),
		Kind:     "oneof",
		Prefix:   f.prefix,
		TextHTML: template.HTML(text),
	}
}

func formatEnum(ed protoreflect.EnumDescriptor, name string, f Formatter) *Token {
	enumVals := []string{}
	for i := 0; i < ed.Values().Len(); i++ {
		enumVals = append(enumVals, string(ed.Values().Get(i).Name()))
	}
	return &Token{
		Comment: formatComment(ed, f),
		Kind:    "enum",
		Prefix:  f.prefix,
		Text:    fmt.Sprintf("%s: (%s)", name, strings.Join(enumVals, "|")),
	}
}

func fieldToToken(fld protoreflect.FieldDescriptor, f Formatter, done *map[string]bool) *Token {
	if oo := fld.ContainingOneof(); oo != nil {
		// In proto3, optional fields have an oneof container.
		if oo.Fields().Len() == 1 {
			return finalToken(fld, f, false)
		}
		if (*done)[string(oo.Name())] {
			return nil
		}
		(*done)[string(oo.Name())] = true
		return formatOneOf(oo, f)
	}

	if ed := fld.Enum(); ed != nil {
		name := string(fld.Name())
		if f.yaml && f.jsonNamesForYAML {
			name = fld.JSONName()
		}
		tok := formatEnum(ed, name, f)
		tok.Comment = formatComment(fld, f)
		return tok
	}

	return finalToken(fld, f, false)
}

func ProcessTokensForHTML(toks []*Token, f Formatter) []*Token {
	for _, tok := range toks {
		tok.URL = kindToURL(tok.Kind, f)

		if tok.MessageHeader {
			tok.Suffix = " {"
			if tok.yaml {
				tok.Suffix = ":"
			}
			tok.Sep = " "
		} else {
			if tok.Default != "" {
				tok.Suffix = template.HTML(" | default: " + tok.Default)
			}
			tok.Sep = ": "
		}

		if tok.TextHTML == "" {
			tok.TextHTML = template.HTML(template.HTMLEscapeString(tok.Text))
		}

		tok.ExtraLine = "\n"
		if tok.NoExtraLine {
			tok.ExtraLine = ""
		}
	}
	return toks
}

func kindToURL(kind string, f Formatter) string {
	if !strings.HasPrefix(kind, "cloudprober.") {
		return ""
	}
	parts := strings.SplitN(kind, ".", 3)
	if len(parts) > 2 {
		kindForURL := strings.ReplaceAll(kind, ".", "_")
		return path.Join(*homeURL, f.relPath, parts[1]+"#"+kindForURL)
	}
	return ""
}
