🎉 {{.Name}} `{{.Version}}` is here!

## Contents

Name | Version | Previous | Description
-----|---------|----------|------------
{{- range .Registry }}
[{{.Repo.Name}}]({{.Repo.URL}}) | {{ template "version" (dict "Repo" .Repo "Version" (index .Versions 0)) }} | {{if eq (index .Versions 0) (index .Versions 1)}}≡{{else}}{{if empty (index .Versions 1)}}✗{{else}}{{ template "version" (dict "Repo" .Repo "Version" (index .Versions 1)) }}{{end}}{{end}} | {{ .Description }}
{{- end }}

≡ means same, ✗ means missing

<!-- Please don't remove the following footer, k6dist needs it! -->
{{.Footer}}

{{- define "version" -}}
{{if hasPrefix "https://github.com" .Repo.URL }}[{{.Version}}](https://github.com/{{.Repo.Owner}}/{{.Repo.Name}}/releases/tag/{{.Version}}){{else}}{{.Version}}{{end}}
{{- end -}}