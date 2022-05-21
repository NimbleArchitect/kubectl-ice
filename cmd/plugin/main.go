package main

import (
	"github.com/NimbleArchitect/kubectl-ice/cmd/plugin/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // required for GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	cli.InitAndExecute()
}
