package server

import (
	"backend-go/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

/* 서버를 정의하고 용처과 응답 구조를 정의 */

func NewHTTPServer(addr string) *http.Server {
	httpsServer := newHTTPServer()

	r := mux.NewRouter()

	r.HandleFunc("/", httpsServer.handleProduce).Methods("Post")
	r.HandleFunc("/", httpsServer.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *log.Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: log.NewLog(),
	}
}

type ProduceRequest struct {
	Record log.Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record log.Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest

	/* 1-1 : request decode */
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/* 2-1 : log data add */
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/* 3-1 : log data offset*/
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/* custom error object 활용 - 404 code 와 sync */
	record, err := s.Log.Read(req.Offset)
	if err == log.ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ConsumeResponse{Record: record}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
