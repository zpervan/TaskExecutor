package main

import (
	"fmt"
	"os/exec"
	"sync/atomic"
	"taskexecutor/agent/tasks"
	"taskexecutor/agent/util"
	"taskexecutor/data"
	"taskexecutor/logger"
	"time"
)

// Configuration
const sleepDuration = 2 * time.Second
const maxNumOfConcurrentTasks = 3

func main() {
	shell, arg := util.DetermineShell()
	var activeTasksCount uint32 = 0
	var tasksToExecute []data.Task = make([]data.Task, 3)
	var taskRequests = tasks.NewHttpClient()


	for {
		time.Sleep(sleepDuration)

		taskRequests.UpdateMany(&tasksToExecute)		

		if activeTasksCount >= maxNumOfConcurrentTasks {
			logger.Info("Queue is currently full")
			continue
		}

		queuedTasks := taskRequests.FetchTasks(maxNumOfConcurrentTasks)

		// Add new tasks to the queue in an empty slot
		for _, task := range queuedTasks {
			for i, slot := range tasksToExecute {
				// If the queue slot is empty, assign a task
				if slot.Empty() {
					tasksToExecute[i] = task
					atomic.AddUint32(&activeTasksCount, 1)
					break
				}
			}
		}

		for i := 0; i < maxNumOfConcurrentTasks; i++ {
			// Execute only waiting tasks
			if tasksToExecute[i].Status != "InQueue" {
				continue
			}

			go func(data *data.Task, i int) {
				data.StartedAt = util.CurrentDateTime()
				data.Status = "Executing"

				taskRequests.UpdateSingle(data)

				logger.Info(fmt.Sprintf("Starting worker thread %d at %s", i, data.StartedAt))

				output, err := exec.Command(shell, arg, data.Command).Output()
				if err != nil {
					logger.Error("Could not execute command. Reason: " + err.Error())
					data.StdErr = err.Error()
				}

				data.StdOut = string(output)
				data.FinishedAt = util.CurrentDateTime()
				data.Status = "Finished"

				atomic.AddUint32(&activeTasksCount, ^uint32(0))
				logger.Info("Finishing command execution at " + data.FinishedAt)
			}(&tasksToExecute[i], i)
		}
	}
}
