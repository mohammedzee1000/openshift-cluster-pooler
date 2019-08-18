package main

import (
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/apiserver/common"
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
	router.HandleFunc("/pools/list", common.ListPoolNames).Methods("GET")
	router.HandleFunc("/pool/{poolname}/activate-cluster", common.ActivateCluster).Methods("GET")
	router.HandleFunc("/cluster/{clusterid}/describe", common.GetClusterInfo).Methods("GET")
	srv := http.Server{
		Addr:              addr,
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
