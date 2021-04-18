package engine

import (
	"net/http"
)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type HandlerFunc func(c *Context)

func New() *Engine {
	router := newRouter()
	group := &RouterGroup{
		e: &Engine{router: router},
	}
	return &Engine{
		RouterGroup: group,
		router:      router,
		groups:      []*RouterGroup{group},
	}
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}

// 实现 http 下的 Handler interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}

func (e *Engine) Get(pattern string, handler HandlerFunc) {
	e.router.Get(pattern, handler)
}

func (e *Engine) Post(pattern string, handler HandlerFunc) {
	e.router.Post(pattern, handler)
}

///  group

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.e
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		parent: rg,
		e:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) addRoute(method, partPath string, handler HandlerFunc) {
	pattern := rg.prefix + partPath
	rg.e.router.addRoute(method, pattern, handler)
}

func (rg *RouterGroup) Get(partPattern string, handler HandlerFunc) {
	rg.addRoute("GET", partPattern, handler)
}

func (rg *RouterGroup) Post(partPattern string, handler HandlerFunc) {
	rg.addRoute("POST", partPattern, handler)
}
