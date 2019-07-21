package main

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/cli/poolmanager"
	"log"
)

func main()  {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	poolmanager.NewPoolManager().Run()
}