package serverhandlers

import (
	"encoding/json"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/servercontext"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/pools"
	"net/http"
)

func SavePool(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := types.NewStringResponse("v1beta")
	ctx, err := servercontext.NewAPIServerContext()
	if err != nil {
		d.Error = types.NewContextError(err)
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	p := pools.NewEmptyPool()
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		d.Error = types.NewUnmarshallError(err)
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	err = p.SaveInDB(ctx, false)
	if err != nil {
		d.Error = types.NewFormattedErrorMsg(err, "failed to save pool to DB")
		_ = json.NewEncoder(w).Encode(d)
	}
	d.Data = "Success"
	_ = json.NewEncoder(w).Encode(d)
}
