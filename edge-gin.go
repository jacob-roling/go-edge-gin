package edgegin

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/jrolingdev/go-edge"
)

type GinEdgeRender struct {
	Edge     edge.Edge
	Context  any
	Template edge.Template
}

// Default creates a GinEdgeRender instance with default options.
func Default() *GinEdgeRender {
	return &GinEdgeRender{
		Edge: edge.Default(),
	}
}

func New(config *edge.Config) *GinEdgeRender {
	return &GinEdgeRender{
		Edge: edge.New(*config),
	}
}

func (r GinEdgeRender) Instance(name string, data any) render.Render {
	var template edge.Template
	filename := path.Join(r.Edge.BaseDirectory, name+".edge")

	// Always read template files from disk if in debug mode, use cache otherwise.
	if gin.Mode() == "debug" {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		template = r.Edge.Compile(string(bytes))
	} else {
		template = r.Edge.Cache[filename]
	}

	r.Template = template

	return r
}

// Render should write the content type, then render the template to the response.
func (r GinEdgeRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	// Make sure we return the error from Template.Exec if it fails, before calling w.Write
	output, err := edge.Exec(r.Template, r.Context)

	if err != nil {
		return err
	}

	// Don't care about status but return err if Write fails
	_, err = w.Write([]byte(output))
	return err
}

// WriteContentType writes header information about content RaymondRender outputs.
// This will now implement gin's render.Render interface.
func (r GinEdgeRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
	}
}
