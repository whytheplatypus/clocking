package auth

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func ListenAndServe(conf *oauth2.Config, next func(rw http.ResponseWriter, r *http.Request)) (tkn *oauth2.Token, err error) {
	srv := &http.Server{
		Addr: ":8080",
	}
	srv.Handler = Handler(conf, func(rw http.ResponseWriter, r *http.Request) {
		tkn = r.Context().Value("Token").(*oauth2.Token)
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		next(rw, r)
		srv.Shutdown(ctx)
	})
	err = srv.ListenAndServeTLS("./cert.pem", "./key.pem")
	return
}
