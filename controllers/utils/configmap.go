package utils

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ODFInfoConfigMapName    = "odf-info"
	ConfigMapResourceType   = "ConfigMap"
	ClientInfoConfigMapName = "odf-client-info"
)

// FetchConfigMap fetches a ConfigMap with a given name from a given namespace
func FetchConfigMap(ctx context.Context, c client.Client, name, namespace string) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := c.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, configMap)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to fetch ConfigMap %s in namespace %s: %v", name, namespace, err)
	}
	return configMap, nil
}

// GetODFInfoConfigMap fetches the odf-info ConfigMap from the given namespace. This will only work on the managed cluster
func GetODFInfoConfigMap(ctx context.Context, c client.Client, namespace string) (*corev1.ConfigMap, error) {
	return FetchConfigMap(ctx, c, ODFInfoConfigMapName, namespace)
}

func SplitKeyForNamespacedName(key string) types.NamespacedName {
	// key = openshift-storage_ocs-storagecluster.config.yaml
	splitKey := strings.Split(key, ".")               // [openshift-storage_ocs-storagecluster,config,yaml]
	namespacedName := strings.Split(splitKey[0], "_") // [openshift-storage,ocs-storagecluster]
	return types.NamespacedName{Namespace: namespacedName[0], Name: namespacedName[1]}
}