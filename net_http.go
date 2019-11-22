package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rs/cors"
)

type JsonRequest struct {
	Json     interface{}
	Response http.ResponseWriter
	Complete chan interface{}
}

func PostJson(host string, port uint, object interface{}) error {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = http.Post(fmt.Sprintf("http://%v:%v", host, port), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	return nil
}

// You can receive requests concurrently.
// You can complete requests concurrently and in any order.
func Listen(port uint, logger *log.Logger) chan JsonRequest {
	requests := make(chan JsonRequest)

	go func() {
		serveMux := http.NewServeMux()
		serveMux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
			if request.Method != "POST" {
				logger.Print("Only POST requests are allowed")
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
			jsonRequest := JsonRequest{jsonObject, response, complete}
			requests <- jsonRequest
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
