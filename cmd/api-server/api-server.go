package main

import (
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/serverhandlers"
	"log"
	"net/http"
	"os"
	"time"
)

func main()  {
	addr := os.Getenv("HOST_ON")
	if len(addr) <= 0 {
		addr = ":20000"
	}
	router := mux.NewRouter()
	router.HandleFunc("/pools/list", serverhandlers.ListPoolNames).Methods("GET")
	router.HandleFunc("/pool/{poolname}/get-cluster", serverhandlers.ActivateCluster).Methods("GET")
	router.HandleFunc("/pool/{poolname}/short-describe", serverhandlers.GetPoolShortDescription).Methods("GET")
	router.HandleFunc("/cluster/{clusterid}/describe", serverhandlers.GetClusterInfo).Methods("GET")
	router.HandleFunc("/cluster/{clusterid}/return", serverhandlers.DeactivateCluster).Methods("GET")
	router.Use(serverhandlers.LoggingMiddleware)
	srv := http.Server{
		Addr:              addr,
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
