package util

import (
	"fmt"

	"github.com/kataras/iris/v12"
)

// Response ...
func Response(data interface{}, status int, c iris.Context) {
	c.StatusCode(status)
	c.JSON(iris.Map{"data": data})
}

// ResponseError ...
func ResponseError(data error, status int, c iris.Context) {
	c.StatusCode(status)
	c.JSON(c.Tr(fmt.Sprintf("%+v", data)))
}
