package main

import (
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	scheme             = runtime.NewScheme()
	SchemeGroupVersion = schema.GroupVersion{Group: "suse.tests.dev", Version: "v1"}
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(AddKnownTypes(scheme))
}

func AddKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion,
		&GPUTracker{},
		&GPUTrackerList{},
	)
	metav1.AddToGroupVersion(s, SchemeGroupVersion)
	return nil
}

func main() {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), manager.Options{
		Scheme: scheme,
	})
	if err != nil {
		os.Exit(1)
	}

	if err = (&GPUTrackerReconciler{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		os.Exit(1)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		os.Exit(1)
	}
}
