package serverhandlers

import (
	"encoding/json"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/types"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func ListPoolNames(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := types.NewPoolNameList("v1beta")
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = types.NewContextError(err)
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
		return
	}
	p, err := pools.List(ctx, false)
	if err != nil {
		d.Error = types.NewFormattedErrorMsg(err, "failed to list pool names")
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