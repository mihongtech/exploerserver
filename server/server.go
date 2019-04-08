package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/linkchain/common/util/log"

	"github.com/linkchain-explorer/server/pool"
	"github.com/linkchain-explorer/server/resp"
)

type Server struct {
}

func NewServer() *Server {
	server := &Server{}
	return server
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    ":9100",
		Handler: mux,
	}

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Content-Type", "application/json;charset=utf-8")
		if strings.ToUpper(request.Method) == "OPTIONS" {
			writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
			return
		}
		if strings.ToUpper(request.Method) != "POST" {
			errorResult(writer, resp.MethodNotAllowedErr)
			return
		}
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			errorResult(writer, resp.InternalServerErr)
			return
		}
		defer request.Body.Close()

		params, e := getParams(request.URL.Path, body)
		if e != nil {
			errorResult(writer, e)
			return
		}

		handler, ok := pool.Handler[request.URL.Path]
		if !ok {
			errorResult(writer, resp.BadRequestErr)
			return
		}
		res, err := handler(params)
		if err != nil {
			errorResult(writer, err)
			return
		}

		result(writer, res)
	})

	log.Info(fmt.Sprintf("Server start at %s", httpServer.Addr))
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
	}
}

func getParams(path string, s []byte) (interface{}, error) {
	p, ok := pool.Params[path]
	if !ok {
		return nil, nil
	}
	params := reflect.New(p.Elem()).Interface()
	err := json.Unmarshal(s, params)
	if err != nil {
		return nil, resp.BadRequestErr
	}
	return params, nil
}

func getHandler(path string) (pool.RequestHandler, error) {
	handler, ok := pool.Handler[path]
	if !ok {
		return nil, resp.BadRequestErr
	}
	return handler, nil
}

func errorResult(writer http.ResponseWriter, err error) {
	result(writer, err)
}

func result(writer http.ResponseWriter, result interface{}) {
	re, _ := json.Marshal(result)
	writer.Write(re)
}
