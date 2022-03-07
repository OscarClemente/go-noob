package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/OscarClemente/go-noob/db"
	"github.com/OscarClemente/go-noob/models"
)

const loadersKey = "dataloaders"

type Loaders struct {
	UserById UserLoader
}

func Middleware(DB db.Database, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserById: UserLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([]*models.User, []error) {
					fmt.Println("Running dataloader")
					users, err := DB.GetUsersById(ids)
					if err != nil {
						return nil, []error{err}
					}

					fmt.Println("Dataloader got: ", len(users))

					return users, nil
				},
			},
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
