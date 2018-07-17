package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/whytheplatypus/clocking/clocks/timely/auth"
	"golang.org/x/oauth2"
)

const RedirectURL = "https://localhost:8080"

func main() {
	conf := &oauth2.Config{
		RedirectURL:  RedirectURL,
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.timelyapp.com/1.1/oauth/authorize",
			TokenURL: "https://api.timelyapp.com/1.1/oauth/token",
		},
	}

	tkn, err := auth.ListenAndServe(conf, func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Well done. your token is saved"))
	})
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	f, err := os.Create("/home/whytheplatypus/.timely.tkn")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	WriteTkn(f, tkn)
}

func WriteTkn(w io.Writer, tkn *oauth2.Token) error {
	t := &Token{*tkn}
	_, err := fmt.Fprintln(w, t)
	return err
}

type Token struct {
	oauth2.Token
}

func (tok *Token) String() string {
	//stok, err := json.MarshalIndent(tok, "", "	")
	stok, err := json.Marshal(tok)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(stok)
}
