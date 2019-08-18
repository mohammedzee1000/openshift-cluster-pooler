package common

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/apiserver/apiresponse"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
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

func GetClusterInfo(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	d := apiresponse.NewApiResponse("v1beta")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = apiresponse.NewContextError(err)
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	vars := mux.Vars(r)
	if val, ok := vars["clusterid"]; ok {
		found := false
		cl, err := clusters.List(ctx)
		if err != nil {
			d.Error = apiresponse.NewFormattedErrorMsg(err, "unable to list clusters")
			_ = json.NewEncoder(w).Encode(d)
			return
		}
		for i:=0 ; i < cl.Len(); i++ {
			if !found && cl.ItemAt(i).ClusterID == val {
				found = true
				d.Data = cl.ItemAt(i)
				_ = json.NewEncoder(w).Encode(d)
				return
			}
		}
		if !found {
			d.Error = apiresponse.NewFormattedErrorMsg(nil, "could not find cluster with uuid")
		}
	} else {
		d.Error = apiresponse.NewMissingParameterError("clusterid")
	}
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