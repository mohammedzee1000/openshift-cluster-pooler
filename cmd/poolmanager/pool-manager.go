package main

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/cli/poolmanager"
)

func main() {
	p := poolmanager.NewPoolManager()
	p.Run()
}
