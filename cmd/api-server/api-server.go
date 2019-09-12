package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/middleware"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/serverhandlers"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func RoutePath(admin bool, version string, pathcomponents ...string) string  {
	var actualpathcomponents []string
	var adminStr string
	if admin {
		adminStr = "/admin"
	}
	actualpathcomponents = append(actualpathcomponents, pathcomponents...)
	p := strings.Join(actualpathcomponents, "/")
	return fmt.Sprintf("/%s%s/%s", version, adminStr, p)
}

func main()  {
	addr := os.Getenv("HOST_ON")
	if len(addr) <= 0 {
		addr = ":20000"
	}
	router := mux.NewRouter()
	RoutePath(false, "v1beta", "cluster", "{clusterid}", "return")
	router.HandleFunc(RoutePath(false, "v1beta", "pools", "list"), serverhandlers.ListPoolNames).Methods("GET")
	router.HandleFunc(RoutePath(false, "v1beta", "pool", "{poolname}", "get-cluster"), serverhandlers.ActivateCluster).Methods("GET")
	router.HandleFunc(RoutePath(false, "v1beta", "pool", "{poolname}", "short-describe"), serverhandlers.GetPoolShortDescription).Methods("GET")
	router.HandleFunc(RoutePath(false, "v1beta", "cluster", "{clusterid}", "describe"), serverhandlers.GetClusterInfo).Methods("GET")
	router.HandleFunc(RoutePath(false, "v1beta", "cluster", "{clusterid}", "return"), serverhandlers.DeactivateCluster).Methods("GET")
	router.HandleFunc(RoutePath(true, "v1beta", "pools", "save"), serverhandlers.SavePool)
	router.HandleFunc(RoutePath(true, "v1beta", "pools", "{poolname}","delete"), serverhandlers.DeletePool)
	router.Use(middleware.LoggingMiddleware)
	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
	srv := http.Server{
		Addr:              addr,
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
