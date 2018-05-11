package server

import (
	"net/http"
	"fmt"
	"os"
	"encoding/json"
	"apt-model/jo"
	"time"
)

var server http.Server

func Start() {
	http.HandleFunc("/calculate", requestHandler)
	server = http.Server{Addr:":1337"}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println(err)
			os.Exit(100)
		}
	}()
}

func Stop() error {
	return server.Shutdown(nil)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	request := new(jo.Request)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(request); err != nil {
		fmt.Println(err)
		return // TODO: send response with error
	}
	r.Body.Close()
	response := new(jo.Response)
	select {
	case response = <-calculateRequest(request):
	case <-time.After(time.Duration(request.Timeout) * time.Second):
		fmt.Println("Timeout")
		return // TODO: send response with timeout
	}
	if data, err := json.Marshal(&response); err != nil {
		fmt.Println(err)
		// TODO: send response with error
	} else {
		if _, err := w.Write(data); err != nil {
			fmt.Println(err)
			// TODO: send response with error
		}
	}
}

func calculateRequest(request *jo.Request) <-chan *jo.Response {
	ch := make(chan *jo.Response)
	go func(ch chan<- *jo.Response) {
		nodes := request.Nodes.Parse()
		routes := nodes.Calculate()
		response := jo.Response{Routes:jo.ConvertRoutes(routes)}
		ch <- &response
		close(ch)
	}(ch)
	return ch
}