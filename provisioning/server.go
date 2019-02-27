package main

import (
	"log"
	"net/http"

	"../provisioning/Routes/GETID"
	"../provisioning/Routes/GenerateToken"
	"../provisioning/Routes/StoreToken"

	"github.com/gorilla/mux"

	_ "net/http/pprof"
)

func StartServer() {
	router := mux.NewRouter()
	go func() {
		log.Fatal(http.ListenAndServe(":6060", http.DefaultServeMux))
	}()
	//http.HandleFunc("/", getIds)
	router.HandleFunc("/storeid/{token}", StoreToken.StoreId)
	router.HandleFunc("/generateToken", generatetoken.GenerateToken).Methods("GET")
	router.HandleFunc("/getid/{id}", getoperation.GetId)

	log.Fatal(http.ListenAndServe(":9001", router))
}
