package main

import (
	"golang.org/x/oauth2"
	"flag"
	"net/http"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"io/ioutil"
	"golang.org/x/net/context"
)

var useProxy = true

var igConf *oauth2.Config

var clientID, clientSec string
var urlRoot, urlPort string

func init() {
	flag.StringVar(&clientID, "clientID", "9ac09abdcca64f43bd73ce348eb22139", "clientID")
	flag.StringVar(&clientSec, "clientSec", "25b60514e43a4d2cb912b8e83087a6a2", "clientSec")
	flag.StringVar(&urlRoot, "urlRoot", "127.0.0.1", "urlRoot")
	flag.StringVar(&urlPort, "urlPort", "10098", "urlPort")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		autoCodeUrl := GetInstagramAuthCodeURL("http://" + urlRoot + ":" + urlPort + "/authCallback")
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
	http.ListenAndServe(":" + urlPort, nil)
}

func GetInstagramAuthCodeURL(redirectURL string) string {
	return getInstagramAuthConfig(redirectURL).AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func GetInstagramAuthToken(code string) (*oauth2.Token, error) {
	//return igConf.Exchange(oauth2.NoContext, code)
	return igConf.Exchange(context.WithValue(oauth2.NoContext, oauth2.HTTPClient, GetIGAPIHttpClient()), code)
}

func getInstagramAuthConfig(redirectURL string) *oauth2.Config {
	if igConf == nil {
		igConf = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSec,
			Scopes:       []string{"basic", "follower_list", "public_content", "comments", "relationships", "likes"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.instagram.com/oauth/authorize",
				TokenURL: "https://api.instagram.com/oauth/access_token",
			},
		}
	}
	igConf.RedirectURL = redirectURL
	return igConf
}

var proxyHttpClient *http.Client

func GetIGAPIHttpClient() *http.Client {
	if !useProxy {
		return http.DefaultClient
	}
	if proxyHttpClient == nil {
		dailer, err := proxy.SOCKS5("tcp", "localhost:1080", nil, &net.Dialer{})
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{Dial: dailer.Dial}
		httpClient := &http.Client{Transport: transport}
		resp, err := httpClient.Get("https://www.instagram.com")
		if err != nil {
			panic(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		proxyHttpClient = httpClient
	}
	return proxyHttpClient
}
