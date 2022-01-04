package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/zepyrshut/breakfast-n-go/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, type is %T", v))

	}

}
