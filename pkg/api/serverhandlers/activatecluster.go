package serverhandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func ActivateCluster(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := types.NewClusterInfo("v1beta")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = types.NewContextError(err)
		d.Data = nil
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
		d.Data = c
	} else {
		d.Error = types.NewMissingParameterError("poolname")
	}
	_ = json.NewEncoder(w).Encode(d)
}