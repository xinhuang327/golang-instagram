package instagram

import (
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
)

const (
	RT_User = "user"
	RT_Tag = "tag"
	RT_Location = "location"
)

func (api *Api) SubscribeRealtime(objType string, arg string, callbackUrl string) (res *RealtimeSubscriptionResposne, err error) {
	res = new(RealtimeSubscriptionResposne)
	params := url.Values{}
	params.Set("object", objType)
	params.Set("aspect", "media")
	params.Set("object_id", arg)
	params.Set("callback_url", callbackUrl)
	params.Set("verify_token", "golang")
	err = api.post("/subscriptions", params, res)
	return
}

func (api *Api) GetRealtimeSubscriptions() (res *RealtimeSubscriptionListResposne, err error) {
	res = new(RealtimeSubscriptionListResposne)
	err = api.get("/subscriptions", nil, res)
	return
}

func (api *Api) RealtimeCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RealtimeCallback got")
	if r.Method == "GET" && r.URL.Query().Get("hub.mode") == "subscribe"{
		fmt.Fprint(w, r.URL.Query().Get("hub.challenge"))
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Errorf("RealtimeCallback error: %s", err.Error())
		}
		fmt.Println(string(body))
	}
}