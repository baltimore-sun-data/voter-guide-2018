package main

import "html/template"

var templates = template.Must(template.New("name").Parse(`

{{ define "metadata" }}
    <ul>
    {{ range .Contests }}
        <li>{{ .Name }}</li>
    {{ end }}
    </ul>
{{ end }}

`))
