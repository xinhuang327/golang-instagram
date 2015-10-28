package main

import (
	"github.com/xinhuang327/golang-instagram/instagram"
	"flag"
	"net/http"
	"fmt"
)

func main() {
	var clientID, clientSec string
	var urlRoot, urlPort string
	var callbackUrl string
	flag.StringVar(&clientID, "clientID", "", "clientID")
	flag.StringVar(&clientSec, "clientSec", "", "clientID")
	flag.StringVar(&urlRoot, "urlRoot", "127.0.0.1", "urlRoot")
	flag.StringVar(&urlPort, "urlPort", "10097", "urlPort")
	flag.StringVar(&callbackUrl, "callbackUrl", "", "callbackUrl")
	flag.Parse()

	api := instagram.Api{}
	api.ClientId = clientID
	api.ClientSecret = clientSec
	callbackPath := "/realtimeCallback"
	if callbackUrl == "" {
		callbackUrl = "http://"+urlRoot+":"+urlPort+"/"+callbackPath
	}
	http.HandleFunc(callbackPath, api.RealtimeCallback)

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		objType := r.URL.Query().Get("objType")
		arg := r.URL.Query().Get("arg")
		resp, err := api.SubscribeRealtime(objType, arg, callbackUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", resp)
	})
	http.ListenAndServe(":"+urlPort, nil)
}