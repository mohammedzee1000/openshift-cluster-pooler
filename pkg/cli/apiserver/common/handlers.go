package common

import (
	"encoding/json"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/cli/apiserver/apierror"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func ListPools(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		_ = json.NewEncoder(w).Encode(apierror.NewContextError())
	}
	p, err := pools.List(ctx)
	if err != nil {
		_ = json.NewEncoder(w).Encode(apierror.NewListError(err.Error()))
	}
	var poolnamelist []string
	for _, item := range p {
		poolnamelist = append(poolnamelist, item.Name)
	}
	_ = json.NewEncoder(w).Encode(poolnamelist)
}