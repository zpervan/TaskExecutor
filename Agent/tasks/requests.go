package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"taskexecutor/data"
	"taskexecutor/logger"
)

// @TODO: Think of a better struct name 
type Client struct {
	client  *http.Client
	address string
}

func NewHttpClient() *Client {
	logger.Info("Initializing HTTP client for requests")

	requests := new(Client)
	requests.client = &http.Client{}

	if os.Getenv("TASKS_BACKEND_API") != "" {
		requests.address = os.Getenv("TASKS_BACKEND_API")
	} else {
		requests.address = "http://localhost:3500"
	}

	return requests
}

func (r *Client) FetchTasks(numberOfTasks ...int) []data.Task {
	var request string

	if len(numberOfTasks) > 0 && numberOfTasks[0] > 0 {
		request = fmt.Sprintf("%s/tasks?status=InQueue&limit=%d", r.address, numberOfTasks[0])
	} else {
		request = fmt.Sprintf("%s/tasks?status=InQueue", r.address)
	}

	resp, err := http.Get(request)
	if err != nil {
		logger.Error("Could not fetch tasks. Reason: " + err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Could not fetch tasks. Reason: " + err.Error())
	}

	var queuedTasks []data.Task
	err = json.Unmarshal(body, &queuedTasks)
	logger.Info(fmt.Sprintf("Number of queued tasks: %d", len(queuedTasks)))
	if err != nil {
		logger.Error("Could not fetch tasks. Reason: " + err.Error())
	}

	return queuedTasks
}

func (r *Client) UpdateSingle(task *data.Task) error {
	taskJson, err := json.Marshal(task)

	if err != nil {
		logger.Error("Could not update task. Reason: " + err.Error())
		return err
	}

	request, err := http.NewRequest(http.MethodPut, r.address+"/tasks", bytes.NewBuffer(taskJson))
	if err != nil {
		logger.Error("Could not update task. Reason: " + err.Error())
		return err
	}

	_, err = r.client.Do(request)

	if err != nil {
		logger.Error("Could not update task. Reason: " + err.Error())
		return err
	}

	logger.Info(fmt.Sprintf("Task with ID %s updated successfully", task.Id))
	return nil
}

func (r *Client) UpdateMany(tasks *[]data.Task) {
	for i, task := range *tasks {
		if task.Status == "Finished" {
			err := r.UpdateSingle(&task)

			if err != nil {
				logger.Error(fmt.Sprintf("Could not update task. Reason: %s", err))
				continue
			}

			(*tasks)[i] = data.Task{}
		}
	}
}
