package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"
	"sigs.k8s.io/e2e-framework/third_party/helm"
)

func TestMain(m *testing.M) {
	c, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	curDir := c

	testEnv, _ := env.NewFromFlags()
	kindClusterName := envconf.RandomName("e2e-testing", 17)
	namespace := envconf.RandomName("kind-ns", 16)
	imageName := "operator-template" // TODO build the image

	fmt.Println("setting up kind cluster named " + kindClusterName)
	testEnv.Setup(
		envfuncs.CreateClusterWithConfig(kind.NewProvider(), kindClusterName, "kind-config.yaml", kind.WithImage("kindest/node:v1.22.2")),
		envfuncs.CreateNamespace(namespace),
		envfuncs.LoadImageToCluster(kindClusterName, imageName),
		DeployOperator(m, curDir, namespace),
	)

	testEnv.Finish(
		envfuncs.DeleteNamespace(namespace),
		envfuncs.ExportClusterLogs(kindClusterName, "./logs"),
		envfuncs.DestroyCluster(kindClusterName),
	)
	os.Exit(testEnv.Run(m))
}

func DeployOperator(m *testing.M, curDir, namespace string) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		manager := helm.New(c.KubeconfigFile())
		chart := path.Join(curDir, "../charts/operator-template")
		if err := manager.RunInstall(helm.WithName("operator-template"), helm.WithNamespace(namespace), helm.WithChart(chart) /* helm.WithArgs(args...),*/, helm.WithWait(), helm.WithTimeout("10m")); err != nil {
			var e *exec.ExitError
			if errors.As(err, &e) {
				fmt.Println(string(e.Stderr))
			}
			return ctx, errors.Join(err, errors.New("failed to invoke helm install operation for chart due to an error"))
		}
		return ctx, nil
	}
}
