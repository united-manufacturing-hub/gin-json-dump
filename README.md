# Gin JSON dump

This middleware dumps the request and response in JSON format.

## Installation

```bash
go get github.com/united-manufacturing-hub/gin-json-dump
```

## Usage

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/united-manufacturing-hub/gin-json-dump"
)

func main() {
    r := gin.Default()
    r.Use(ginjsondump.Dump())
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello World!",
        })
    })
    r.Run()
}
```

### Thanks

This middleware is inspired by [gin-dump](https://github.com/tpkeeper/gin-dump)
