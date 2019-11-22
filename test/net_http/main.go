package main

import (
	"fmt"

	"github.com/kvu787/util"
)

func main() {
	logger, err := util.NewLogger(util.LoggerOptions{true, "net_http_test.log", nil})
	if err != nil {
		panic(err)
	}

	requests := util.Listen(3000, logger)
	go func() {
		for {
			request := <-requests
			logger.Println(request.Json)
			request.Complete <- 0
			logger.Println("done")
		}
	}()

	for i := 0; i < 5; i++ {
		go func(i int) {
			err = util.PostJson("localhost", 3000, fmt.Sprintf("hello, world %v", i))
			if err != nil {
				panic(err)
			}
		}(i)
	}

	for {
	}
}
