package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"subscriptionservice/cmd/migration"
	"subscriptionservice/interal/app"
	"subscriptionservice/interal/server"
	"subscriptionservice/pkg/config"
	"subscriptionservice/pkg/sql/postgres"

	_ "github.com/swaggo/http-swagger"
)

const (
    envLocal = "local"
    envDev   = "dev"
    envProd  = "prod"
)


// @title           Swagger Example API
// @version         2.0
// @description     This is a sample server celler server.

// @host      localhost:8080
// @BasePath  /
func main() {
	loger := setupLogger("local")
	loger = loger.With(slog.String("env", "local"))
    migration.Migrations()
	
    loger.Info("initializing server") 
    
    confPost, err := config.ReturnedDatabase()
    if err!= nil {
        loger.Error("error initializing config", slog.String("error", err.Error()))
        panic(err)
    }

    loger.Info("connecting to database " + confPost.DBName)
    postgres, err := postgres.NewPostgres(confPost.User, confPost.Password, confPost.DBName, confPost.Host, confPost.Port)
    if err != nil {
       loger.Error("error initializing postgres", slog.String("error", err.Error()))
       panic(err)
    }

    loger.Info("initializing server config")
    confServer, err := config.RetuneServerConfig()
    if err!= nil {
        loger.Error("error initializing config", slog.String("error", err.Error()))
        panic(err)
    }

    loger.Info("initializing server app")  
    app := app.NewApp(loger, postgres)
    server := server.NewSubscrServer(loger, app)

    loger.Info("Initializing server endpoints")
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        loger.Info("Received request: " + r.URL.String())
        fmt.Fprint(w, "Server listening on " + r.URL.Host)
    })
   
    http.HandleFunc("/create", server.CreateSubscr)
    http.HandleFunc("/search", server.GetSubscr)
    http.HandleFunc("/delete", server.DeleteSubscr)
    http.HandleFunc("/update", server.UpdateSubscr)
    http.HandleFunc("/summsubscr", server.GetSummSubscr)

    loger.Info("Starting server..." + confServer.Host + " " + confServer.Port)
    err = http.ListenAndServe(confServer.Host + ":" + confServer.Port, nil) 
    if err != nil {
        loger.Error("error starting server", slog.String("error",err.Error()))
        panic(err)
    }
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

    switch env {
    case envLocal:
        logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envDev:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envProd:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    default:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }

    return logger
}

