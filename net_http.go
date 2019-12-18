package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

var upgrader = websocket.Upgrader{
	// Allow any origin
	CheckOrigin: func(request *http.Request) bool { return true },
}

type JsonRequest struct {
	Json     interface{}
	Response http.ResponseWriter
	complete chan interface{}
}

func (jsonRequest JsonRequest) Complete() {
	close(jsonRequest.complete)
}

func (jsonRequest JsonRequest) RespondJson(object interface{}) error {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	_, err = io.Copy(jsonRequest.Response, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	return nil
}

func PostJson(host string, port uint, object interface{}) (*http.Response, error) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(fmt.Sprintf("http://%v:%v", host, port), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	return response, err
}

// Starts one goroutine to listen for requests.
// Starts one goroutine to handle each request. Exits when Complete is called.
//
// You can receive requests concurrently.
// You can complete requests concurrently and in any order.
func Listen(port uint, logger *log.Logger) chan JsonRequest {
	requests := make(chan JsonRequest)

	go func() {
		serveMux := http.NewServeMux()
		serveMux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
			if request.Method != "POST" {
				logger.Printf("Only POST requests are allowed, got %v", request.Method)
				return
			}

			contentType := request.Header.Get("Content-Type")
			if contentType != "application/json; charset=UTF-8" {
				logger.Print("Content-Type header must be 'application/json; charset=UTF-8', got %v, contentType")
				return
			}

			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				logger.Print(err)
				return
			}

			var jsonObject interface{}
			err = json.Unmarshal(body, &jsonObject)
			if err != nil {
				logger.Print(err)
				return
			}

			complete := make(chan interface{})
			requests <- JsonRequest{jsonObject, response, complete}
			<-complete
		})

		corsHandler := cors.Default().Handler(serveMux)
		server := http.Server{Addr: fmt.Sprintf(":%v", port), Handler: corsHandler}
		err := server.ListenAndServe()
		if err != nil {
			logger.Print(err)
		}
		close(requests)
	}()

	return requests
}

type WebSocket struct {
	Conn     *websocket.Conn
	complete chan interface{}
}

func (webSocket *WebSocket) Complete() {
	close(webSocket.complete)
}

// Starts one goroutine to listen for sockets.
// Starts one goroutine to handle each socket. Exits when Complete is called.
//
// You can receive and use sockets concurrently.
func ListenWebSocket(port uint, logger *log.Logger) chan *WebSocket {
	connections := make(chan *WebSocket)

	go func() {
		serveMux := http.NewServeMux()
		serveMux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
			conn, err := upgrader.Upgrade(response, request, nil)
			if err != nil {
				logger.Print(err)
				return
			}

			complete := make(chan interface{})
			connections <- &WebSocket{conn, complete}
			<-complete
		})

		corsHandler := cors.Default().Handler(serveMux)
		server := http.Server{Addr: fmt.Sprintf(":%v", port), Handler: corsHandler}
		err := server.ListenAndServe()
		if err != nil {
			logger.Print(err)
		}
		close(connections)
	}()

	return connections
}
