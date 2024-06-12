package client

import (
	"Food/internal/errors"
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
)

const (
	defaultBaseURL = "https://7o9mo.wiremockapi.cloud/"
	defaultTimeout = 60 * time.Second
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client ...
type Client struct {
	httpClient HTTPClient
	baseURL    string
	debug      bool
}

// NewClient creates a new Spend-Juice API client with the default base URL.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    defaultBaseURL,
		debug:      os.Getenv("ENV") != "production",
	}
}

// SetHTTPClient sets the HTTP client that will be used for API calls.
func (cl *Client) SetHTTPClient(httpClient HTTPClient) {
	cl.httpClient = httpClient
}

// SetBaseURL overrides the default base URL. For internal use.
func (cl *Client) SetBaseURL(baseURL string) {
	cl.baseURL = strings.TrimRight(baseURL, "/")
}

// Get Base URL
func (cl *Client) GetBaseURL() string {
	return cl.baseURL
}

func (cl *Client) Get(path string, params interface{}, response interface{}) (err error) {
	if params != nil {

		_, err = valid.ValidateStruct(params)
		if err != nil {
			return
		}

		v, _ := query.Values(params)
		path = path + "?" + v.Encode()
	}

	url := cl.baseURL + "/" + strings.TrimLeft(path, "/")

	if cl.debug {
		log.Printf("cellulant: Call: %s %s", "GET", url)
		log.Printf("cellulant: Request Params: %#v", params)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	return cl.request(req, response)
}

func (cl *Client) Post(path string, params interface{}, response interface{}) (err error) {
	url := cl.baseURL + "/" + strings.TrimLeft(path, "/")

	var req *http.Request
	var bodyBuffered io.Reader

	if params != nil {
		_, err = valid.ValidateStruct(params)
		if err != nil {
			return
		}

		data, _ := json.Marshal(params)
		bodyBuffered = bytes.NewBuffer([]byte(data))

	}

	if cl.debug {
		log.Printf("cellulant: Call: %s %s", "POST", url)
		log.Printf("cellulant: Request Params: %#v", params)
	}

	req, err = http.NewRequest(http.MethodPost, url, bodyBuffered)

	if err != nil {
		return
	}

	return cl.request(req, response)
}

func (cl *Client) request(req *http.Request, response interface{}) (err error) {

	req.Header.Set("Content-Type", "application/json")

	r, err := cl.httpClient.Do(req)

	if err != nil {
		return
	}

	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		e := errors.ErrResponse{}
		err = json.NewDecoder(r.Body).Decode(&e)

		if err != nil {
			return err
		}

		return e
	}

	err = json.NewDecoder(r.Body).Decode(response)
	if err != nil {
		log.Printf("Failed to decode response: %v", err)
	}

	return
}
