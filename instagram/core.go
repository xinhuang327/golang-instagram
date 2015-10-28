// Package instagram provides a minimialist instagram API wrapper.
package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	// "io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
)

var (
	baseUrl = "https://api.instagram.com/v1"
)

type Api struct {
	ClientId     string
	ClientSecret string
	AccessToken  string
}

// Create an API with either a ClientId OR an accessToken. Only one is required. Access tokens are preferred because they keep rate limiting down.
func New(clientId string, accessToken string) *Api {
	if clientId == "" && accessToken == "" {
		panic("ClientId or AccessToken must be given to create an Api")
	}

	return &Api{
		ClientId:    clientId,
		AccessToken: accessToken,
	}
}

// -- Implementation of request --

func buildGetRequest(urlStr string, params url.Values) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// If we are getting, then we can't merge query params
	if params != nil {
		if u.RawQuery != "" {
			return nil, fmt.Errorf("Cannot merge query params in urlStr and params")
		}
		u.RawQuery = params.Encode()
	}

	return http.NewRequest("GET", u.String(), nil)
}

func buildDeleteRequest(urlStr string, params url.Values) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// If we are getting, then we can't merge query params
	if params != nil {
		if u.RawQuery != "" {
			return nil, fmt.Errorf("Cannot merge query params in urlStr and params")
		}
		u.RawQuery = params.Encode()
	}

	return http.NewRequest("DELETE", u.String(), nil)
}

func buildPostRequest(urlStr string, params url.Values) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// If we are getting, then we can't merge query params
	postBody := ""
	if params != nil {
		postBody = params.Encode()
	}

	return http.NewRequest("POST", u.String(), strings.NewReader(postBody))
}

func (api *Api) extendParams(p url.Values) url.Values {
	if p == nil {
		p = url.Values{}
	}
	if api.AccessToken != "" {
		p.Set("access_token", api.AccessToken)
	} else {
		p.Set("client_id", api.ClientId)
	}
	if api.ClientSecret != "" {
		p.Set("client_secret", api.ClientSecret)
	}
	return p
}

func (api *Api) get(path string, params url.Values, r interface{}) error {
	params = api.extendParams(params)
	req, err := buildGetRequest(urlify(path), params)
	if err != nil {
		return err
	}
	return api.do(req, r)
}

func (api *Api) delete(path string, params url.Values, r interface{}) error {
	params = api.extendParams(params)
	req, err := buildDeleteRequest(urlify(path), params)
	if err != nil {
		return err
	}
	return api.do(req, r)
}

func (api *Api) post(path string, params url.Values, r interface{}) error {
	params = api.extendParams(params)
	req, err := buildPostRequest(urlify(path), params)
	if err != nil {
		return err
	}
	return api.do(req, r)
}

var httpClient = http.DefaultClient

func SetHttpClient(client *http.Client) {
	httpClient = client
}

func (api *Api) do(req *http.Request, r interface{}) error {
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return apiError(resp)
	}

	return decodeResponse(resp.Body, r)
}

var PrintRawAPIResponse bool

func decodeResponse(body io.Reader, to interface{}) error {
	var err error
	if PrintRawAPIResponse{
		b, _ := ioutil.ReadAll(body)
		fmt.Println("Body:",string(b))
		err = json.Unmarshal(b, to)
	} else {
		err = json.NewDecoder(body).Decode(to)
	}

	if err != nil {
		return fmt.Errorf("instagram: error decoding body; %s", err.Error())
	}
	return nil
}

func apiError(resp *http.Response) error {
	m := new(MetaResponse)
	if err := decodeResponse(resp.Body, m); err != nil {
		return err
	}

	var err MetaError
	if m.Meta != nil {
		err = MetaError(*m.Meta)
	} else {
		err = MetaError(Meta{Code: resp.StatusCode, ErrorMessage: resp.Status})
	}
	return &err
}

func urlify(path string) string {
	return baseUrl + path
}

type MetaError Meta

func (m *MetaError) Error() string {
	return fmt.Sprintf("Error making api call: Code %d %s %s", m.Code, m.ErrorType, m.ErrorMessage)
}

func ensureParams(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	return v
}
