package main

import (
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/apiserver/apiresponse"
	"log"
	"net/http"
)

func main()  {
	router := mux.NewRouter()
	router.HandleFunc("/pools/list", apiresponse.ListPoolNames).Methods("GET")
	//TODO add route here
	log.Fatal(http.ListenAndServe(":20000", router))
}
