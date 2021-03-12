package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/assignments-fixed-ssunni12/servers/gateway/handlers"
	"github.com/assignments-fixed-ssunni12/servers/gateway/models/users"
	"github.com/assignments-fixed-ssunni12/servers/gateway/sessions"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

//main is the main entry point for the server
func main() {
	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
	  the server should listen on. If empty, default to ":80"
	- Create a new mux for the web server.
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/
	ADDR := os.Getenv("ADDR")
	if len(ADDR) == 0 {
		ADDR = ":80"
	}
	TLSCERT := os.Getenv("TLSCERT")
	if len(TLSCERT) == 0 {
		fmt.Println("TLSCERT environment variable wasn't set")
		os.Exit(1)
	}
	TLSKEY := os.Getenv("TLSKEY")
	if len(TLSKEY) == 0 {
		fmt.Println("TLSKEY environment variable wasn't set")
		os.Exit(1)
	}
	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		fmt.Println("SESSIONKEY Env variable wasn't set")
		os.Exit(1)
	}
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "redisServer:6379"
	}
	dsn := os.Getenv("DSN")
	if len(dsn) == 0 {
		fmt.Println("DSN Env variable wasn't set")
		os.Exit(1)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Error opening database: %v", err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error opening database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	usersStore, err := users.NewMySQLStore(dsn)
	if err != nil {
		fmt.Printf("Error opening database: %v", err)
		os.Exit(1)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	sessionStore := sessions.NewRedisStore(redisClient, time.Hour)
	context := handlers.NewContext(sessionKey, sessionStore, usersStore)

	mux := http.NewServeMux()
	log.Printf("Cert: %s\nKey: %s", TLSCERT, TLSKEY)
	log.Printf("Server running, listening on %s", ADDR)

	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)
	wrappedMux := handlers.NewHeaderHandler(mux)

	log.Fatal(http.ListenAndServeTLS(ADDR, TLSCERT, TLSKEY, wrappedMux))

}
