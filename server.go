package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-simple-graphql/database"
	"go-simple-graphql/graph"
	"go-simple-graphql/graph/generated"
	"go-simple-graphql/graph/middleware"
	"go-simple-graphql/utils"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func NewGraphQLHandler() *chi.Mux {
	var router *chi.Mux = chi.NewRouter()
	router.Use(middleware.NewMiddleware())

	server := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	return router
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var handler *chi.Mux = NewGraphQLHandler()

	err := database.Connect(utils.GetValue("DB_NAME"))
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v\n", err)
	}

	fmt.Println("Connected to the database")

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
