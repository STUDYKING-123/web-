package test

import (
	"code/router"
	"fmt"
	"net/http"
	"testing"
)

func TestEngine(t *testing.T) {
	engine := router.NewEngine()
	g := engine.Group("user")
	g.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Get login")
	})
	g.POST("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Get login")
	})
	g.DEL("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Get login")
	})
	g.PUT("/lsjadl", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Get login")
	})
	engine.Run(":8081")
}
