package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/suumiizxc/raw_rest1/api/account"
	"github.com/suumiizxc/raw_rest1/config"
	"github.com/suumiizxc/raw_rest1/response"
)

func ContentTypeApplicationJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx := context.WithValue(r.Context(), "data", "data:test")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("access_token")
		var account account.Account
		acc, _ := config.RS.Get(token).Result()
		err := json.Unmarshal([]byte(acc), &account)
		fmt.Println("account : ", account)
		resp := response.Response{}
		if err != nil {
			resp.Error = "Failed in autherize token"
			resp.Message = "Failed in autherize token"
			w.WriteHeader(http.StatusNotImplemented)
			w.Write(resp.ConvertByte())
			return
		}

		ctx := context.WithValue(r.Context(), "user", account.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
