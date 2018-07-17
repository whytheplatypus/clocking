package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

func Handler(conf *oauth2.Config, n func(http.ResponseWriter, *http.Request)) http.Handler {
	return &TokenHandler{
		Config: conf,
		Next:   n,
	}
}

type TokenHandler struct {
	Config *oauth2.Config
	Next   func(http.ResponseWriter, *http.Request)
}

func (th *TokenHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	/*breq, _ := httputil.DumpRequest(r, true)
	log.Println(string(breq))*/
	if code := r.FormValue("code"); code != "" {
		th.ExchangeHandler(th.Next)(rw, r)
		return
	}
	if err := r.FormValue("error"); err != "" {
		http.Error(
			rw,
			err,
			http.StatusBadRequest,
		)
		return
	}
	th.RedirectHandler(rw, r)
}

func (th *TokenHandler) RedirectHandler(rw http.ResponseWriter, r *http.Request) {
	authURL := th.Config.AuthCodeURL("bobert")
	http.Redirect(rw, r, authURL, http.StatusSeeOther)
}

func (th *TokenHandler) ExchangeHandler(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		tok, err := th.Config.Exchange(r.Context(), strings.TrimSpace(code))
		if err != nil {
			log.Println("[ERROR] bad token exchange", err, tok)
			http.Error(
				rw,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		ctx := context.WithValue(r.Context(), "Token", tok)
		next(rw, r.WithContext(ctx))
		return
	}
}
