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
	user := userInputToModel(input)
	err := r.DB.AddUser(user)

	return userToOutputModel(user), err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	user := userInputToModel(input)
	updatedUser, err := r.DB.UpdateUser(user.ID, *user)

	return userToOutputModel(&updatedUser), err
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
	users, err := r.DB.GetAllUsers()
	if err != nil || users == nil {
		return nil, err
	}

	var output []*model.User
	for _, user := range users.Users {
		output = append(output, userToOutputModel(&user))
	}
	return output, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func reviewInputToModel(input model.ReviewInput) *models.Review {
	id, _ := strconv.Atoi(input.ID)
	user, _ := strconv.Atoi(input.UserID)
	return &models.Review{
		ID:      id,
		Game:    input.Game,
		Title:   input.Title,
		Content: input.Content,
		Rating:  input.Rating,
		User:    user,
	}
}
func reviewToOutputModel(input *models.Review) *model.Review {
	return &model.Review{
		ID:      strconv.FormatInt(int64(input.ID), 10),
		Game:    input.Game,
		Title:   input.Title,
		Content: input.Content,
		Rating:  input.Rating,
		User:    strconv.FormatInt(int64(input.User), 10),
	}
}
func userInputToModel(input model.UserInput) *models.User {
	id, _ := strconv.Atoi(input.ID)
	var friend int
	if len(input.Friends) > 0 {
		friend, _ = strconv.Atoi(input.Friends[0])
	}
	return &models.User{
		ID:      id,
		Name:    input.Name,
		Email:   input.Email,
		Friends: friend,
	}
}
func userToOutputModel(input *models.User) *model.User {
	return &model.User{
		ID:      fmt.Sprint(input.ID),
		Name:    input.Name,
		Email:   input.Email,
		Friends: strconv.FormatInt(int64(input.Friends), 10),
	}
}
