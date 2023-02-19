package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strconv"
	"taskexecutor/data"
	"taskexecutor/logger"

	"github.com/gorilla/mux"
)

// For local (non-Docker) development/testing
const localDatabaseUrl = "mongodb://admin:password@localhost:27018"
const noLimit = -1

type Handler struct {
	Database *Database
}

func NewHandler() *Handler {
	var databaseUrl string

	if os.Getenv("TASKS_DB_ADDRESS") != "" {
		databaseUrl = os.Getenv("TASKS_DB_ADDRESS")
	} else {
		databaseUrl = localDatabaseUrl
	}

	tasksDb, err := Connect(databaseUrl)

	if err != nil {
		logger.Error("Couldn't connect to database. Reason:" + err.Error())
	}

	handler := new(Handler)
	handler.Database = &tasksDb

	return handler
}

func (h *Handler) GetTasks(w http.ResponseWriter, req *http.Request) {
	logger.Info("Fetching tasks")

	ctx := req.Context()

	tasks, err := h.Database.GetAllTasks(&ctx)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while querying tasks. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Warn(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while querying tasks. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}
}

func (h *Handler) GetTaskById(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	vars := mux.Vars(req)
	key := vars["id"]

	logger.Info("Fetching task with ID " + key)

	task, err := h.Database.GetTaskById(&ctx, key)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while querying tasks by ID. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}

	if task.Empty() {
		warningMessage := fmt.Sprintf("%d - Couldn't not find task by ID %s", http.StatusBadRequest, key)
		logger.Warn(warningMessage)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(warningMessage))

		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while encoding fetched task data by ID. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}
}

func (h *Handler) QueryTasksByStatus(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	status := req.URL.Query().Get("status")

	limitNumber := noLimit
	limitString := req.URL.Query().Get("limit")

	if limitString != "" {
		limitNumber, _ = strconv.Atoi(limitString)
	}

	logger.Info(fmt.Sprintf("Fetching tasks with status %s", status))

	queriedTasks, err := h.Database.QueryTaskByStatus(&ctx, status, limitNumber)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while fetching task data by status. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
	}

	err = json.NewEncoder(w).Encode(queriedTasks)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while encoding fetched task data by status. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}
}

func (h *Handler) UpdateTask(w http.ResponseWriter, req *http.Request) {
	logger.Info("Updating task")

	ctx := req.Context()
	w.Header().Set("Content-Type", "application/json")

	var updatedTask data.Task
	err := json.NewDecoder(req.Body).Decode(&updatedTask)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while updating decoding the task with ID %s. Reason: %s", http.StatusInternalServerError, updatedTask.Id, err)
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}

	err = h.Database.UpdateTask(&ctx, &updatedTask)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while updating updating the task with ID %s. Reason: %s", http.StatusInternalServerError, updatedTask.Id, err)
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}

	err = json.NewEncoder(w).Encode(&updatedTask)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Error while updating encoding the task with ID %s. Reason: %s", http.StatusInternalServerError, updatedTask.Id, err)
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	w.Header().Set("Content-Type", "application/json")

	var newTask data.Task
	err := json.NewDecoder(req.Body).Decode(&newTask)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Could not decode the request while creating a new task. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}

	// Generate and fetch system data (i.e. Id)
	newTask.Id = uuid.New().String()
	newTask.Status = "InQueue"

	// Append request into database
	err = h.Database.InsertTask(&ctx, &newTask)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Could not insert new task into the database. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}

	err = json.NewEncoder(w).Encode(&newTask)
	if err != nil {
		errorMessage := fmt.Sprintf("%d - Could not encode the request while creating a new task. Reason: %s", http.StatusInternalServerError, err.Error())
		logger.Error(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))

		return
	}
}
