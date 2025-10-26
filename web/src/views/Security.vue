<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">Security & Alerts</h1>

    <!-- Security Score -->
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6 mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">Security Score</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">Overall security posture of your infrastructure</p>
        </div>
        <div class="text-center">
          <div class="text-5xl font-bold text-green-600">87</div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Good</p>
        </div>
      </div>
    </div>

    <!-- Security Alerts -->
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">Security Alerts</h3>
      </div>
      <ul class="divide-y divide-gray-200 dark:divide-gray-700">
        <li v-for="alert in alerts" :key="alert.id" class="px-6 py-4">
          <div class="flex items-start space-x-3">
            <div
              class="flex-shrink-0 w-2 h-2 mt-2 rounded-full"
              :class="{
                'bg-red-400': alert.severity === 'critical',
                'bg-yellow-400': alert.severity === 'warning',
                'bg-blue-400': alert.severity === 'info',
              }"
            ></div>
            <div class="flex-1">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ alert.title }}</p>
                <div class="flex items-center space-x-2">
                  <span
                    class="px-2 py-1 text-xs font-semibold rounded-full"
                    :class="{
                      'bg-red-100 dark:bg-red-900/50 text-red-800 dark:text-red-300': alert.severity === 'critical',
                      'bg-yellow-100 dark:bg-yellow-900/50 text-yellow-800 dark:text-yellow-300': alert.severity === 'warning',
                      'bg-blue-100 dark:bg-blue-900/50 text-blue-800 dark:text-blue-300': alert.severity === 'info',
                    }"
                  >
                    {{ alert.severity }}
                  </span>
                  <button
                    @click="dismissAlert(alert.id)"
                    class="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300"
                    title="Dismiss alert"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                  </button>
                </div>
              </div>
              <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ alert.description }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 mt-2">{{ formatTime(alert.timestamp) }}</p>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { formatDistanceToNow } from 'date-fns'
import axios from 'axios'
import { API_CONFIG } from '@/config/api'

const alerts = ref([
  {
    id: 1,
    title: 'Collector service failing',
    description: 'Collector pods in CrashLoopBackOff state - Kubernetes config issue',
    severity: 'critical',
    timestamp: new Date(Date.now() - 15 * 60 * 1000),
  },
  {
    id: 2,
    title: 'CPU spike detected',
    description: 'API pod showing elevated CPU usage (95.2%)',
    severity: 'warning',
    timestamp: new Date(Date.now() - 3 * 60 * 1000),
  },
  {
    id: 3,
    title: 'ML Service anomaly detection active',
    description: 'Machine learning service successfully detecting anomalies',
    severity: 'info',
    timestamp: new Date(Date.now() - 5 * 60 * 1000),
  },
])

const dismissAlert = (id: number) => {
  alerts.value = alerts.value.filter(alert => alert.id !== id)
}

const removeStaleAlerts = () => {
  const ONE_HOUR = 60 * 60 * 1000
  const now = Date.now()

  // Remove alerts older than 1 hour
  alerts.value = alerts.value.filter(alert => {
    const age = now - alert.timestamp.getTime()
    return age < ONE_HOUR
  })
}

const fetchMLAnomalies = async () => {
  try {
    // Remove stale alerts first
    removeStaleAlerts()

    const response = await axios.get(`${API_CONFIG.ML_SERVICE_URL}/api/v1/anomalies`, {
      timeout: API_CONFIG.TIMEOUT,
    })
    const anomalies = response.data.anomalies

    // Get current ML anomaly IDs from the response
    const currentMLAnomalyIds = new Set(anomalies.map((_: any, index: number) => 100 + index))

    // Remove old ML anomalies that are no longer present
    alerts.value = alerts.value.filter(alert => {
      // Keep non-ML alerts (id < 100)
      if (alert.id < 100) return true
      // Keep ML alerts that are still in the current response
      return currentMLAnomalyIds.has(alert.id)
    })

    // Add or update ML-detected anomalies
    anomalies.forEach((anomaly: any, index: number) => {
      const existingAlert = alerts.value.find(a => a.id === 100 + index)
      const newAlert = {
        id: 100 + index,
        title: `ML Anomaly: ${anomaly.type.replace('_', ' ')}`,
        description: `Resource: ${anomaly.resource}, Value: ${anomaly.value}%`,
        severity: anomaly.severity,
        timestamp: new Date(anomaly.timestamp),
      }

      if (!existingAlert) {
        alerts.value.unshift(newAlert)
      } else {
        // Update existing alert
        Object.assign(existingAlert, newAlert)
      }
    })
  } catch (error) {
    if (axios.isAxiosError(error) && error.code === 'ECONNABORTED') {
      console.warn('ML Service request timed out')
    } else {
      console.error('Failed to fetch ML anomalies:', error)
    }
  }
}

onMounted(() => {
  fetchMLAnomalies()
  // Refresh every 30 seconds
  setInterval(fetchMLAnomalies, 30000)
})

const formatTime = (timestamp: Date) => {
  return formatDistanceToNow(timestamp, { addSuffix: true })
}
</script>
