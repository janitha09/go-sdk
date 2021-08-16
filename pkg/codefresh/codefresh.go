package codefresh

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

//go:generate mockery -name Codefresh -filename codefresh.go

//go:generate mockery -name UsersAPI -filename users.go

type (
	Codefresh interface {
		Pipelines() IPipelineAPI
		Tokens() ITokenAPI
		RuntimeEnvironments() IRuntimeEnvironmentAPI
		Workflows() IWorkflowAPI
		Progresses() IProgressAPI
		Clusters() IClusterAPI
		Contexts() IContextAPI
		Users() UsersAPI
		Argo() ArgoAPI
		Gitops() GitopsAPI
		ArgoRuntime() IArgoRuntimeAPI
		GitSource() IGitSourceAPI
	}
)

func New(opt *ClientOptions) Codefresh {
	httpClient := &http.Client{}
	if opt.Client != nil {
		httpClient = opt.Client
	}

	return &codefresh{
		host:   opt.Host,
		token:  opt.Auth.Token,
		client: httpClient,
	}
}

func (c *codefresh) Pipelines() IPipelineAPI {
	return newPipelineAPI(c)
}

func (c *codefresh) Users() UsersAPI {
	return newUsersAPI(c)
}

func (c *codefresh) Tokens() ITokenAPI {
	return newTokenAPI(c)
}

func (c *codefresh) RuntimeEnvironments() IRuntimeEnvironmentAPI {
	return newRuntimeEnvironmentAPI(c)
}

func (c *codefresh) Workflows() IWorkflowAPI {
	return newWorkflowAPI(c)
}

func (c *codefresh) Progresses() IProgressAPI {
	return newProgressAPI(c)
}

func (c *codefresh) Clusters() IClusterAPI {
	return newClusterAPI(c)
}

func (c *codefresh) Contexts() IContextAPI {
	return newContextAPI(c)
}

func (c *codefresh) Argo() ArgoAPI {
	return newArgoAPI(c)
}

func (c *codefresh) Gitops() GitopsAPI {
	return newGitopsAPI(c)
}

func (c *codefresh) ArgoRuntime() IArgoRuntimeAPI {
	return newArgoRuntimeAPI(c)
}

func (c *codefresh) GitSource() IGitSourceAPI {
	return newGitSourceAPI(c)
}

func (c *codefresh) requestAPI(opt *requestOptions) (*http.Response, error) {
	return c.requestAPIWithContext(context.Background(), opt)
}

func (c *codefresh) requestAPIWithContext(ctx context.Context, opt *requestOptions) (*http.Response, error) {
	var body []byte
	finalURL := fmt.Sprintf("%s%s", c.host, opt.path)
	if opt.qs != nil {
		finalURL += toQS(opt.qs)
	}
	if opt.body != nil {
		body, _ = json.Marshal(opt.body)
	}
	request, err := http.NewRequestWithContext(ctx, opt.method, finalURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", c.token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("origin", c.host)

	response, err := c.client.Do(request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func buildQSFromMap(qs map[string]string) string {
	var arr = []string{}
	for k, v := range qs {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	return "?" + strings.Join(arr, "&")
}

func toQS(qs interface{}) string {
	v, _ := query.Values(qs)
	qsStr := v.Encode()
	if qsStr != "" {
		return "?" + qsStr
	}
	var qsMap map[string]string
	rs, _ := json.Marshal(qs)
	err := json.Unmarshal(rs, &qsMap)
	if err != nil {
		return ""
	}
	return buildQSFromMap(qsMap)
}

func (c *codefresh) decodeResponseInto(resp *http.Response, target interface{}) error {
	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *codefresh) getBodyAsString(resp *http.Response) (string, error) {
	body, err := c.getBodyAsBytes(resp)
	return string(body), err
}

func (c *codefresh) getBodyAsBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
