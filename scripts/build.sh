#!/usr/bin/env bash
set +eux
go build cmd/test-db-cli/db-cli.go
go build cmd/pool-manager/pool-manager.go
go build cmd/api-server/api-server.go
