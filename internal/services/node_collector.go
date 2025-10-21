package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/sentinel-mesh/sentinel-mesh/internal/models"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/metrics"
)

// NodeCollector collects metrics from Kubernetes nodes
type NodeCollector struct {
	kubeClient    kubernetes.Interface
	metricsClient versioned.Interface
	metrics       *metrics.Metrics
	logger        *logrus.Entry
}

// NewNodeCollector creates a new node collector
func NewNodeCollector(kubeClient kubernetes.Interface, metricsClient versioned.Interface, m *metrics.Metrics, logger *logrus.Entry) *NodeCollector {
	return &NodeCollector{
		kubeClient:    kubeClient,
		metricsClient: metricsClient,
		metrics:       m,
		logger:        logger.WithField("collector", "node"),
	}
}

// Collect collects node metrics
func (nc *NodeCollector) Collect(ctx context.Context) ([]*models.Metric, error) {
	start := time.Now()
	
	// Get node list
	nodeList, err := nc.kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}
	
	// Get node metrics
	var nodeMetrics *v1beta1.NodeMetricsList
	nodeMetrics, err = nc.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		nc.logger.WithError(err).Warn("Failed to get node metrics, continuing without usage data")
	}
	
	var metrics []*models.Metric
	timestamp := time.Now()
	
	for _, node := range nodeList.Items {
		nodeMetrics := nc.collectNodeMetrics(&node, nodeMetrics, timestamp)
		metrics = append(metrics, nodeMetrics...)
	}
	
	nc.logger.WithFields(logrus.Fields{
		"node_count":   len(nodeList.Items),
		"metric_count": len(metrics),
		"duration_ms":  time.Since(start).Milliseconds(),
	}).Debug("Collected node metrics")
	
	return metrics, nil
}

