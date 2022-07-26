package cmd

import (
	"github.com/Luukvdm/operator-template/src/controllers"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"log"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
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

	scheme := runtime.NewScheme()

	logger := createLogger(args)

	inCluster := args.insideCluster
	clusterCnf := createConfig(inCluster)

	mgr, err := ctrl.NewManager(clusterCnf, ctrl.Options{
		Scheme: scheme,
		// MetricsBindAddress:     metricsAddr,
		Port: 9445,
		// HealthProbeBindAddress: probeAddr,
		LeaderElection: false,
		// LeaderElectionID:       "",
		Logger: logger,
	})
	if err != nil {
		logger.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.MyReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		mgr.GetLogger().Error(err, "failed to create controller")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healtz", healthz.Ping); err != nil {
		mgr.GetLogger().Error(err, "failed to set up health check")
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		mgr.GetLogger().Error(err, "failed to set up ready check")
		os.Exit(1)
	}

	log.Println("starting controllers")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		mgr.GetLogger().Error(err, "failed to start manager")
		os.Exit(1)
	}
}

func createConfig(inCluster bool) *rest.Config {
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

	return cnf
}

func createLogger(args Args) logr.Logger {
	zLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to instantiate logger:\n%s", err)
	}
	defer zLogger.Sync()
	return zapr.NewLogger(zLogger)
}
