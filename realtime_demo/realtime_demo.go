package main

import (
	"github.com/xinhuang327/golang-instagram/instagram"
	"flag"
	"net/http"
	"fmt"
	"time"
)

func main() {
	instagram.PrintRawAPIResponse = true

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

	http.HandleFunc("/list",  func(w http.ResponseWriter, r *http.Request) {
		resp, err := api.GetRealtimeSubscriptions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", resp)
	})

	http.HandleFunc("/delete",  func(w http.ResponseWriter, r *http.Request) {
		resp, err := api.DeleteAllRealtimeSubscriptions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", resp)
	})

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

	http.HandleFunc("/ensureSubscribed", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			for {
				resp, err := api.GetRealtimeSubscriptions()
				if err == nil {
					if len(resp.RealtimeSubscriptionList) == 0 {
						res, err := api.SubscribeRealtime("user", "", callbackUrl)
						if err == nil && res.RealtimeSubscription.CallbackUrl != ""{
							fmt.Println("Renewed.")
						} else {
							fmt.Println(err, res)
						}
					}
				} else {
					fmt.Println(err)
				}
				<-time.After(10*time.Minute)
			}
		}()
	})

	http.ListenAndServe(":"+urlPort, nil)
}