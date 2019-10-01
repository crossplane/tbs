/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// +kubebuilder:scaffold:imports
	linodev1alpha1 "github.com/displague/stack-linode/api/v1alpha1"
	"github.com/displague/stack-linode/controllers"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {

	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme, MetricsBindAddress: metricsAddr})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupLog.Info("Adding schemes")

	// add all resources to the manager's runtime scheme

	if err := addToScheme(mgr.GetScheme()); err != nil {
		setupLog.Error(err, "Cannot add APIs to scheme")
		os.Exit(1)
	}
	setupLog.Info("Adding controllers")

	// Setup all Controllers
	if err := controllerSetupWithManager(mgr); err != nil {
		setupLog.Error(err, "Cannot add controllers to manager")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func controllerSetupWithManager(mgr manager.Manager) error {
	if err := (&controllers.InstanceController{}).SetupWithManager(mgr); err != nil {
		return err
	}

	return nil
}

// addToScheme adds all resources to the runtime scheme.
func addToScheme(scheme *runtime.Scheme) error {
	if err := linodev1alpha1.AddToScheme(scheme); err != nil {
		return err
	}

	if err := corev1.AddToScheme(scheme); err != nil {
		return err
	}

	return nil
}
