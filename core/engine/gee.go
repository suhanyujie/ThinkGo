package engine

import (
	"net/http"
	"path"
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

///  group 路由分组

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

// 静态文件服务
// relativePath 参数仅用于路由匹配，不用于寻找文件。
func (rg *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(rg.prefix, relativePath)
	// url 中以 absolutePath 开头的静态文件请求。
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.GetParam("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (rg *RouterGroup) Static(relativePath string, root string) {
	handler := rg.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	rg.Get(urlPattern, handler)
}
