package middleware

import (
	"io"
	"net/http"
	"pgpro2024/internal/service"
	"strconv"
)

type MyRouter struct {
	Service service.ServiceInterface
	Server  *http.Server
}

func NewMyRouter(host string, port string, service service.ServiceInterface) *MyRouter {
	r := &MyRouter{}
	r.Service = service
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", r.ping)
	mux.HandleFunc("GET /get/{id}", r.getAvailibleCommand)
	mux.HandleFunc("GET /getAll", r.getListAvailibleCommands)
	mux.HandleFunc("PATCH /create", r.createCommand)
	mux.HandleFunc("POST /start", r.startCommand)
	mux.HandleFunc("POST /status", r.getStatusPID)
	mux.HandleFunc("DELETE /stop", r.stopPID)
	r.Server = &http.Server{
		Addr:    host + ":" + port,
		Handler: mux,
	}
	return r
}

func (r *MyRouter) getAvailibleCommand(response http.ResponseWriter, request *http.Request) {
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

func (r *MyRouter) getListAvailibleCommands(response http.ResponseWriter, request *http.Request) {
	data, e := r.Service.GetListAvailibleCommands()
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(data)
}

func (r *MyRouter) createCommand(response http.ResponseWriter, request *http.Request) {
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

func (r *MyRouter) startCommand(response http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if len(data) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ans, e := r.Service.StartCommand(data)
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
	response.Write(ans)
}

func (r *MyRouter) getStatusPID(response http.ResponseWriter, request *http.Request) {
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

func (r *MyRouter) stopPID(response http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if len(data) == 0 || err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	e := r.Service.StopProcess(data)
	if e.E != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(e.Err)
		return
	}
}

func (r *MyRouter) ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
}

func (r *MyRouter) ListenAndServe() {
	r.Server.ListenAndServe()
}
