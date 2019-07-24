package main

import (
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/cli/clientapiserver"
	"log"
	"net/http"
)

func main()  {
	router := mux.NewRouter()
	router.HandleFunc("/pools/list", clientapiserver.ListPools).Methods("GET")
	//TODO add route here
	log.Fatal(http.ListenAndServe(":20000", router))
}
