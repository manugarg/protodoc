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

package main

import (
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"sort"

	"github.com/cloudprober/cloudprober/logger"
	"github.com/manugarg/protodoc/internal/protodoc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	outFmt        = flag.String("format", "yaml", "textpb or yaml")
	outDir        = flag.String("out_dir", "proto_docs", "Output directory for the documentation.")
	protoRootDir  = flag.String("proto_root_dir", ".", "Root directory for the proto files.")
	packagePrefix = flag.String("package_prefix", "", "Package prefix to resolve import paths")
)

type msgTokens struct {
	Name   string
	Tokens []*protodoc.Token
}

var docTmpl = template.Must(template.New("index").Parse(protodoc.DocTmpl))

func writeDoc(pkg string, mTokens []*msgTokens, l *logger.Logger) {
	outF, err := os.OpenFile(filepath.Join(*outDir, pkg+".html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		l.Criticalf("Error opening output file: %v", err)
	}
	defer outF.Close()

	if err := docTmpl.Execute(outF, mTokens); err != nil {
		l.Criticalf("Error executing template: %v", err)
	}
}

func main() {
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	l := &logger.Logger{}

	protodoc.BuildFileDescRegistry(protodoc.Files, *protoRootDir, *packagePrefix, l)

	// Top level message
	m, err := protodoc.Files.FindDescriptorByName("cloudprober.ProberConfig")
	if err != nil {
		panic(err)
	}

	f := protodoc.Formatter{}.WithYAML(*outFmt == "yaml")

	toks, nextMessageNames := protodoc.DumpMessage(m.(protoreflect.MessageDescriptor), f.WithDepth(2))
	mTokens := &msgTokens{Name: "", Tokens: protodoc.ProcessTokensForHTML(toks)}
	writeDoc("index", []*msgTokens{mTokens}, l)

	msgToDoc := map[string][]*protodoc.Token{}

	for len(nextMessageNames) > 0 {
		var nextLoop []protoreflect.FullName
		for _, msgName := range nextMessageNames {
			m, err := protodoc.Files.FindDescriptorByName(protoreflect.FullName(msgName))
			if err != nil {
				panic(err)
			}

			toks, next := protodoc.DumpMessage(m.(protoreflect.MessageDescriptor), f.WithDepth(1).WithPrefix("  "))
			msgToDoc[string(msgName)] = toks
			nextLoop = append(nextLoop, next...)
		}
		nextMessageNames = nextLoop
	}

	var msgs []string
	for key := range msgToDoc {
		msgs = append(msgs, key)
	}

	packages := protodoc.ArrangeIntoPackages(msgs, l)

	for pkg, msgs := range packages {
		sort.Strings(msgs)
		mtoks := []*msgTokens{}
		for _, msg := range msgs {
			mtoks = append(mtoks, &msgTokens{Name: msg, Tokens: protodoc.ProcessTokensForHTML(msgToDoc[msg])})
		}
		writeDoc(pkg, mtoks, l)
	}

	l.Infof("Documentation generated in %s", *outDir)
}
