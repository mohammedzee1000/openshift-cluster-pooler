package servercontext

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
)

func NewAPIServerContext() (*generic.Context, error)  {
	return generic.NewContext("API Server")
}
