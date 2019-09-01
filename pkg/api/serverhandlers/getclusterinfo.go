package serverhandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/servercontext"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"net/http"
)

func GetClusterInfo(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	d := types.NewClusterInfo("v1beta")
	_, _ = servercontext.NewAPIServerContext()
	ctx, err := servercontext.NewAPIServerContext()
	if err != nil {
		d.Error = types.NewContextError(err)
		d.Cluster = nil
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
				d.Cluster = cl.ItemAt(i)
				p, err := pools.PoolByName(ctx, d.Cluster.PoolName, false)
				if err != nil {
					d.Error = types.NewFormattedErrorMsg(err, "unable to get pool %s to which cluster belongs", d.Cluster.PoolName)
				}
				d.ExpiresOn = p.ExpiresOn(d.Cluster)
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

