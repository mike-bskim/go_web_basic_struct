package main

import (
	"GO/Toy_Prj/basic_struct/internal/config"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case (*chi.Mux):
		fmt.Printf("(*chi.Mux) type is %T\n", v)
		// do nothing; test passed
	case http.Handler:
		fmt.Printf("http.Handler type is %T\n", v)
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, but is %T", v))
	}
}
