package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/server/middleware"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	
	return &Router{
		engine: engine,
	}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) Group(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{
		group: r.engine.Group(path, handlers...),
	}
}

func (r *Router) Use(middleware ...gin.HandlerFunc) {
	r.engine.Use(middleware...)
}

func (r *Router) GET(path string, handlers ...gin.HandlerFunc) {
	r.engine.GET(path, handlers...)
}

func (r *Router) POST(path string, handlers ...gin.HandlerFunc) {
	r.engine.POST(path, handlers...)
}

func (r *Router) PUT(path string, handlers ...gin.HandlerFunc) {
	r.engine.PUT(path, handlers...)
}

func (r *Router) DELETE(path string, handlers ...gin.HandlerFunc) {
	r.engine.DELETE(path, handlers...)
}

func (r *Router) PATCH(path string, handlers ...gin.HandlerFunc) {
	r.engine.PATCH(path, handlers...)
}

func (r *Router) OPTIONS(path string, handlers ...gin.HandlerFunc) {
	r.engine.OPTIONS(path, handlers...)
}

func (r *Router) HEAD(path string, handlers ...gin.HandlerFunc) {
	r.engine.HEAD(path, handlers...)
}

func (r *Router) Static(path, root string) {
	r.engine.Static(path, root)
}

func (r *Router) StaticFile(path, filepath string) {
	r.engine.StaticFile(path, filepath)
}

func (r *Router) StaticFS(path string, fs http.FileSystem) {
	r.engine.StaticFS(path, fs)
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

type RouterGroup struct {
	group *gin.RouterGroup
}

func (rg *RouterGroup) Group(path string, handlers ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{
		group: rg.group.Group(path, handlers...),
	}
}

func (rg *RouterGroup) Use(middleware ...gin.HandlerFunc) {
	rg.group.Use(middleware...)
}

func (rg *RouterGroup) GET(path string, handlers ...gin.HandlerFunc) {
	rg.group.GET(path, handlers...)
}

func (rg *RouterGroup) POST(path string, handlers ...gin.HandlerFunc) {
	rg.group.POST(path, handlers...)
}

func (rg *RouterGroup) PUT(path string, handlers ...gin.HandlerFunc) {
	rg.group.PUT(path, handlers...)
}

func (rg *RouterGroup) DELETE(path string, handlers ...gin.HandlerFunc) {
	rg.group.DELETE(path, handlers...)
}

func (rg *RouterGroup) PATCH(path string, handlers ...gin.HandlerFunc) {
	rg.group.PATCH(path, handlers...)
}

func (rg *RouterGroup) OPTIONS(path string, handlers ...gin.HandlerFunc) {
	rg.group.OPTIONS(path, handlers...)
}

func (rg *RouterGroup) HEAD(path string, handlers ...gin.HandlerFunc) {
	rg.group.HEAD(path, handlers...)
}

func (rg *RouterGroup) Static(path, root string) {
	rg.group.Static(path, root)
}

func (rg *RouterGroup) StaticFile(path, filepath string) {
	rg.group.StaticFile(path, filepath)
}

func (rg *RouterGroup) StaticFS(path string, fs http.FileSystem) {
	rg.group.StaticFS(path, fs)
}

func (rg *RouterGroup) WithAuth() *RouterGroup {
	rg.group.Use(middleware.Auth())
	return rg
}

func (rg *RouterGroup) WithDomainAuth() *RouterGroup {
	rg.group.Use(middleware.DomainAuth())
	return rg
}
