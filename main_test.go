package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/lib/pq"

	"os"

	"github.com/OscarClemente/go-noob/db"

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

	defer func() {
		if err := database.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	request := httptest.NewRequest("GET", "/query", nil)
	recorder := httptest.NewRecorder()

	GQLApp().router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 but receieved %d", recorder.Code)
	}

	expectedPayload := `[{"id":1,"name":"Punk IPA","consumed":true,"rating":68.29}]`
	actualPayload := recorder.Body.String()

	if actualPayload != expectedPayload {
		t.Fatalf("expected %+v but receieved %+v", expectedPayload, actualPayload)
	}
}

type App struct {
	router *http.ServeMux
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
	defer database.Conn.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: database}}))

	router := http.NewServeMux()
	router.Handle("/query", srv)

	return &App{router: router}
}
