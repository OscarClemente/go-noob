package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/OscarClemente/go-noob/graph/generated"
	"github.com/OscarClemente/go-noob/graph/model"
	"github.com/OscarClemente/go-noob/models"
)

func (r *mutationResolver) CreateReview(ctx context.Context, input model.ReviewInput) (*models.Review, error) {
	review := reviewInputToModel(input)
	err := r.DB.AddReview(review)

	return review, err
}

func (r *mutationResolver) UpdateReview(ctx context.Context, input model.ReviewInput) (*models.Review, error) {
	review := reviewInputToModel(input)
	updatedReview, err := r.DB.UpdateReview(review.ID, *review)

	return &updatedReview, err
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*models.User, error) {
	user := userInputToModel(input)
	err := r.DB.AddUser(user)

	return user, err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UserInput) (*models.User, error) {
	user := userInputToModel(input)
	updatedUser, err := r.DB.UpdateUser(user.ID, *user)

	return &updatedUser, err
}

func (r *queryResolver) Reviews(ctx context.Context) ([]*models.Review, error) {
	reviews, err := r.DB.GetAllReviews()
	if err != nil || reviews == nil {
		return nil, err
	}

	var output []*models.Review
	for _, review := range reviews.Reviews {
		output = append(output, &review)
	}
	return output, err
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users, err := r.DB.GetAllUsers()
	if err != nil || users == nil {
		return nil, err
	}

	var output []*models.User
	for _, user := range users.Users {
		output = append(output, &user)
	}
	return output, err
}

func (r *reviewResolver) User(ctx context.Context, obj *models.Review) (*models.User, error) {
	user, err := r.DB.GetUserById(obj.UserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userResolver) Friends(ctx context.Context, obj *models.User) (*models.User, error) {
	user, err := r.DB.GetUserById(obj.Friends)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Review returns generated.ReviewResolver implementation.
func (r *Resolver) Review() generated.ReviewResolver { return &reviewResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type reviewResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

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
		UserID:  user,
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
