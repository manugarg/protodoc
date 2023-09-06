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

var DocTmpl = `
<style>
.comment {
    color: #888;
}
.protodoc {
    border: 1px solid #ddd;
    border-left: 3px solid #e6522c;
    border-radius: 0;
    padding-left: 10px;
}
</style>
{{- range . -}}
{{- if .Name -}}<h3 id="{{ .Name | replace "." "_" }}">{{ .Name }} <a class="anchor" href="#{{ .Name | replace "." "_" }}">#</a></h3>{{- end }}
<pre class="protodoc">

{{ range .Tokens -}}
  {{- if .Comment }}<div class="comment">{{.Comment}}</div>{{ end -}}
  {{- if .URL }}
    {{- .Prefix}}{{.TextHTML}}{{.Sep}}<<a href="{{.URL}}">{{- .Kind}}</a>>{{.Suffix}}
  {{- else if .Kind }}
    {{- .Prefix}}{{.TextHTML}}{{.Sep}}<{{.Kind}}>{{.Suffix}}
  {{- else }}
    {{- .Prefix}}{{.TextHTML}}{{.Suffix}}
  {{- end }}
  {{- .ExtraLine }}
{{ end -}}
</pre>
{{- end -}}
`
