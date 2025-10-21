package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/sentinel-mesh/sentinel-mesh/internal/models"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/metrics"
)

// PodCollector collects metrics from Kubernetes pods
type PodCollector struct {
	kubeClient    kubernetes.Interface
	metricsClient versioned.Interface
	metrics       *metrics.Metrics
	logger        *logrus.Entry
}

// NewPodCollector creates a new pod collector
func NewPodCollector(kubeClient kubernetes.Interface, metricsClient versioned.Interface, m *metrics.Metrics, logger *logrus.Entry) *PodCollector {
	return &PodCollector{
		kubeClient:    kubeClient,
		metricsClient: metricsClient,
		metrics:       m,
		logger:        logger.WithField("collector", "pod"),
	}
}

// Collect collects pod metrics
func (pc *PodCollector) Collect(ctx context.Context) ([]*models.Metric, error) {
	start := time.Now()
	
	// Get pod list from all namespaces
	podList, err := pc.kubeClient.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	
	// Get pod metrics
	var podMetrics *v1beta1.PodMetricsList
	podMetrics, err = pc.metricsClient.MetricsV1beta1().PodMetricses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		pc.logger.WithError(err).Warn("Failed to get pod metrics, continuing without usage data")
	}
	
	var metrics []*models.Metric
	timestamp := time.Now()
	
	for _, pod := range podList.Items {
		podMetrics := pc.collectPodMetrics(&pod, podMetrics, timestamp)
		metrics = append(metrics, podMetrics...)
	}
	
	pc.logger.WithFields(logrus.Fields{
		"pod_count":    len(podList.Items),
		"metric_count": len(metrics),
		"duration_ms":  time.Since(start).Milliseconds(),
	}).Debug("Collected pod metrics")
	
	return metrics, nil
}

