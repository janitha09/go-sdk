package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IRuntimeAPI interface {
		List(ctx context.Context) ([]model.Runtime, error)
		Create(ctx context.Context, runtimeName, cluster, runtimeVersion string) (*model.RuntimeCreationResponse, error)
	}

	argoRuntime struct {
		codefresh *codefresh
	}

	graphqlRuntimesResponse struct {
		Data struct {
			Runtimes model.RuntimePage
		}
		Errors []graphqlError
	}

	graphQlRuntimeCreationResponse struct {
		Data struct {
			Runtime model.RuntimeCreationResponse
		}
		Errors []graphqlError
	}
)

func newArgoRuntimeAPI(codefresh *codefresh) IRuntimeAPI {
	return &argoRuntime{codefresh: codefresh}
}

func (r *argoRuntime) List(ctx context.Context) ([]model.Runtime, error) {
	jsonData := map[string]interface{}{
		"query": `{
				runtimes {
					edges {
						node {
							metadata {
								name
								namespace
							}
							self {
								healthStatus
								version
							}
							cluster
						}
					}
				}
			}`,
	}

	res := &graphqlRuntimesResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting runtime list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	runtimes := make([]model.Runtime, len(res.Data.Runtimes.Edges))
	for i := range res.Data.Runtimes.Edges {
		runtimes[i] = *res.Data.Runtimes.Edges[i].Node
	}

	return runtimes, nil
}

func (r *argoRuntime) Create(ctx context.Context, runtimeName, cluster, runtimeVersion string) (*model.RuntimeCreationResponse, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation CreateRuntime(
				$name: String!
				$cluster: String!
				$runtimeVersion: String!
			) {
				runtime(name: $name, cluster: $cluster, runtimeVersion: $runtimeVersion) {
					name
					newAccessToken
				}
			}
		`,
		"variables": map[string]interface{}{
			"name":           runtimeName,
			"cluster":        cluster,
			"runtimeVersion": runtimeVersion,
		},
	}

	res := &graphQlRuntimeCreationResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting runtime list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return &res.Data.Runtime, nil
}
