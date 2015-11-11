package gopkg

import (
	"fmt"
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

	fmt.Println(pkg, ok)
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
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Package.Name}}</title>
	</head>
	<body>
		<pre>go get -u {{.Package.Name}}</pre>
    <ul>
    {{range .Versions}}
      <li>{{.Name}} {{.Versions.Match .Package.Repository.Rev}}</li>
    {{end}}
    </ul>
	</body>
</html>`
