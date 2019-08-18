package common

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/apiserver/apiresponse"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func ListPoolNames(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := apiresponse.NewApiResponse("v1beta")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = apiresponse.NewContextError(err)
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	p, err := pools.List(ctx, false)
	if err != nil {
		d.Error = apiresponse.NewFormattedErrorMsg(err, "failed to list pool names")
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	var poolnamelist []string
	for _, item := range p {
		poolnamelist = append(poolnamelist, item.Name)
	}
	d.Error = ""
	d.Data = poolnamelist
	_ = json.NewEncoder(w).Encode(d)
}

func ActivateCluster(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := apiresponse.NewApiResponse("v1beta")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = apiresponse.NewContextError(err)
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	vars := mux.Vars(r)
	if val, ok := vars["poolname"]; ok {
		p, err := pools.PoolByName(ctx, val, false)
		if err != nil {
			d.Error = apiresponse.NewFormattedErrorMsg(err, "failed to find pool named %s", val)
			_ = json.NewEncoder(w).Encode(d)
			return
		}
		c, err := p.ActivateCluster(ctx)
		if err != nil {
			d.Error = apiresponse.NewFormattedErrorMsg(err, "could not access cluster information")
		}
		d.Data = c
	} else {
		d.Error = apiresponse.NewMissingParameterError("poolname")
	}
	_ = json.NewEncoder(w).Encode(d)
}