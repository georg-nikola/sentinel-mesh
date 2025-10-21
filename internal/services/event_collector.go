package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/sentinel-mesh/sentinel-mesh/internal/models"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/metrics"
)

// EventCollector collects events from Kubernetes
type EventCollector struct {
	kubeClient    kubernetes.Interface
	metrics       *metrics.Metrics
	logger        *logrus.Entry
	lastTimestamp time.Time
}

// NewEventCollector creates a new event collector
func NewEventCollector(kubeClient kubernetes.Interface, m *metrics.Metrics, logger *logrus.Entry) *EventCollector {
	return &EventCollector{
		kubeClient:    kubeClient,
		metrics:       m,
		logger:        logger.WithField("collector", "event"),
		lastTimestamp: time.Now().Add(-5 * time.Minute), // Start collecting events from 5 minutes ago
	}
}

// Collect collects Kubernetes events
func (ec *EventCollector) Collect(ctx context.Context) ([]*models.Event, error) {
	start := time.Now()
	
	// Get events from all namespaces
	eventList, err := ec.kubeClient.CoreV1().Events("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	
	var events []*models.Event
	newLastTimestamp := ec.lastTimestamp
	
	for _, event := range eventList.Items {
		// Only collect events newer than our last timestamp
		if event.LastTimestamp.Time.After(ec.lastTimestamp) {
			modelEvent := ec.convertKubernetesEvent(&event)
			events = append(events, modelEvent)
			
			// Update the newest timestamp
			if event.LastTimestamp.Time.After(newLastTimestamp) {
				newLastTimestamp = event.LastTimestamp.Time
			}
		}
	}
	
	// Update last timestamp for next collection
	ec.lastTimestamp = newLastTimestamp
	
	ec.logger.WithFields(logrus.Fields{
		"total_events":     len(eventList.Items),
		"collected_events": len(events),
		"duration_ms":      time.Since(start).Milliseconds(),
	}).Debug("Collected events")
	
	return events, nil
}

// convertKubernetesEvent converts a Kubernetes event to our internal event model
func (ec *EventCollector) convertKubernetesEvent(kubeEvent *corev1.Event) *models.Event {
	event := &models.Event{
		ID:        generateEventID(kubeEvent),
		Type:      kubeEvent.Type,
		Reason:    kubeEvent.Reason,
		Message:   kubeEvent.Message,
		Source:    fmt.Sprintf("%s/%s", kubeEvent.Source.Component, kubeEvent.Source.Host),
		Timestamp: kubeEvent.LastTimestamp.Time,
		Severity:  determineSeverity(kubeEvent),
		Object: models.EventObject{
			Kind:      kubeEvent.InvolvedObject.Kind,
			Name:      kubeEvent.InvolvedObject.Name,
			Namespace: kubeEvent.InvolvedObject.Namespace,
			UID:       string(kubeEvent.InvolvedObject.UID),
		},
		Labels: map[string]string{
			"namespace":  kubeEvent.Namespace,
			"component":  kubeEvent.Source.Component,
			"host":       kubeEvent.Source.Host,
			"kind":       kubeEvent.InvolvedObject.Kind,
			"name":       kubeEvent.InvolvedObject.Name,
			"reason":     kubeEvent.Reason,
			"type":       kubeEvent.Type,
		},
	}
	
	// Add resource version and field path if available
	if kubeEvent.InvolvedObject.ResourceVersion != "" {
		event.Labels["resource_version"] = kubeEvent.InvolvedObject.ResourceVersion
	}
	if kubeEvent.InvolvedObject.FieldPath != "" {
		event.Labels["field_path"] = kubeEvent.InvolvedObject.FieldPath
	}
	
	// Add count if available
	if kubeEvent.Count > 0 {
		event.Labels["count"] = fmt.Sprintf("%d", kubeEvent.Count)
	}
	
	return event
}

// generateEventID generates a unique ID for an event
func generateEventID(kubeEvent *corev1.Event) string {
	// Use a combination of namespace, name, reason, and timestamp to create a unique ID
	baseID := fmt.Sprintf("%s-%s-%s-%s-%d",
		kubeEvent.Namespace,
		kubeEvent.InvolvedObject.Name,
		kubeEvent.Reason,
		kubeEvent.InvolvedObject.UID,
		kubeEvent.FirstTimestamp.Unix(),
	)
	
	// If the base ID would be too long, use UUID
	if len(baseID) > 100 {
		return uuid.New().String()
	}
	
	return baseID
}

// determineSeverity determines the severity level based on the event type and reason
func determineSeverity(kubeEvent *corev1.Event) models.EventSeverity {
	// Determine severity based on event type first
	if kubeEvent.Type == corev1.EventTypeWarning {
		// Check if it's a critical warning based on reason
		switch kubeEvent.Reason {
		case "Failed", "FailedScheduling", "FailedMount", "FailedAttachVolume",
			"FailedDetachVolume", "FailedCreate", "FailedDelete", "Unhealthy",
			"BackOff", "Evicted", "Preempting", "OutOfDisk", "FreeDiskSpaceFailed":
			return models.EventSeverityError
		case "Killing", "NetworkNotReady", "NodeNotReady", "Rebooted",
			"NodeAllocatableEnforced", "SystemOOM", "ContainerGCFailed":
			return models.EventSeverityCritical
		default:
			return models.EventSeverityWarning
		}
	}
	
	// Normal events are generally informational
	if kubeEvent.Type == corev1.EventTypeNormal {
		// But some normal events might be more important
		switch kubeEvent.Reason {
		case "Started", "Created", "Scheduled", "Pulled", "SuccessfulCreate",
			"SuccessfulDelete", "SuccessfulMount", "SuccessfulAttach":
			return models.EventSeverityInfo
		case "NodeReady", "LeaderElection", "RegisteredNode":
			return models.EventSeverityInfo
		default:
			return models.EventSeverityInfo
		}
	}
	
	// Default to info for unknown types
	return models.EventSeverityInfo
}

// GetEventMetrics returns metrics about collected events
func (ec *EventCollector) GetEventMetrics() map[string]int {
	// This could be extended to maintain internal counters
	return map[string]int{
		"total_collected": 0, // Would need to maintain this counter
	}
}

// Reset resets the collector's state
func (ec *EventCollector) Reset() {
	ec.lastTimestamp = time.Now().Add(-5 * time.Minute)
	ec.logger.Info("Event collector state reset")
}

// SetLastTimestamp sets the last timestamp for event collection
func (ec *EventCollector) SetLastTimestamp(timestamp time.Time) {
	ec.lastTimestamp = timestamp
	ec.logger.WithField("timestamp", timestamp).Info("Updated last timestamp for event collection")
}