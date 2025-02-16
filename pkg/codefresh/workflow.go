package codefresh

import (
	"errors"
	"fmt"
	"time"
)

type (
	IWorkflowAPI interface {
		WaitForStatus(string, string, time.Duration, time.Duration) error
		Get(string) (*Workflow, error)
	}

	workflow struct {
		codefresh *codefresh
	}

	Workflow struct {
		ID                 string    `json:"id"`
		Status             string    `json:"status"`
		UserYamlDescriptor string    `json:"userYamlDescriptor"`
		Progress           string    `json:"progress"`
		Created            time.Time `json:"created"`
		Updated            time.Time `json:"updated"`
		Finished           time.Time `json:"finished"`
	}
)

func newWorkflowAPI(codefresh *codefresh) IWorkflowAPI {
	return &workflow{codefresh}
}

func (w *workflow) Get(id string) (*Workflow, error) {
	wf := &Workflow{}
	resp, err := w.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/builds/%s", id),
		method: "GET",
	})
	// failed in api call
	if err != nil {
		return nil, err
	}
	err = w.codefresh.decodeResponseInto(resp, wf)
	// failed to decode
	if err != nil {
		return nil, err
	}
	return wf, nil
}

func (w *workflow) WaitForStatus(id string, status string, interval time.Duration, timeout time.Duration) error {
	err := waitFor(interval, timeout, func() (bool, error) {

		wf := &Workflow{}
		resp, err := w.codefresh.requestAPI(&requestOptions{
			path:   fmt.Sprintf("/api/builds/%s", id),
			method: "GET",
		})
		// failed in api call
		if err != nil {
			return false, err
		}
		err = w.codefresh.decodeResponseInto(resp, wf)
		// failed to decode
		if err != nil {
			return false, err
		}
		// status match
		if wf.Status == status {
			return true, nil
		}
		// status dosent match
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func waitFor(interval time.Duration, timeout time.Duration, execution func() (bool, error)) error {
	t := time.After(timeout)
	tick := time.Tick(interval)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-t:
			return errors.New("timed out")
		case <-tick:
			ok, err := execution()
			if err != nil {
				return err
			} else if ok {
				return nil
			}
		}
	}
}
