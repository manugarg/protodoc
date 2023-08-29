// Copyright 2023 The Cloudprober Authors.
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
	//	homeURL      = flag.String("home_url", "", "Home URL for the documentation.")
	outFmt        = flag.String("format", "yaml", "textpb or yaml")
	outDir        = flag.String("out_dir", "docs", "Output directory for the documentation.")
	protoRootDir  = flag.String("proto_root_dir", ".", "Root directory for the proto files.")
	packagePrefix = flag.String("package_prefix", "", "Package prefix to resolve import paths")
)

func main() {
	flag.Parse()

	l := &logger.Logger{}

	protodoc.BuildFileDescRegistry(protodoc.Files, *protoRootDir, *packagePrefix, l)

	// Top level message
	m, err := protodoc.Files.FindDescriptorByName("cloudprober.ProberConfig")
	if err != nil {
		panic(err)
	}

	f := protodoc.Formatter{}.WithYAML(*outFmt == "yaml")

	toks, nextMessageNames := protodoc.DumpMessage(m.(protoreflect.MessageDescriptor), f.WithDepth(2))
	outF, err := os.OpenFile(filepath.Join(*outDir, "index.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		l.Criticalf("Error opening output file: %v", err)
	}
	template.Must(template.New("index").Parse(protodoc.IndexTmpl)).Execute(outF, protodoc.ProcessTokensForHTML(toks))

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
	type msgTokens struct {
		Name   string
		Tokens []*protodoc.Token
	}

	for pkg, msgs := range packages {
		sort.Strings(msgs)
		toks := []*msgTokens{}
		for _, msg := range msgs {
			toks = append(toks, &msgTokens{Name: msg, Tokens: protodoc.ProcessTokensForHTML(msgToDoc[msg])})
		}
		outF, err := os.OpenFile(filepath.Join(*outDir, pkg+".html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			l.Criticalf("Error opening output file: %v", err)
		}
		template.Must(template.New("package").Parse(protodoc.PackageTmpl)).Execute(outF, toks)
	}

}
