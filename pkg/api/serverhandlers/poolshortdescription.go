package serverhandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func GetPoolShortDescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := types.NewPoolShortDescription("v1beta")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = types.NewContextError(err)
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
		p.GetClusters()
		d.Description = p.Description
	} else {
		d.Error = types.NewMissingParameterError("poolname")
	}
	_ = json.NewEncoder(w).Encode(d)
}