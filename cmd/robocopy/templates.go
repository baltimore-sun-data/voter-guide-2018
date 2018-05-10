package main

import "html/template"

var templates = template.Must(template.New("name").Parse(`

{{ define "contest" }}
<table>
    <tr>
        <td>LastUpdate</td>
        <td>{{ .LastUpdate }}</td>
    </tr>
    <tr>
        <td>Contest</td>
        <td>{{ .Contest }}</td>
    </tr>
    <tr>
        <td>District</td>
        <td>{{ .District }}</td>
    </tr>
    <tr>
        <td>Jurisdiction</td>
        <td>{{ .Jurisdiction }}</td>
    </tr>
    <tr>
        <td>Party</td>
        <td>{{ .Party }}</td>
    </tr>
    <tr>
        <td>VoteFor</td>
        <td>{{ .VoteFor }}</td>
    </tr>
    <tr>
        <td>PrimaryDescription</td>
        <td>{{ .PrimaryDescription }}</td>
    </tr>
    <tr>
        <td>SecondaryDescription</td>
        <td>{{ .SecondaryDescription }}</td>
    </tr>
    <tr>
        <td>FullDescription</td>
        <td>{{ .FullDescription }}</td>
    </tr>
</table>
<ul>
    {{ range .Options }}
    <li>
        <table>
            <tr>
                <td>Text</td>
                <td>{{ .Text }}</td>
            </tr>
            <tr>
                <td>TotalVotes</td>
                <td>{{ .TotalVotes }}</td>
            </tr>
            <tr>
                <td>Jurisdiction</td>
                <td>{{ .Jurisdiction }}</td>
            </tr>
            <tr>
                <td>District</td>
                <td>{{ .District }}</td>
            </tr>
        </table>
        SubResults
        <ul>
        {{ range .SubResults }}
        <li>
        <table>
            <tr>
                <td>Jurisdiction</td>
                <td>{{ .Jurisdiction }}</td>
            </tr>
            <tr>
                <td>District</td>
                <td>{{ .District }}</td>
            </tr>
            <tr>
                <td>TotalVotes</td>
                <td>{{ .TotalVotes }}</td>
            </tr>
        </table>
        </li>
        {{ end }}
        </ul>
    </li>
    {{ end }}
</ul>
{{ end }}

`))
