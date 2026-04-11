package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	commonrest "goddd/internal/common/interfaces/rest"
	orderrest "goddd/internal/order/interfaces/rest"
)

func main() {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Go Template", "1.0"))

	routes := []commonrest.RouteCollection{
		orderrest.NewOrderRoutes(nil, nil, nil, nil, nil),
	}

	for _, r := range routes {
		r.Register(api)
	}

	spec, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("openapi.json", spec, 0644); err != nil {
		log.Fatal(err)
	}
}
