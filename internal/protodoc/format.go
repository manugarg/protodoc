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
	"fmt"
	"html/template"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func formatComment(fld protoreflect.FieldDescriptor, f Formatter) string {
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
		for _, line := range strings.Split(comment, "\n") {
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
		tok := finalToken(oof.Get(i), f, true)
		s := fmt.Sprintf("%s &lt;%s&gt;", tok.Text, tok.Kind)
		if strings.HasPrefix(tok.Kind, "cloudprober.") {
			s = fmt.Sprintf("%s &lt;<a href=\"%s\">%s</a>&gt;", tok.Text, kindToURL(tok.Kind), tok.Kind)
		}
		oneofFields = append(oneofFields, s)
	}

	text := "["
	for i, tok := range oneofFields {
		if i != 0 && i%2 == 0 {
			text += "<br>\n" + strings.ReplaceAll(f.prefix+" ", " ", "&nbsp;")
		}
		if i == len(oneofFields)-1 {
			text += tok + "]"
			break
		}
		text += tok + " | "
	}
	return &Token{
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
		Kind:   "enum",
		Prefix: f.prefix,
		Text:   fmt.Sprintf("%s: (%s)", name, strings.Join(enumVals, "|")),
	}
}

func fieldToToken(fld protoreflect.FieldDescriptor, f Formatter, done *map[string]bool) *Token {
	ed := fld.Enum()
	if ed != nil {
		name := string(fld.Name())
		if f.yaml {
			name = fld.JSONName()
		}
		return formatEnum(ed, name, f)
	}

	oo := fld.ContainingOneof()
	if oo != nil {
		if (*done)[string(oo.Name())] {
			return nil
		}
		(*done)[string(oo.Name())] = true
		return formatOneOf(oo, f)
	}

	return finalToken(fld, f, false)
}