package ginjsondump

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"net/url"
)

type Options struct {
	// Dump request body
	ShowReqBody bool
	// Dump response body
	ShowRespBody bool
	// Dump request headers
	ShowReqHeaders bool
	// Dump response headers
	ShowRespHeaders bool
	// Dump response status
	ShowRespStatus bool
	// DumpFunc is a function that will be called with the dump data
	DumpFunc func(dumpStr DumpData)
	// Dump request path
	ShowReqPath bool
	// Dump request method
	ShowReqMethod bool
	// Dump request query
	ShowReqQuery bool
}

func Dump() gin.HandlerFunc {
	return DumpWithOptions(Options{
		ShowReqBody:     true,
		ShowRespBody:    true,
		ShowReqHeaders:  true,
		ShowRespHeaders: true,
		ShowRespStatus:  true,
		DumpFunc: func(dumpData DumpData) {
			dump, err := jsoniter.MarshalIndent(dumpData, "", "  ")
			if err != nil {
				return
			}
			fmt.Println(string(dump))
		},
		ShowReqPath:   true,
		ShowReqMethod: true,
		ShowReqQuery:  true,
		Formatted:     true,
	})
}

type DumpData struct {
	Request  DumpRequest
	Response DumpResponse
}

func (d *DumpData) ToJSONString() string {
	marshal, err := jsoniter.MarshalToString(d)
	if err != nil {
		return err.Error()
	}
	return marshal
}
func (d *DumpData) ToJSONFormattedString() string {
	marshal, err := jsoniter.MarshalIndent(d, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(marshal)
}

type DumpRequest struct {
	Headers http.Header
	Body    map[string]interface{}
	Path    string
	Method  string
	Query   url.Values
}

type DumpResponse struct {
	Headers http.Header
	Body    map[string]interface{}
	Status  int
}

func DumpWithOptions(options Options) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dumpData DumpData

		handleRequest(ctx, &dumpData, options)

		// Execute the next handler
		ctx.Writer = &bodyWriter{ctx.Writer, bytes.NewBufferString("")}
		ctx.Next()

		handleResponse(ctx, &dumpData, options)

		// Dump the request and response details
		if options.DumpFunc != nil {
			options.DumpFunc(dumpData)
		}
	}
}

func handleRequest(ctx *gin.Context, d *DumpData, options Options) {
	if options.ShowReqHeaders {
		d.Request.Headers = ctx.Request.Header
	}

	if options.ShowReqBody {
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			d.Request.Body = map[string]interface{}{
				"error": err.Error(),
			}
		} else {
			// try to parse json
			var body map[string]interface{}
			err = jsoniter.Unmarshal(buf, &body)
			if err != nil {
				d.Request.Body = map[string]interface{}{
					"value": string(buf),
				}
			} else {
				d.Request.Body = body
			}
		}
	}

	if options.ShowRespStatus {
		d.Response.Status = ctx.Writer.Status()
	}

	if options.ShowReqPath {
		d.Request.Path = ctx.Request.URL.Path
	}

	if options.ShowReqMethod {
		d.Request.Method = ctx.Request.Method
	}

	if options.ShowReqQuery {
		d.Request.Query = ctx.Request.URL.Query()
	}
}

func handleResponse(ctx *gin.Context, d *DumpData, options Options) {
	if options.ShowRespHeaders {
		d.Response.Headers = ctx.Writer.Header()
	}
	if options.ShowRespBody && bodyAllowedForStatus(ctx.Writer.Status()) {
		bw, ok := ctx.Writer.(*bodyWriter)
		if !ok {
			d.Response.Body = map[string]interface{}{
				"error": "response body is not a bodyWriter",
			}
		} else {
			// try to parse json
			var body map[string]interface{}
			err := jsoniter.Unmarshal(bw.bodyCache.Bytes(), &body)
			if err != nil {
				d.Response.Body = map[string]interface{}{
					"value": string(bw.bodyCache.Bytes()),
				}
			} else {
				d.Response.Body = body
			}
		}

	}
}

type bodyWriter struct {
	gin.ResponseWriter
	bodyCache *bytes.Buffer
}

// rewrite Write()
func (w bodyWriter) Write(b []byte) (int, error) {
	w.bodyCache.Write(b)
	return w.ResponseWriter.Write(b)
}

// bodyAllowedForStatus reports whether a given response status code
// permits a body. See RFC 7230, section 3.3.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}
