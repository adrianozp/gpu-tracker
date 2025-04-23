package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GPUTrackerReconciler struct {
	client client.Client
	scheme *runtime.Scheme
}

func getEnvOrDefault(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getTimeEnvOrDefault(key string, defaultVal time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return time.Duration(intVal)
}

func (r *GPUTrackerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var gpuTracker GPUTracker
	if err := r.client.Get(ctx, req.NamespacedName, &gpuTracker); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	labelKey := getEnvOrDefault("LABEL_KEY", "node-type")
	labelValue := getEnvOrDefault("LABEL_VALUE", "gpu-node")

	var nodeList corev1.NodeList
	if err := r.client.List(ctx, &nodeList, client.MatchingLabels{labelKey: labelValue}); err != nil {
		return ctrl.Result{}, err
	}

	var nodeNames []string
	for _, node := range nodeList.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	newValue := strings.Join(nodeNames, ",")
	if gpuTracker.GPUNodes != newValue {
		gpuTracker.GPUNodes = newValue
		if err := r.client.Update(ctx, &gpuTracker); err != nil {
			return ctrl.Result{}, err
		}
	}

	updateSeconds := getTimeEnvOrDefault("UPDATE_SECONDS", 30)
	return ctrl.Result{RequeueAfter: updateSeconds * time.Second}, nil
}

func (r *GPUTrackerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&GPUTracker{}).
		Complete(r)
}
