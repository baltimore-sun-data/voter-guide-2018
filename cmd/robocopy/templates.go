package main

import "html/template"

var templates = template.Must(template.New("name").Parse(`

{{ define "metadata" }}
    <table>
    {{ $m := . }}
    {{ range $c := .Contests }}
        <tr class="contest">
            <td class="jurisdiction">{{ $c.Jurisdiction $m }}</td>
            <td class="name">{{ $c.Name }}</td>
            <td class="district">{{ $c.DistrictName $m }}</td>
            <td class="party">{{ $c.Party $m }}</td>
            <td>Pick {{ $c.VoteFor }}</td>
        </tr>
    {{ end }}
    </ul>
{{ end }}

`))
