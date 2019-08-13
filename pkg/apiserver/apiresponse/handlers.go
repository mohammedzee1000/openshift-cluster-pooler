package apiresponse

import (
	"encoding/json"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"net/http"
)

func ListPoolNames(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	d := NewApiResponse()
	ctx, err := generic.NewContext("clientapiserver")
	if err != nil {
		d.Error = NewContextError()
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
	}
	p, err := pools.List(ctx, false)
	if err != nil {
		d.Error = NewListError(err.Error())
		d.Data = nil
		_ = json.NewEncoder(w).Encode(d)
	}
	var poolnamelist []string
	for _, item := range p {
		poolnamelist = append(poolnamelist, item.Name)
	}
	d.Error = NewNoError()
	d.Data = poolnamelist
	_ = json.NewEncoder(w).Encode(poolnamelist)
}