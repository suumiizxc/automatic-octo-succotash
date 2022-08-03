package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/suumiizxc/raw_rest1/api/account"
	"github.com/suumiizxc/raw_rest1/api/content"
	"github.com/suumiizxc/raw_rest1/config"
	lg "github.com/suumiizxc/raw_rest1/logger"
	middleware "github.com/suumiizxc/raw_rest1/middleware"

	_ "github.com/swaggo/http-swagger/example/gorilla/docs"
)

// func contentTypeApplicationJsonMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		ctx := context.WithValue(r.Context(), "data", "data:test")
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

var (
	trace    = "dev"
	httpPort = 4000
)

func main() {
	config.OpenConnection()
	config.RedisConfig()
	defer config.DB.Close()
	mux := mux.NewRouter()

	// route: /account
	accountH := &account.Account{}
	accountRouter := mux.PathPrefix("/account").Subrouter()
	accountRouter.Use(middleware.ContentTypeApplicationJsonMiddleware)
	accountRouter.HandleFunc("/create", accountH.RegisterAccount).Methods("POST")
	accountRouter.HandleFunc("/login", accountH.LoginAccount).Methods("POST")

	// route authenticated: /account
	accountAuthRouter := accountRouter.NewRoute().Subrouter()
	accountAuthRouter.Use(middleware.Authenticate)
	accountAuthRouter.HandleFunc("/profile", accountH.ProfileAccount).Methods("GET")

	// route: /content
	contentH := &content.Content{}
	contentRouter := mux.PathPrefix("/content").Subrouter()
	contentRouter.HandleFunc("/create", contentH.CreateContent).Methods("POST")
	contentRouter.HandleFunc("/get-by-id/{id}", contentH.GetContentById).Methods("GET")

	// route: /content/author
	contentAuthorH := &content.ContentAuthor{}
	contentAuthorRouter := contentRouter.NewRoute().PathPrefix("/author").Subrouter()
	contentAuthorRouter.HandleFunc("/create", contentAuthorH.CreateContentAuthor).Methods("POST")

	// route: /content/file
	contentFileH := &content.ContentFile{}
	contentFileRouter := contentRouter.NewRoute().PathPrefix("/file").Subrouter()
	contentFileRouter.HandleFunc("/create", contentFileH.CreateContentFile).Methods("POST")

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
