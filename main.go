package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/suumiizxc/raw_rest1/api/account"
	"github.com/suumiizxc/raw_rest1/config"
	lg "github.com/suumiizxc/raw_rest1/logger"

	_ "github.com/swaggo/http-swagger/example/gorilla/docs"
)

func contentTypeApplicationJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx := context.WithValue(r.Context(), "data", "data:test")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var (
	trace    = "dev"
	httpPort = 4000
)

func main() {
	config.OpenConnection()

	defer config.DB.Close()
	mux := mux.NewRouter()

	accountH := &account.Account{}
	accountRouter := mux.PathPrefix("/account").Subrouter()
	accountRouter.Use(contentTypeApplicationJsonMiddleware)
	accountRouter.HandleFunc("/create", accountH.RegisterAccount).Methods("POST")

	if trace == "dev" {
		fmt.Printf("listening on %v\n", httpPort)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux))
	} else if trace == "prod" {
		currentTime := time.Now()
		Today := currentTime.Format("2006-01-02")
		fmt.Println("Today : ", Today)
		logPath := fmt.Sprintf("logger/log/%s.log", Today)

		lgr := lg.Logger{}
		lgr.OpenLogFile(logPath)

		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		fmt.Printf("Logging to %v\n", logPath)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), lgr.LogRequest(mux)))
	}
}
