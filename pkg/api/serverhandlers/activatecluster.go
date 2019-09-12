package serverhandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/servercontext"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/pools"
	"net/http"
)

func ActivateCluster(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := types.NewClusterInfo("v1beta")
	ctx, err := servercontext.NewAPIServerContext()
	if err != nil {
		d.Error = types.NewContextError(err)
		d.Cluster = nil
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	vars := mux.Vars(r)
	if val, ok := vars["poolname"]; ok {
		p, err := pools.PoolByName(ctx, val, false)
		if err != nil {
			d.Error = types.NewFormattedErrorMsg(err, "failed to find pool named %s", val)
			_ = json.NewEncoder(w).Encode(d)
			return
		}
		c, err := p.ActivateCluster(ctx)
		if err != nil {
			d.Error = types.NewFormattedErrorMsg(err, "could not access cluster information")
		}
		d.Cluster = c
		d.ExpiresOn = p.ExpiresOn(c)
	} else {
		d.Error = types.NewMissingParameterError("poolname")
	}
	_ = json.NewEncoder(w).Encode(d)
}