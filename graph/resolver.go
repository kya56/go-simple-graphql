package graph

import "go-simple-graphql/graph/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	blogService service.BlogService
	userService service.UserService
}
