package test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	ginjsondump "github.com/united-manufacturing-hub/gin-json-dump/pkg"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// performRequest https://github.com/tpkeeper/gin-dump/blob/434e893eb6032f4932ccf61049b654e1fba1f78c/gindump_test.go#L21-L27
func performRequest(r http.Handler, method, contentType string, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Test_Matrix(t *testing.T) {
	tests := []struct {
		method         string
		contentType    string
		path           string
		requestBody    params
		expectedStatus int
		responseBody   string
	}{
		{"GET", "application/json", "/dump", params{Key: "hello", Value: "world"}, http.StatusOK, `{"data":"gin-dump","ok":true}`},
		{"POST", "application/json", "/dump", params{Key: "hello", Value: "world"}, http.StatusOK, `{"data":"gin-dump","ok":true}`},
		{"PUT", "application/json", "/dump", params{Key: "hello", Value: "world"}, http.StatusOK, `{"data":"gin-dump","ok":true}`},
		{"DELETE", "application/json", "/dump", params{Key: "hello", Value: "world"}, http.StatusOK, `{"data":"gin-dump","ok":true}`},
		{"PATCH", "application/json", "/dump", params{Key: "hello", Value: "world"}, http.StatusOK, `{"data":"gin-dump","ok":true}`},
		{"HEAD", "application/json", "/dump", params{Key: "hello", Value: "world"}, http.StatusOK, `{"data":"gin-dump","ok":true}`},
		{"OPTIONS", "application/json", "/dump", params{Key: "hello", Value: "world"}, http.StatusOK, `{"data":"gin-dump","ok":true}`},
		{"GET", "application/json", "/dumpX", params{Key: "hello", Value: "world"}, http.StatusNotFound, `404 page not found`},
	}

	router := gin.New()
	router.Use(ginjsondump.Dump())

	router.GET("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})
	router.POST("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})
	router.PUT("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})
	router.DELETE("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})
	router.PATCH("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})
	router.HEAD("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})
	router.OPTIONS("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})

	for _, test := range tests {
		b, err := jsoniter.Marshal(test.requestBody)
		if err != nil {
			t.Error(err)
			return
		}

		body := bytes.NewBuffer(b)
		w := performRequest(router, test.method, test.contentType, test.path, body)
		if w.Code != test.expectedStatus {
			t.Errorf("expected status %d, got %d", test.expectedStatus, w.Code)
		}
		if w.Body.String() != test.responseBody {
			t.Errorf("expected body %s, got %s", test.responseBody, w.Body.String())
		}
	}

}
