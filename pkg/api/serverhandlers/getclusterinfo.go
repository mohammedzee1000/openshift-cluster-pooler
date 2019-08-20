package serverhandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func GetClusterInfo(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	d := types.NewClusterInfo("v1beta")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = types.NewContextError(err)
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	vars := mux.Vars(r)
	if val, ok := vars["clusterid"]; ok {
		found := false
		cl, err := clusters.List(ctx)
		if err != nil {
			d.Error = types.NewFormattedErrorMsg(err, "unable to list clusters")
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
			d.Error = types.NewFormattedErrorMsg(nil, "could not find cluster with uuid")
		}
	} else {
		d.Error = types.NewMissingParameterError("clusterid")
	}
	_ = json.NewEncoder(w).Encode(d)
}

