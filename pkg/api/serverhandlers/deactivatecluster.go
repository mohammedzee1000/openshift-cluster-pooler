package serverhandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/servercontext"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"net/http"
)

func DeactivateCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	d := types.NewStringResponse("v1beta")
	ctx, err := servercontext.NewAPIServerContext()
	if err != nil {
		d.Error = types.NewContextError(err)
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
		for i := 0; i < cl.Len(); i++ {
			if !found && cl.ItemAt(i).ClusterID == val {
				curr := cl.ItemAt(i)
				found = true
				d.Data = "Found"
				err = curr.Return(ctx)
				if err != nil {
					d.Error = types.NewFormattedErrorMsg(err, "failed to return cluster")
					_ = json.NewEncoder(w).Encode(d)
					return
				}
				d.Data = "OK"
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
	return
}
