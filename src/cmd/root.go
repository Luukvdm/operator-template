package cmd

import (
	"github.com/Luukvdm/operator-template/src/mgr"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/rest"
	"log"
	"os"
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
	l := createZapLogger(args.isDebug)
	l.V(int(zap.InfoLevel)).Info("starting operator-template")
	if args.isDebug {
		l.V(int(zap.DebugLevel)).Info("running in debug mode")
	}

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

	mgr.StartManager(l, cnf)

}

func createZapLogger(isDebug bool) logr.Logger {
	var zLogger *zap.Logger
	var err error
	if isDebug {
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			os.Stderr,
			zap.DebugLevel,
		)
		zLogger = zap.New(core)
	} else {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			os.Stderr,
			zap.InfoLevel,
		)
		zLogger = zap.New(core)
	}
	if err != nil {
		log.Fatalf("failed to create logger:\n%s", err)
	}

	defer zLogger.Sync()

	return zapr.NewLogger(zLogger)
}
