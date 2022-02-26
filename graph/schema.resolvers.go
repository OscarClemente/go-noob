package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/OscarClemente/go-noob/graph/generated"
	"github.com/OscarClemente/go-noob/graph/model"
	"github.com/OscarClemente/go-noob/models"
)

func reviewInputToModel(input model.ReviewInput) *models.Review {
	user, _ := strconv.Atoi(input.UserID)
	return &models.Review{
		Game:    input.Game,
		Title:   input.Title,
		Content: input.Content,
		Rating:  input.Rating,
		User:    user,
	}
}

func reviewToOutputModel(input *models.Review) *model.Review {
	return &model.Review{
		Game:    input.Game,
		Title:   input.Title,
		Content: input.Content,
		Rating:  input.Rating,
		User:    strconv.FormatInt(int64(input.User), 10),
	}
}

func (r *mutationResolver) CreateReview(ctx context.Context, input model.ReviewInput) (*model.Review, error) {
	review := reviewInputToModel(input)
	err := r.DB.AddReview(review)

	return reviewToOutputModel(review), err
}

func (r *mutationResolver) UpdateReview(ctx context.Context, input model.ReviewInput) (*model.Review, error) {
	review := reviewInputToModel(input)
	updatedReview, err := r.DB.UpdateReview(review.ID, *review)

	return reviewToOutputModel(&updatedReview), err
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Reviews(ctx context.Context) ([]*model.Review, error) {
	reviews, err := r.DB.GetAllReviews()
	if err != nil || reviews == nil {
		return nil, err
	}

	var output []*model.Review
	for _, review := range reviews.Reviews {
		output = append(output, reviewToOutputModel(&review))
	}
	return output, err
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
