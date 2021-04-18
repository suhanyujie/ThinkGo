package engine

import "testing"

func TestAddRoute(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/hello", nil)
	router.addRoute("GET", "/hello/:name", nil)
	t.Log(router)
}
