//go:build generate

package api

// will generate yaml and deepcopy for all resources in sub folders
//go:generate controller-gen object paths=./...
//go:generate controller-gen rbac:roleName=my-role crd paths="./..." output:dir=../../charts/operator-template/templates
