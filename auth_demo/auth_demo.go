package main

import (
	"golang.org/x/oauth2"
	"flag"
	"net/http"
	"fmt"
)

var igConf *oauth2.Config

var clientID, clientSec string
var urlRoot, urlPort string

func init() {
	flag.StringVar(&clientID, "clientID", "", "clientID")
	flag.StringVar(&clientSec, "clientSec", "", "clientSec")
	flag.StringVar(&urlRoot, "urlRoot", "http://127.0.0.1", "urlRoot")
	flag.StringVar(&urlPort, "urlPort", "10099", "urlPort")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		autoCodeUrl := GetInstagramAuthCodeURL(urlRoot+"/authCallback")
		http.Redirect(w, r, autoCodeUrl, http.StatusFound)
		return
	})
	http.HandleFunc("/authCallback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := GetInstagramAuthToken(r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, tok.AccessToken)
	})
	http.ListenAndServe(":"+urlPort, nil)
}

func GetInstagramAuthCodeURL(redirectURL string) string {
	return getInstagramAuthConfig(redirectURL).AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func GetInstagramAuthToken(code string) (*oauth2.Token, error) {
	return igConf.Exchange(oauth2.NoContext, code)
}

func getInstagramAuthConfig(redirectURL string) *oauth2.Config {
	if igConf == nil {
		igConf = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSec,
			Scopes:       []string{"basic", "comments", "relationships", "likes"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.instagram.com/oauth/authorize",
				TokenURL: "https://api.instagram.com/oauth/access_token",
			},
		}
	}
	igConf.RedirectURL = redirectURL
	return igConf
}
