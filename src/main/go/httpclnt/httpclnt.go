package httpclnt

import (
	"context"
	"errors"
	"fmt"
	"github.com/engswee/flashpipe/logger"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"time"
)

type HTTPExecuter struct {
	basicUserId   string
	basicPassword string
	host          string
	httpClient    *http.Client
	AuthType      string
}

// New returns an initialised HTTPExecuter instance.
func New(oauthHost string, oauthPath string, clientId string, clientSecret string, userId string, password string, host string) *HTTPExecuter {
	e := new(HTTPExecuter)
	e.host = host
	if oauthHost != "" {
		logger.Debug("Initialising HTTP client with OAuth 2.0")

		tokenURL := fmt.Sprintf("https://%v%v", oauthHost, oauthPath)
		logger.Debug(fmt.Sprintf("Getting OAuth 2.0 client with token URL %v", tokenURL))

		// Reference https://pkg.go.dev/golang.org/x/oauth2/clientcredentials#pkg-overview
		conf := &clientcredentials.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			TokenURL:     tokenURL,
		}

		ctx := context.Background()
		e.httpClient = conf.Client(ctx)
		e.AuthType = "OAUTH"
	} else {
		logger.Debug("Initialising HTTP client with Basic Authentication")
		e.httpClient = &http.Client{Timeout: 5 * time.Second}
		e.basicUserId = userId
		e.basicPassword = password
		e.AuthType = "BASIC"
	}
	return e
}

func (e *HTTPExecuter) ExecRequestWithCookies(method string, path string, body io.Reader, headers map[string]string, cookies []*http.Cookie) (resp *http.Response, err error) {

	url := fmt.Sprintf("https://%v%v", e.host, path)
	logger.Debug(fmt.Sprintf("Executing HTTP request: %v %v", method, url))

	// Create new HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	// Set basic authentication if needed
	if e.basicUserId != "" {
		req.SetBasicAuth(e.basicUserId, e.basicPassword)
	}

	// Set HTTP headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Set cookies
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	// Execute HTTP request
	return e.httpClient.Do(req)
}

func (e *HTTPExecuter) ExecRequest(method string, path string, body io.Reader, headers map[string]string) (resp *http.Response, err error) {
	return e.ExecRequestWithCookies(method, path, body, headers, nil)
}

func (e *HTTPExecuter) ExecGetRequest(path string, headers map[string]string) (resp *http.Response, err error) {
	return e.ExecRequestWithCookies(http.MethodGet, path, http.NoBody, headers, nil)
}

func (e *HTTPExecuter) ReadRespBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (e *HTTPExecuter) LogError(resp *http.Response, callType string) (err error) {
	resBody, err := e.ReadRespBody(resp)
	if err != nil {
		return
	}

	body := string(resBody)
	if body != "" {
		logger.Error(fmt.Sprintf("Response Body = %s", body))
	}

	return errors.New(fmt.Sprintf("%v call failed with response code = %d", callType, resp.StatusCode))
}