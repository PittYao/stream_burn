package gin

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
