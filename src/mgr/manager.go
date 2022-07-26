package mgr

import (
	"github.com/Luukvdm/operator-template/src/api/v1alpha1"
	"github.com/Luukvdm/operator-template/src/controllers"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"log"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

func StartManager(cnf *rest.Config) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	logger := createLogger()

	mgr, err := ctrl.NewManager(cnf, ctrl.Options{
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

	if err = (&controllers.MyResourceReconciler{
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

func createLogger() logr.Logger {
	zLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to instantiate logger:\n%s", err)
	}
	defer zLogger.Sync()
	return zapr.NewLogger(zLogger)
}