// collectPodMetrics collects metrics for a single pod
func (pc *PodCollector) collectPodMetrics(pod *corev1.Pod, podMetrics *v1beta1.PodMetricsList, timestamp time.Time) []*models.Metric {
	var metrics []*models.Metric
	
	labels := map[string]string{
		"pod":       pod.Name,
		"namespace": pod.Namespace,
		"node":      pod.Spec.NodeName,
		"phase":     string(pod.Status.Phase),
	}
	
	// Add pod labels
	for k, v := range pod.Labels {
		labels[fmt.Sprintf("label_%s", k)] = v
	}
	
	// Add owner information
	if len(pod.OwnerReferences) > 0 {
		owner := pod.OwnerReferences[0]
		labels["owner_kind"] = owner.Kind
		labels["owner_name"] = owner.Name
	}
	
	// Pod phase
	phaseValue := 0.0
	switch pod.Status.Phase {
	case corev1.PodRunning:
		phaseValue = 1.0
	case corev1.PodPending:
		phaseValue = 2.0
	case corev1.PodSucceeded:
		phaseValue = 3.0
	case corev1.PodFailed:
		phaseValue = 4.0
	case corev1.PodUnknown:
		phaseValue = 5.0
	}
	
	metrics = append(metrics, &models.Metric{
		Name:      "pod_phase",
		Value:     phaseValue,
		Labels:    labels,
		Timestamp: timestamp,
		Source:    "kubernetes",
		Type:      models.MetricTypeGauge,
	})
	
	// Pod ready condition
	readyValue := 0.0
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			if condition.Status == corev1.ConditionTrue {
				readyValue = 1.0
			}
			break
		}
	}
	
	metrics = append(metrics, &models.Metric{
		Name:      "pod_ready",
		Value:     readyValue,
		Labels:    labels,
		Timestamp: timestamp,
		Source:    "kubernetes",
		Type:      models.MetricTypeGauge,
	})
	
	// Container metrics
	for _, container := range pod.Spec.Containers {
		containerLabels := make(map[string]string)
		for k, v := range labels {
			containerLabels[k] = v
		}
		containerLabels["container"] = container.Name
		containerLabels["image"] = container.Image
		
		// Container resource requests
		if cpu := container.Resources.Requests[corev1.ResourceCPU]; !cpu.IsZero() {
			metrics = append(metrics, &models.Metric{
				Name:      "container_cpu_request_cores",
				Value:     float64(cpu.MilliValue()) / 1000,
				Labels:    containerLabels,
				Timestamp: timestamp,
				Source:    "kubernetes",
				Type:      models.MetricTypeGauge,
			})
		}
		
		if memory := container.Resources.Requests[corev1.ResourceMemory]; !memory.IsZero() {
			metrics = append(metrics, &models.Metric{
				Name:      "container_memory_request_bytes",
				Value:     float64(memory.Value()),
				Labels:    containerLabels,
				Timestamp: timestamp,
				Source:    "kubernetes",
				Type:      models.MetricTypeGauge,
			})
		}
		
		// Container resource limits
		if cpu := container.Resources.Limits[corev1.ResourceCPU]; !cpu.IsZero() {
			metrics = append(metrics, &models.Metric{
				Name:      "container_cpu_limit_cores",
				Value:     float64(cpu.MilliValue()) / 1000,
				Labels:    containerLabels,
				Timestamp: timestamp,
				Source:    "kubernetes",
				Type:      models.MetricTypeGauge,
			})
		}
		
		if memory := container.Resources.Limits[corev1.ResourceMemory]; !memory.IsZero() {
			metrics = append(metrics, &models.Metric{
				Name:      "container_memory_limit_bytes",
				Value:     float64(memory.Value()),
				Labels:    containerLabels,
				Timestamp: timestamp,
				Source:    "kubernetes",
				Type:      models.MetricTypeGauge,
			})
		}
	}
	
	// Container status metrics
	for _, containerStatus := range pod.Status.ContainerStatuses {
		containerLabels := make(map[string]string)
		for k, v := range labels {
			containerLabels[k] = v
		}
		containerLabels["container"] = containerStatus.Name
		containerLabels["image"] = containerStatus.Image
		containerLabels["image_id"] = containerStatus.ImageID
		
		// Container ready status
		readyValue := 0.0
		if containerStatus.Ready {
			readyValue = 1.0
		}
		
		metrics = append(metrics, &models.Metric{
			Name:      "container_ready",
			Value:     readyValue,
			Labels:    containerLabels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
		
		// Container restart count
		metrics = append(metrics, &models.Metric{
			Name:      "container_restart_count",
			Value:     float64(containerStatus.RestartCount),
			Labels:    containerLabels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeCounter,
		})
		
		// Container state
		stateValue := 0.0
		if containerStatus.State.Running != nil {
			stateValue = 1.0
		} else if containerStatus.State.Waiting != nil {
			stateValue = 2.0
		} else if containerStatus.State.Terminated != nil {
			stateValue = 3.0
		}
		
		metrics = append(metrics, &models.Metric{
			Name:      "container_state",
			Value:     stateValue,
			Labels:    containerLabels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	// Pod usage metrics (if available)
	if podMetrics != nil {
		for _, podMetric := range podMetrics.Items {
			if podMetric.Name == pod.Name && podMetric.Namespace == pod.Namespace {
				// Pod-level CPU and memory usage
				podLabels := make(map[string]string)
				for k, v := range labels {
					podLabels[k] = v
				}
				
				var totalCPU, totalMemory int64
				for _, container := range podMetric.Containers {
					if cpu := container.Usage[corev1.ResourceCPU]; !cpu.IsZero() {
						totalCPU += cpu.MilliValue()
					}
					if memory := container.Usage[corev1.ResourceMemory]; !memory.IsZero() {
						totalMemory += memory.Value()
					}
				}
				
				if totalCPU > 0 {
					metrics = append(metrics, &models.Metric{
						Name:      "pod_cpu_usage_cores",
						Value:     float64(totalCPU) / 1000,
						Labels:    podLabels,
						Timestamp: timestamp,
						Source:    "kubernetes",
						Type:      models.MetricTypeGauge,
					})
				}
				
				if totalMemory > 0 {
					metrics = append(metrics, &models.Metric{
						Name:      "pod_memory_usage_bytes",
						Value:     float64(totalMemory),
						Labels:    podLabels,
						Timestamp: timestamp,
						Source:    "kubernetes",
						Type:      models.MetricTypeGauge,
					})
				}
				
				// Container-level usage metrics
				for _, container := range podMetric.Containers {
					containerLabels := make(map[string]string)
					for k, v := range labels {
						containerLabels[k] = v
					}
					containerLabels["container"] = container.Name
					
					if cpu := container.Usage[corev1.ResourceCPU]; !cpu.IsZero() {
						metrics = append(metrics, &models.Metric{
							Name:      "container_cpu_usage_cores",
							Value:     float64(cpu.MilliValue()) / 1000,
							Labels:    containerLabels,
							Timestamp: timestamp,
							Source:    "kubernetes",
							Type:      models.MetricTypeGauge,
						})
					}
					
					if memory := container.Usage[corev1.ResourceMemory]; !memory.IsZero() {
						metrics = append(metrics, &models.Metric{
							Name:      "container_memory_usage_bytes",
							Value:     float64(memory.Value()),
							Labels:    containerLabels,
							Timestamp: timestamp,
							Source:    "kubernetes",
							Type:      models.MetricTypeGauge,
						})
					}
				}
				break
			}
		}
	}
	
	// Pod age
	podAge := timestamp.Sub(pod.CreationTimestamp.Time)
	metrics = append(metrics, &models.Metric{
		Name:      "pod_age_seconds",
		Value:     podAge.Seconds(),
		Labels:    labels,
		Timestamp: timestamp,
		Source:    "kubernetes",
		Type:      models.MetricTypeGauge,
	})
	
	// Pod conditions
	for _, condition := range pod.Status.Conditions {
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
			Name:      "pod_condition",
			Value:     value,
			Labels:    conditionLabels,
			Timestamp: timestamp,
			Source:    "kubernetes",
			Type:      models.MetricTypeGauge,
		})
	}
	
	return metrics
}