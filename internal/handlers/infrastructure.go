package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// InfrastructureHandler handles infrastructure-related requests
type InfrastructureHandler struct {
	clientset *kubernetes.Clientset
}

// NodeInfo represents node information
type NodeInfo struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	CPU        int    `json:"cpu"`
	Memory     int    `json:"memory"`
	Pods       int    `json:"pods"`
	K8sVersion string `json:"k8sVersion"`
}

// ServiceInfo represents service information
type ServiceInfo struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Replicas   string `json:"replicas"`
	Port       string `json:"port"`
	Version    string `json:"version"`
	Deployment string `json:"deployment"`
	Namespace  string `json:"namespace"`
}

// InfrastructureResponse represents the full infrastructure status
type InfrastructureResponse struct {
	Nodes    []NodeInfo    `json:"nodes"`
	Services []ServiceInfo `json:"services"`
	Timestamp string        `json:"timestamp"`
}

// NewInfrastructureHandler creates a new infrastructure handler
func NewInfrastructureHandler() (*InfrastructureHandler, error) {
	// Create in-cluster client config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &InfrastructureHandler{
		clientset: clientset,
	}, nil
}

// GetInfrastructure returns current infrastructure status
func (h *InfrastructureHandler) GetInfrastructure(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Fetch nodes
	nodes, err := h.fetchNodes(ctx)
	if err != nil {
		// Return fallback data if unable to fetch
		nodes = []NodeInfo{{
			Name:   "Unable to fetch nodes",
			Status: "Unknown",
			CPU:    0,
			Memory: 0,
			Pods:   0,
		}}
	}

	// Fetch services
	services, err := h.fetchServices(ctx)
	if err != nil {
		// Return fallback data if unable to fetch
		services = []ServiceInfo{}
	}

	response := InfrastructureResponse{
		Nodes:     nodes,
		Services:  services,
		Timestamp: getTodayDate(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// fetchNodes retrieves node information from the cluster
func (h *InfrastructureHandler) fetchNodes(ctx context.Context) ([]NodeInfo, error) {
	nodes, err := h.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeInfos []NodeInfo
	for _, node := range nodes.Items {
		status := "Ready"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
				status = "NotReady"
				break
			}
		}

		nodeInfo := NodeInfo{
			Name:       node.Name,
			Status:     status,
			CPU:        0, // Would need metrics server for real CPU data
			Memory:     0, // Would need metrics server for real memory data
			Pods:       len(node.Status.Images), // Rough estimate
			K8sVersion: node.Status.NodeInfo.KubeletVersion,
		}
		nodeInfos = append(nodeInfos, nodeInfo)
	}

	return nodeInfos, nil
}

// fetchServices retrieves service information from the cluster
func (h *InfrastructureHandler) fetchServices(ctx context.Context) ([]ServiceInfo, error) {
	deployments, err := h.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var services []ServiceInfo
	for _, deployment := range deployments.Items {
		status := "Running"
		if deployment.Status.Replicas == 0 {
			status = "Down"
		}

		replicas := "0/0"
		if deployment.Spec.Replicas != nil {
			replicas = fmt.Sprintf("%d/%d", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
		}

		// Get version from image tag if available
		version := "v1.0.0"
		if len(deployment.Spec.Template.Spec.Containers) > 0 {
			image := deployment.Spec.Template.Spec.Containers[0].Image
			if idx := strings.LastIndex(image, ":"); idx != -1 {
				version = image[idx+1:]
			}
		}

		service := ServiceInfo{
			Name:       deployment.Name,
			Status:     status,
			Replicas:   replicas,
			Port:       "8080",
			Version:    version,
			Deployment: deployment.Name,
			Namespace:  deployment.Namespace,
		}
		services = append(services, service)
	}

	return services, nil
}

// getTodayDate returns today's date as a string
func getTodayDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
