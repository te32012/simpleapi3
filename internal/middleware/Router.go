package middleware

import (
	"io"
	"net/http"
	"pgpro2024/internal/service"
	"strconv"
)

type Router struct {
	Service service.ServiceInterface
}

func NewRouter() *Router {
	r := &Router{}
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", r.ping)
	mux.HandleFunc("GET /get/{id}", r.getAvailibleCommand)
	mux.HandleFunc("GET /getAll", r.getListAvailibleCommands)
	mux.HandleFunc("PATCH /create", r.createCommand)
	mux.HandleFunc("POST /start/{id}", r.startCommand)
	mux.HandleFunc("GET /status", r.getStatusPID)
	mux.HandleFunc("DELETE /stop/{id}", r.stopPID)
	return r
}

func (r *Router) getAvailibleCommand(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if len(request.PathValue("id")) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	data, e := r.Service.GetAvailibleCommandById(id)
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(data)
}

func (r *Router) getListAvailibleCommands(response http.ResponseWriter, request *http.Request) {
	data, e := r.Service.GetListAvailibleCommands()
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(data)
}

func (r *Router) createCommand(response http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if len(data) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ans, e := r.Service.CreateCommand(data)
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(ans)
}

func (r *Router) startCommand(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if len(request.PathValue("id")) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(request.Body)
	if len(data) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ans, e := r.Service.StartCommand(id, data)
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(ans)
}

func (r *Router) getStatusPID(response http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if len(data) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ans, e := r.Service.GetStatusProcess(data)
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(ans)
}

func (r *Router) stopPID(response http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if len(data) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ans, e := r.Service.StopProcess(data)
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(ans)
}

func (r *Router) ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
}
