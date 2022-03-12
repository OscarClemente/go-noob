package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/lib/pq"

	"os"

	"github.com/OscarClemente/go-noob/db"
	"github.com/OscarClemente/go-noob/models"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/OscarClemente/go-noob/graph"
	"github.com/OscarClemente/go-noob/graph/generated"
)

func Test_SimpleHttpWebApp(t *testing.T) {
	dbUser := "postgres"
	dbPassword := "postgres_pass"
	dbName := "noob_db"
	database := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Username(dbUser).Password(dbPassword).Database(dbName).Port(8087).StartTimeout(15 * time.Second))
	if err := database.Start(); err != nil {
		t.Fatal(err)
	}

	jsonData := map[string]string{
		"query": `
            { 
                reviews {
                    id
                    title
                }
            }
        `,
	}
	jsonValue, _ := json.Marshal(jsonData)

	request := httptest.NewRequest("POST", "/query", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	recorder := httptest.NewRecorder()

	app := GQLApp()
	db.Migrations(app.db)
	tempSeedData(&app.db)
	fmt.Println("Populated----------------------")
	app.router.ServeHTTP(recorder, request)

	fmt.Println("SERVED----------------------")

	defer app.db.Conn.Close()
	defer func() {
		if err := database.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 but receieved %d, %q", recorder.Code, recorder.Body)
	}

	expectedPayload := `{"data":{"reviews":[{"id":1,"title":"Chill exploration"},{"id":2,"title":"Best game"},{"id":3,"title":"Be a detective"},{"id":4,"title":"Boom boom pow"},{"id":5,"title":"Classic shooter"}]}}`
	actualPayload := recorder.Body.String()

	if actualPayload != expectedPayload {
		t.Fatalf("received %+v", actualPayload)
	}
}

type App struct {
	router *http.ServeMux
	db     db.Database
}

func (a *App) Start() error {
	return http.ListenAndServe("localhost:8080", a.router)
}

const defaultPort = "8080"

func GQLApp() *App {
	dbUser := "postgres"
	dbPassword := "postgres_pass"
	dbName := "noob_db"
	database, err := db.Initialize(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: database}}))

	router := http.NewServeMux()
	router.Handle("/query", srv)

	return &App{router: router, db: database}
}

func tempSeedData(db *db.Database) {
	reviews, err := db.GetAllReviews()
	if err != nil || reviews == nil || len(reviews.Reviews) > 0 {
		return
	}
	db.AddReview(&models.Review{
		ID:      1,
		Game:    "Sable",
		Title:   "Chill exploration",
		Content: "Cool game in gamepass, entertaining and chill.",
		Rating:  4,
		UserID:  1,
	})
	db.AddReview(&models.Review{
		ID:      2,
		Game:    "Outer wilds",
		Title:   "Best game",
		Content: "Cool game in gamepass, nice world and great wow factor.",
		Rating:  5,
		UserID:  1,
	})
	db.AddReview(&models.Review{
		ID:      3,
		Game:    "Return of the Obra Dinn",
		Title:   "Be a detective",
		Content: "You feel like a real detective, great bell soundtrack.",
		Rating:  4,
		UserID:  1,
	})
	db.AddReview(&models.Review{
		ID:      4,
		Game:    "Dusk",
		Title:   "Boom boom pow",
		Content: "High speed gun carnage.",
		Rating:  5,
		UserID:  2,
	})
	db.AddReview(&models.Review{
		ID:      5,
		Game:    "Doom",
		Title:   "Classic shooter",
		Content: "Nothing else to say.",
		Rating:  5,
		UserID:  2,
	})
	db.AddUser(&models.User{
		ID:    1,
		Name:  "Player1",
		Email: "player1@xbox.com",
	})
	db.AddUser(&models.User{
		ID:    2,
		Name:  "Player2",
		Email: "player2@steam.com",
	})
	db.AddFriendToUser(1, 2)
	db.AddFriendToUser(2, 1)
}
