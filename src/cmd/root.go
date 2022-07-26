package cmd

import (
	"github.com/Luukvdm/operator-template/src/mgr"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
)

type OperatorCmd struct {
	command *cobra.Command
	Args
}

type Args struct {
	isDebug       bool
	insideCluster bool
}

// Execute starts the base command
func Execute() {
	cmd := OperatorCmd{
		command: &cobra.Command{
			Use:   "otemplate",
			Short: "Template for K8S operators with some example controllers",
			Long: `A template for Kubernetes operators that contains some example controllers.
This template uses:
- Controller-gen (from github.com/kubernetes-sigs/controller-tools/) for generating CRD yaml
- Helm charts
- Kubernetes client-go
- Cobra for the CLI`,
		},
	}

	cmd.command.PersistentFlags().BoolP("debug", "d", cmd.Args.isDebug, "Enable debug mode")
	cmd.command.PersistentFlags().Bool("local-config", cmd.Args.insideCluster, "Use the local kubeconfig instead of getting it from the cluster")

	cmd.command.Run = func(_ *cobra.Command, _ []string) {
		run(cmd.Args)
	}

	if err := cmd.command.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func run(args Args) {
	log.Println("starting operator-template")

	inCluster := args.insideCluster

	var cnf *rest.Config
	var err error
	if inCluster {
		cnf, err = rest.InClusterConfig()
	} else {
		cnf, err = ctrl.GetConfig()
	}

	if err != nil {
		log.Fatalln("failed to create cluster config:\n%w", err)
	}

	mgr.StartManager(cnf)

}
