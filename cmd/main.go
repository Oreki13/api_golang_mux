package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	api "restapi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sign := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		_ = <-sign
		done <- true
	}()

	httpHanler := api.StartHttp()
	srv := &http.Server{
		Handler:      httpHanler,
		Addr:         ":3000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	log.Println("api started")
	<-done
	srv.Shutdown(context.Background())
	log.Println("api stopped")

}
