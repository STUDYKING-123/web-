package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 本文件目的是将net/http中调用进行一些封装
type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type node struct {
	pattern     string
	handlerFunc HandlerFunc
	children    map[string]*node
}
type router struct {
	grRouter map[string]*groupRouter // 为了提供Group方法
}

// engine.Group("/user")则会创建一个以user开头的
type groupRouter struct {
	prefix string
	trees  map[string]*node
}

func (r *router) Group(prefix string) *groupRouter {
	if r.grRouter == nil {
		r.grRouter = make(map[string]*groupRouter)
	}
	gpRouter := &groupRouter{
		prefix: prefix,
		trees:  make(map[string]*node),
	}
	r.grRouter[prefix] = gpRouter
	return gpRouter
}
func (r *groupRouter) addroute(method, name string, handleFunc HandlerFunc) {
	if r.trees == nil {
		r.trees = make(map[string]*node)
	}
	var paths []string
	if r.prefix != "" {
		paths = append(paths, r.prefix)
		names := strings.Split(name, "/")[1:] // 将/user/login 分割成 user login
		for _, v := range names {
			paths = append(paths, v)
		}
	} else {
		names := strings.Split(name, "/")[1:] // 将/user/login 分割成 user login
		for _, v := range names {
			paths = append(paths, v)
		}
	}
	root, ok := r.trees[method]
	if !ok {
		root = &node{pattern: "/", children: make(map[string]*node)}
		r.trees[method] = root
	}
	for _, s := range paths {
		if root.children[s] == nil {
			root.children = make(map[string]*node)
		}
		chid, ok := root.children[s]
		if !ok {
			chid = newnode(s)
			//fmt.Println(chid.pattern)
			root.children[s] = chid
			root = chid
		} else {
			root = chid
		}
	}
	root.handlerFunc = handleFunc
	//fmt.Println(paths)
}
func newnode(pattern string) *node {
	return &node{
		pattern:  pattern,
		children: make(map[string]*node),
	}
}

// 支持restful四种请求，将addroute隐藏，对外暴露下面四个
func (r *groupRouter) Get(name string, handlerFunc HandlerFunc) {
	r.addroute(http.MethodGet, name, handlerFunc)
}
func (r *groupRouter) POST(name string, handlerFunc HandlerFunc) {
	r.addroute(http.MethodPost, name, handlerFunc)
}
func (r *groupRouter) PUT(name string, handlerFunc HandlerFunc) {
	r.addroute(http.MethodPut, name, handlerFunc)
}
func (r *groupRouter) DEL(name string, handlerFunc HandlerFunc) {
	r.addroute(http.MethodDelete, name, handlerFunc)
}

var _ http.Handler = &Engine{}

type Engine struct {
	router
}

// 对外暴露Engine,所有操作通过Engine完成，将路由添加查找部分抽离出来给router
func NewEngine() *Engine {
	return &Engine{
		router{},
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Path) 路径参数
	//fmt.Println(r.Method)
	path := r.URL.Path
	paths := strings.Split(path, "/")[1:]
	group := paths[0]
	fmt.Println(group)
	grRouter, ok := e.grRouter[group]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("找不到页面"))
		return
	}
	root, ok := grRouter.trees[r.Method]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("找不到页面"))
		return
	}
	for _, s := range paths {
		chid, ok := root.children[s]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("找不到页面"))
			return
		}
		root = chid
	}
	root.handlerFunc(w, r)
}

func (e *Engine) Run(address string) {

	err := http.ListenAndServe(address, e)
	if err != nil {
		log.Fatal(err)
	}
}
