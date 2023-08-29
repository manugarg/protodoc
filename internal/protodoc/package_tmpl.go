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

var PackageTmpl = `
<html>
<head>
    <title>Cloudprober Configuration</title>
    <style>
        body {
            font-family: monospace;
        }
        .comment {
            color: #888;
            font-family: monospace;
            white-space: pre;
        }
    </style>
</head>
<body>
{{range .}}
<h2 id="{{ .Name }}">{{ .Name }}</h2>
    {{range .Tokens}}
        {{- if .Comment}}
            <div class="comment">{{.Comment}}</div>
        {{- end -}}
        {{- if .URL}}
            {{- .PrefixHTML}}{{- .TextHTML}}: <<a href="{{.URL}}">{{- .Kind}}</a>>{{- .Suffix }}
        {{- else if .Kind}}
            {{- .PrefixHTML}}{{- .TextHTML}}: <{{- .Kind}}>{{- .Suffix }}
        {{- else}}
            {{- .PrefixHTML}}{{- .TextHTML}}{{- .Suffix }}
        {{- end}}
    {{end}}
{{end}}
</body>
</html>`
