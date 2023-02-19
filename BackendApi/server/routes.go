package server

import (
	"github.com/gorilla/mux"
)

type TaskexecutorRoutes struct {
	Handler *Handler
}

func InitializeRoutes() TaskexecutorRoutes {
	return TaskexecutorRoutes{Handler: NewHandler()}
}

func (tr TaskexecutorRoutes) HandleRequests() *mux.Router {
	r := mux.NewRouter()

	// GET
	r.Path("/tasks").Methods("GET").Queries("status", "{status}", "limit", "{limit:[0-9]+}").HandlerFunc(tr.Handler.QueryTasksByStatus)
	r.Path("/tasks").Methods("GET").Queries("status", "{status}").HandlerFunc(tr.Handler.QueryTasksByStatus)
	r.Path("/tasks").Methods("GET").HandlerFunc(tr.Handler.GetTasks)
	r.Path("/tasks/{id}").Methods("GET").HandlerFunc(tr.Handler.GetTaskById)

	// POST
	r.Path("/tasks").Methods("POST").HandlerFunc(tr.Handler.CreateTask)

	// PUT
	r.Path("/tasks").Methods("PUT").HandlerFunc(tr.Handler.UpdateTask)

	return r
}
