package main

import (
	"context"
	"fmt"
	"log"

	//"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OscarClemente/go-noob/db"
	"github.com/OscarClemente/go-noob/models"

	//"github.com/OscarClemente/go-noob/handler"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/OscarClemente/go-noob/graph"
	"github.com/OscarClemente/go-noob/graph/generated"
)

const defaultPort = "8080"

func main() {
	/*addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}*/
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")
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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	tempSeedData(&database)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
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

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
