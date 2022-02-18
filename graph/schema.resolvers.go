package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/OscarClemente/go-noob/graph/generated"
	"github.com/OscarClemente/go-noob/graph/model"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	return nil, fmt.Errorf("createtodo not yet implemented but query received")
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return nil, fmt.Errorf("todos not yet implemented but query received")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
