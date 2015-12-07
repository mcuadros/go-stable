package gopkg

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Public(ctx *gin.Context) {
	pkg, ok := ctx.Get("package")
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	f := NewFetcher(pkg.(*Package), nil)
	versions, err := f.Versions()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "package", gin.H{
		"Package":  pkg,
		"Versions": versions,
	})
}

func (s *Server) setTemplates() {
	t, _ := template.New("package").Parse(packageHTML)
	s.SetHTMLTemplate(t)
}

const packageHTML = `
{{$parent := .}}
{{$package := .Package}}
{{$current := .Versions.BestMatch .Package.Repository.Rev}}

<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Package.Name}}</title>
		<link rel="stylesheet" type="text/css" href="http://fonts.googleapis.com/css?family=Raleway" />
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/meyer-reset/2.0/reset.css" />
		<link rel="stylesheet" href="file:///Users/mcuadros/gopkg.css" />
		<style>
			* {
				text-decoration: none;
				color: black;
				font-family: Raleway;
			}

			body {
				width: 60%;
				margin: auto;
				padding-top: 10%;
			}

			pre {
				padding: 10px;
				font-family: "Courier New", Courier, monospace;
				border: 1px solid black;
				text-align: center;
			}

			ul.mayor {
				padding: auto;
			}
			
			ul.mayor > li {
				padding: 10px;
				float: left;
				display: block;
			}

			ul.minor > li {
				margin-left: 40px;

			}
			.mayor .minor { text-decoration: none; }
			.mayor .best { font-weight: bold; }

		</style>
	</head>
	<body>
		<pre>go get -u {{.Package.Name}}</pre>
		<ul class="mayor">
			{{range $mayor, $v := .Versions.Mayor}}
			<li>
				{{$mayor}} ->
				<ul class="minor">
				{{range $parent.Versions.Match $mayor}}
				<li><a 
					href="https://{{$package.Name.Change .}}"
					class="{{if eq .Name $v.Name}}best{{else}}other{{end}}"
				>{{.Name}}</a></li>
				{{end}}
				</ul>
			</li>
			{{end}}
    		</ul>
	</body>
</html>`
