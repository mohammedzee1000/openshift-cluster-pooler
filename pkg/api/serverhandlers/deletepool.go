package serverhandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/servercontext"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/pools"
	"net/http"
)

func DeletePool(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := types.NewStringResponse("v1beta")
	ctx, err := servercontext.NewAPIServerContext()
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
		err = p.MarkForRemoval(ctx)
		if err != nil {
			d.Error = types.NewFormattedErrorMsg(err, "failed to mark pool %s for removal", p.Name)
			_ = json.NewEncoder(w).Encode(d)
			return
		}
	} else {
		d.Error = types.NewMissingParameterError("poolname")
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	d.Data = "Success"
	_ = json.NewEncoder(w).Encode(d)
}