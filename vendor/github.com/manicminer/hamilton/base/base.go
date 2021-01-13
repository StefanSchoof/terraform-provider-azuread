package base

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/manicminer/hamilton/auth"
)

const (
	DefaultEndpoint = "https://graph.microsoft.com"
	Version10       = "v1.0"
	VersionBeta     = "beta"
)

type ValidStatusFunc func(response *http.Response) bool

type HttpRequestInput interface {
	GetValidStatusCodes() []int
	GetValidStatusFunc() ValidStatusFunc
}

type Uri struct {
	Entity      string
	Params      url.Values
	HasTenantId bool
}

type GraphClient = *http.Client

type Client struct {
	ApiVersion string
	Endpoint   string
	TenantId   string
	UserAgent  string

	Authorizer auth.Authorizer
	httpClient GraphClient
}

func NewClient(endpoint, tenantId, version string) Client {
	return Client{
		httpClient: http.DefaultClient,
		Endpoint:   endpoint,
		TenantId:   tenantId,
		ApiVersion: version,
	}
}

func (c Client) buildUri(uri Uri) (string, error) {
	url, err := url.Parse(c.Endpoint)
	if err != nil {
		return "", err
	}
	url.Path = "/" + c.ApiVersion
	if uri.HasTenantId {
		url.Path = fmt.Sprintf("%s/%s", url.Path, c.TenantId)
	}
	url.Path = fmt.Sprintf("%s/%s", url.Path, strings.TrimLeft(uri.Entity, "/"))
	if uri.Params != nil {
		url.RawQuery = uri.Params.Encode()
	}
	return url.String(), nil
}

func (c Client) performRequest(req *http.Request, input HttpRequestInput) (*http.Response, int, error) {
	var status int

	if c.Authorizer != nil {
		token, err := c.Authorizer.Token()
		if err != nil {
			return nil, status, err
		}
		token.SetAuthHeader(req)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, status, err
	}

	status = resp.StatusCode
	if !containsStatusCode(input.GetValidStatusCodes(), status) {
		f := input.GetValidStatusFunc()
		if f != nil && f(resp) {
			return resp, status, nil
		}

		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, status, fmt.Errorf("unexpected status %d with response: %s", resp.StatusCode, string(respBody))
	}

	return resp, status, nil
}

func containsStatusCode(expected []int, actual int) bool {
	for _, v := range expected {
		if actual == v {
			return true
		}
	}

	return false
}
