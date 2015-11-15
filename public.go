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
	</head>
	<body>
		<pre>go get -u {{.Package.Name}}</pre>
		<ul>
			{{range $mayor, $v := .Versions.Mayor}}
			<li>
				<span>{{$mayor}} -> {{$v.Name}}</span>
				<ul>
				{{range $parent.Versions.Match $mayor}}
				<li><a href="http://{{$package.Name.Base}}@{{.Name}}">{{.Name}}</a></li>
				{{end}}
				</ul>
			</li>
			{{end}}
    		</ul>
	</body>
</html>`