// collectNodeMetrics collects metrics for a single node
func (nc *NodeCollector) collectNodeMetrics(node *corev1.Node, nodeMetrics *v1beta1.NodeMetricsList, timestamp time.Time) []*models.Metric {
	var metrics []*models.Metric
	
	labels := map[string]string{
		"node":         node.Name,
		"architecture": node.Status.NodeInfo.Architecture,
		"os":           node.Status.NodeInfo.OperatingSystem,
		"kernel":       node.Status.NodeInfo.KernelVersion,
		"runtime":      node.Status.NodeInfo.ContainerRuntimeVersion,
	}
	
	// Add node labels
	for k, v := range node.Labels {
		labels[fmt.Sprintf("label_%s", k)] = v
	}
	
	// Node capacity metrics
	if cpu := node.Status.Capacity[corev1.ResourceCPU]; !cpu.IsZero() {
		metrics = append(metrics, &models.Metric{
			Name:      "node_cpu_capacity_cores",
			Value:     float64(cpu.MilliValue()) / 1000,
			Labels:    labels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	if memory := node.Status.Capacity[corev1.ResourceMemory]; !memory.IsZero() {
		metrics = append(metrics, &models.Metric{
			Name:      "node_memory_capacity_bytes",
			Value:     float64(memory.Value()),
			Labels:    labels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	if storage := node.Status.Capacity[corev1.ResourceStorage]; !storage.IsZero() {
		metrics = append(metrics, &models.Metric{
			Name:      "node_storage_capacity_bytes",
			Value:     float64(storage.Value()),
			Labels:    labels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	if pods := node.Status.Capacity[corev1.ResourcePods]; !pods.IsZero() {
		metrics = append(metrics, &models.Metric{
			Name:      "node_pod_capacity",
			Value:     float64(pods.Value()),
			Labels:    labels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	// Node allocatable metrics
	if cpu := node.Status.Allocatable[corev1.ResourceCPU]; !cpu.IsZero() {
		metrics = append(metrics, &models.Metric{
			Name:      "node_cpu_allocatable_cores",
			Value:     float64(cpu.MilliValue()) / 1000,
			Labels:    labels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	if memory := node.Status.Allocatable[corev1.ResourceMemory]; !memory.IsZero() {
		metrics = append(metrics, &models.Metric{
			Name:      "node_memory_allocatable_bytes",
			Value:     float64(memory.Value()),
			Labels:    labels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	// Node conditions
	for _, condition := range node.Status.Conditions {
		conditionLabels := make(map[string]string)
		for k, v := range labels {
			conditionLabels[k] = v
		}
		conditionLabels["condition"] = string(condition.Type)
		conditionLabels["status"] = string(condition.Status)
		conditionLabels["reason"] = condition.Reason
		
		var value float64
		if condition.Status == corev1.ConditionTrue {
			value = 1
		} else {
			value = 0
		}
		
		metrics = append(metrics, &models.Metric{
			Name:      "node_condition",
			Value:     value,
			Labels:    conditionLabels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	// Node usage metrics (if available)
	if nodeMetrics != nil {
		for _, nodeMetric := range nodeMetrics.Items {
			if nodeMetric.Name == node.Name {
				// CPU usage
				if cpu := nodeMetric.Usage[corev1.ResourceCPU]; !cpu.IsZero() {
					metrics = append(metrics, &models.Metric{
						Name:      "node_cpu_usage_cores",
						Value:     float64(cpu.MilliValue()) / 1000,
						Labels:    labels,
						Timestamp: timestamp,
						Source:    "kubernetes",
						Type:      models.MetricTypeGauge,
					})
				}
				
				// Memory usage
				if memory := nodeMetric.Usage[corev1.ResourceMemory]; !memory.IsZero() {
					metrics = append(metrics, &models.Metric{
						Name:      "node_memory_usage_bytes",
						Value:     float64(memory.Value()),
						Labels:    labels,
						Timestamp: timestamp,
						Source:    "kubernetes",
						Type:      models.MetricTypeGauge,
					})
				}
				break
			}
		}
	}
	
	// Node info metrics
	infoLabels := make(map[string]string)
	for k, v := range labels {
		infoLabels[k] = v
	}
	infoLabels["kubelet_version"] = node.Status.NodeInfo.KubeletVersion
	infoLabels["kube_proxy_version"] = node.Status.NodeInfo.KubeProxyVersion
	infoLabels["container_runtime"] = node.Status.NodeInfo.ContainerRuntimeVersion
	
	metrics = append(metrics, &models.Metric{
		Name:      "node_info",
		Value:     1,
		Labels:    infoLabels,
		Timestamp: timestamp,
		Source:    "kubernetes",
		Type:      models.MetricTypeGauge,
	})
	
	// Node age
	nodeAge := timestamp.Sub(node.CreationTimestamp.Time)
	metrics = append(metrics, &models.Metric{
		Name:      "node_age_seconds",
		Value:     nodeAge.Seconds(),
		Labels:    labels,
		Timestamp: timestamp,
		Source:    "kubernetes",
		Type:      models.MetricTypeGauge,
	})
	
	// Taints
	for _, taint := range node.Spec.Taints {
		taintLabels := make(map[string]string)
		for k, v := range labels {
			taintLabels[k] = v
		}
		taintLabels["taint_key"] = taint.Key
		taintLabels["taint_value"] = taint.Value
		taintLabels["taint_effect"] = string(taint.Effect)
		
		metrics = append(metrics, &models.Metric{
			Name:      "node_taint",
			Value:     1,
			Labels:    taintLabels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	return metrics
}

// getResourceValue safely extracts resource quantity value
func getResourceValue(quantity resource.Quantity) float64 {
	if quantity.IsZero() {
		return 0
	}
	return float64(quantity.Value())
}

// getResourceMilliValue safely extracts resource quantity milli value
func getResourceMilliValue(quantity resource.Quantity) float64 {
	if quantity.IsZero() {
		return 0
	}
	return float64(quantity.MilliValue())
}