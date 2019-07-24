package client_api_server

import (
	"encoding/json"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/cli/apierror"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func ActivateCluster(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	ctx, err := generic.NewContext("client-api-server")
	if err != nil {
		_ = json.NewEncoder(w).Encode(apierror.NewContextError())
	}
	
}