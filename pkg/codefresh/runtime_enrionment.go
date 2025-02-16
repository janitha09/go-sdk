package codefresh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

const (
	KubernetesRunnerType = "kubernetes"
)

type (
	// IRuntimeEnvironmentAPI declers Codefresh runtime environment API
	IRuntimeEnvironmentAPI interface {
		Create(*CreateRuntimeOptions) (*RuntimeEnvironment, error)
		Validate(*ValidateRuntimeOptions) error
		SignCertificate(*SignCertificatesOptions) ([]byte, error)
		Get(string) (*RuntimeEnvironment, error)
		List() ([]*RuntimeEnvironment, error)
		Delete(string) (bool, error)
		Default(string) (bool, error)
	}

	RuntimeEnvironment struct {
		Version               int                   `json:"version"`
		Metadata              RuntimeMetadata       `json:"metadata"`
		Extends               []string              `json:"extends"`
		Description           string                `json:"description"`
		AccountID             string                `json:"accountId"`
		RuntimeScheduler      RuntimeScheduler      `json:"runtimeScheduler"`
		DockerDaemonScheduler DockerDaemonScheduler `json:"dockerDaemonScheduler"`
		Status                struct {
			Message   string    `json:"message"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"status"`
	}

	RuntimeScheduler struct {
		Cluster struct {
			ClusterProvider struct {
				AccountID string `json:"accountId"`
				Selector  string `json:"selector"`
			} `json:"clusterProvider"`
			Namespace string `json:"namespace"`
		} `json:"cluster"`
		UserAccess bool `json:"userAccess"`
		Pvcs       struct {
			Dind struct {
				StorageClassName string `yaml:"storageClassName"`
			} `yaml:"dind"`
		} `yaml:"pvcs"`
	}

	DockerDaemonScheduler struct {
		Cluster struct {
			ClusterProvider struct {
				AccountID string `json:"accountId"`
				Selector  string `json:"selector"`
			} `json:"clusterProvider"`
			Namespace string `json:"namespace"`
		} `json:"cluster"`
		UserAccess bool `json:"userAccess"`
	}

	RuntimeMetadata struct {
		Agent        bool   `json:"agent"`
		Name         string `json:"name"`
		ChangedBy    string `json:"changedBy"`
		CreationTime string `json:"creationTime"`
	}

	CreateRuntimeOptions struct {
		Cluster            string
		Namespace          string
		HasAgent           bool
		StorageClass       string
		RunnerType         string
		DockerDaemonParams string
		NodeSelector       map[string]string
		Annotations        map[string]string
	}

	ValidateRuntimeOptions struct {
		Cluster   string
		Namespace string
	}

	SignCertificatesOptions struct {
		AltName string
		CSR     string
	}

	CreateResponse struct {
		Name string
	}

	runtimeEnvironment struct {
		codefresh *codefresh
	}
)

func newRuntimeEnvironmentAPI(codefresh *codefresh) IRuntimeEnvironmentAPI {
	return &runtimeEnvironment{codefresh}
}

// Create - create Runtime-Environment
func (r *runtimeEnvironment) Create(opt *CreateRuntimeOptions) (*RuntimeEnvironment, error) {
	re := &RuntimeEnvironment{
		Metadata: RuntimeMetadata{
			Name: fmt.Sprintf("%s/%s", opt.Cluster, opt.Namespace),
		},
	}
	body := map[string]interface{}{
		"clusterName":        opt.Cluster,
		"namespace":          opt.Namespace,
		"storageClassName":   opt.StorageClass,
		"runnerType":         opt.RunnerType,
		"dockerDaemonParams": opt.DockerDaemonParams,
		"nodeSelector":       opt.NodeSelector,
		"annotations":        opt.Annotations,
	}
	if opt.HasAgent {
		body["agent"] = true
	}
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/custom_clusters/register",
		method: "POST",
		body:   body,
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 400 {
		return re, nil
	}
	defer resp.Body.Close()
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Error during runtime environment creation, error: %s", string(buffer))
}

func (r *runtimeEnvironment) Validate(opt *ValidateRuntimeOptions) error {
	body := map[string]interface{}{
		"clusterName": opt.Cluster,
		"namespace":   opt.Namespace,
	}
	_, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/custom_clusters/validate",
		method: "POST",
		body:   body,
	})
	return err
}

func (r *runtimeEnvironment) SignCertificate(opt *SignCertificatesOptions) ([]byte, error) {
	body := map[string]interface{}{
		"reqSubjectAltName": opt.AltName,
		"csr":               opt.CSR,
	}
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/custom_clusters/signServerCerts",
		method: "POST",
		body:   body,
	})
	if err != nil {
		return nil, err
	}
	return r.codefresh.getBodyAsBytes(resp)
}

func (r *runtimeEnvironment) Get(name string) (*RuntimeEnvironment, error) {
	re := &RuntimeEnvironment{}
	path := fmt.Sprintf("/api/runtime-environments/%s", url.PathEscape(name))
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   path,
		method: "GET",
		qs: map[string]string{
			"extend": "false",
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	r.codefresh.decodeResponseInto(resp, re)
	return re, nil
}

func (r *runtimeEnvironment) List() ([]*RuntimeEnvironment, error) {
	emptySlice := make([]*RuntimeEnvironment, 0)
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/runtime-environments",
		method: "GET",
	})
	tokensAsBytes, err := r.codefresh.getBodyAsBytes(resp)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(tokensAsBytes, &emptySlice)

	return emptySlice, err
}

func (r *runtimeEnvironment) Delete(name string) (bool, error) {
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/runtime-environments/%s", url.PathEscape(name)),
		method: "DELETE",
	})
	if err != nil {
		return false, err
	}

	if resp.StatusCode < 400 {
		return true, nil
	}
	body, err := r.codefresh.getBodyAsString(resp)
	if err != nil {
		return false, err
	}
	return false, fmt.Errorf(body)
}

func (r *runtimeEnvironment) Default(name string) (bool, error) {
	path := fmt.Sprintf("/api/runtime-environments/default/%s", url.PathEscape(name))
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   path,
		method: "PUT",
	})
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 201 {
		return true, nil
	} else {
		res, err := r.codefresh.getBodyAsString(resp)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("Unknown error: %v", res)
	}
}
