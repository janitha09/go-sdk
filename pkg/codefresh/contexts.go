package codefresh

import "github.com/codefresh-io/go-sdk/pkg/utils"

type (
	IContextAPI interface {
		GetGitContexts() (error, *[]ContextPayload)
		GetGitContextByName(name string) (error, *ContextPayload)
		GetDefaultGitContext() (error, *ContextPayload)
	}

	context struct {
		codefresh *codefresh
	}

	ContextPayload struct {
		Metadata struct {
			Name string `json:"name"`
		}
		Spec struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Type           string `json:"type"`
					Username       string `json:"username"`
					Password       string `json:"password"`
					ApiHost        string `json:"apiHost"`
					ApiPathPrefix  string `json:"apiPathPrefix"`
					SshPrivateKey  string `json:"sshPrivateKey"`
					AppId          string `json:"appId"`
					InstallationId string `json:"installationId"`
					PrivateKey     string `json:"privateKey"`
				} `json:"auth"`
			} `json:"data"`
		} `json:"spec"`
	}

	GitContextsQs struct {
		Type    []string `json:"type"`
		Decrypt string   `json:"decrypt"`
	}
)

func newContextAPI(codefresh *codefresh) IContextAPI {
	return &context{codefresh}
}

func (c context) GetGitContexts() (error, *[]ContextPayload) {
	var result []ContextPayload
	var qs map[string]string

	gitContextsQs := GitContextsQs{
		Type:    []string{"git.github", "git.gitlab", "git.github-app"},
		Decrypt: "true",
	}

	utils.Convert(gitContextsQs, &qs)

	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/contexts",
		qs:     qs,
	})
	if err != nil {
		return err, nil
	}

	err = c.codefresh.decodeResponseInto(resp, &result)

	return err, &result
}

func (c context) GetGitContextByName(name string) (error, *ContextPayload) {
	var result ContextPayload
	var qs = map[string]string{
		"decrypt": "true",
	}

	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/contexts/" + name,
		qs:     qs,
	})
	if err != nil {
		return err, nil
	}

	err = c.codefresh.decodeResponseInto(resp, &result)

	return nil, &result
}

func (c context) GetDefaultGitContext() (error, *ContextPayload) {
	var result ContextPayload

	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/contexts/git/default",
	})

	if err != nil {
		return err, nil
	}

	err = c.codefresh.decodeResponseInto(resp, &result)

	return err, &result
}
