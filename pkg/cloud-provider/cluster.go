package cloud_provider

import (
	"context"
	"fmt"

	cloudprovider "k8s.io/cloud-provider"
)

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (bc *Baiducloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// HasClusterID returns true if a ClusterID is required and set
func (bc *Baiducloud) HasClusterID() bool {
	return true
}

// ListClusters lists the names of the available clusters.
func (bc *Baiducloud) ListClusters(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("ListClusters unimplemented")
}

// Master gets back the address (either DNS name or IP address) of the master node for the cluster.
func (bc *Baiducloud) Master(ctx context.Context, clusterName string) (string, error) {
	return "", fmt.Errorf("Master unimplemented")
}
