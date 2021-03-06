package main

import (
	"fmt"

	"github.com/suhanyujie/go-utils/libs/framwork/base3/core/engine"
)

var (
	Port = 3001
)

func main() {
	r := engine.New()

	r.Get("/test", func(c *engine.Context) {
		c.Json(200, "who are you....")
	})
	r.Get("/query/:keyword", func(c *engine.Context) {
		c.Json(200, fmt.Sprintf("the query param is: %s\n", c.GetParam("keyword")))
	})
	v2Group := r.Group("v2")
	v2Group.Get("/projectList", func(c *engine.Context) {
		c.Json(200, fmt.Sprintf("the route path is: /v2/projectList..."))
	})

	// 静态文件路由
	r.Static("/assets", "D:\\repo\\self\\golang\\ThinkGo")

	// 模板文件
	r.LoadHtmlGlob("examples/templates/*")
	r.Get("/html", func(c *engine.Context) {
		c.Html(200, "index.html", map[string]string{
			"Name": "Samuel",
		})
	})

	fmt.Printf("server is start in port: %d\n", Port)
	addr := fmt.Sprintf(":%d", Port)
	r.Run(addr)
}
