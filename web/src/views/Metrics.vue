<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">
      Metrics & Analytics
    </h1>

    <!-- Time Range Selector -->
    <div class="mb-6 flex justify-end">
      <select class="block w-48 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm">
        <option>Last 1 hour</option>
        <option>Last 6 hours</option>
        <option selected>
          Last 24 hours
        </option>
        <option>Last 7 days</option>
        <option>Last 30 days</option>
      </select>
    </div>

    <!-- Metrics Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Request Rate -->
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Request Rate
        </h3>
        <div class="h-64">
          <Line
            :data="requestRateData"
            :options="chartOptions"
          />
        </div>
      </div>

      <!-- CPU Usage -->
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
          CPU Usage
        </h3>
        <div class="h-64">
          <Line
            :data="cpuUsageData"
            :options="chartOptions"
          />
        </div>
      </div>

      <!-- Memory Usage -->
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Memory Usage
        </h3>
        <div class="h-64">
          <Line
            :data="memoryUsageData"
            :options="chartOptions"
          />
        </div>
      </div>

      <!-- Response Time -->
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
          Response Time
        </h3>
        <div class="h-64">
          <Line
            :data="responseTimeData"
            :options="chartOptions"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'
import axios from 'axios'
import { Line } from 'vue-chartjs'
import { API_CONFIG } from '@/config/api'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const requestRateData = shallowRef({
  labels: ['...'] as string[],
  datasets: [
    {
      label: 'Requests/sec',
      data: [0] as number[],
      borderColor: 'rgb(99, 102, 241)',
      backgroundColor: 'rgba(99, 102, 241, 0.1)',
      fill: true,
      tension: 0.4,
    },
  ],
})

const cpuUsageData = shallowRef({
  labels: ['...'] as string[],
  datasets: [
    {
      label: 'CPU %',
      data: [0] as number[],
      borderColor: 'rgb(239, 68, 68)',
      backgroundColor: 'rgba(239, 68, 68, 0.1)',
      fill: true,
      tension: 0.4,
    },
  ],
})

const memoryUsageData = shallowRef({
  labels: ['...'] as string[],
  datasets: [
    {
      label: 'Memory MB',
      data: [0] as number[],
      borderColor: 'rgb(16, 185, 129)',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true,
      tension: 0.4,
    },
  ],
})

const responseTimeData = shallowRef({
  labels: ['...'] as string[],
  datasets: [
    {
      label: 'Response Time (ms)',
      data: [0] as number[],
      borderColor: 'rgb(245, 158, 11)',
      backgroundColor: 'rgba(245, 158, 11, 0.1)',
      fill: true,
      tension: 0.4,
    },
  ],
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false,
    },
  },
  scales: {
    y: {
      beginAtZero: true,
    },
  },
}

const fetchMetrics = async () => {
  try {
    // Fetch all metrics from our API
    const response = await axios.get(`${API_CONFIG.BASE_URL}/api/v1/metrics`, {
      timeout: API_CONFIG.TIMEOUT,
    })

    if (response.data && response.data.metrics) {
      const metrics = response.data.metrics
      const timestamp = new Date().toLocaleTimeString()

      // Update Request Rate chart (using request_rate metric)
      if (metrics.request_rate && metrics.request_rate.length > 0) {
        const latestValue = metrics.request_rate[metrics.request_rate.length - 1].value

        if (requestRateData.value.labels[0] === '...') {
          requestRateData.value = {
            labels: [timestamp],
            datasets: [{
              label: 'Requests/sec',
              data: [latestValue],
              borderColor: 'rgb(99, 102, 241)',
              backgroundColor: 'rgba(99, 102, 241, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        } else {
          const newLabels = [...requestRateData.value.labels, timestamp]
          const newData = [...requestRateData.value.datasets[0].data, latestValue]

          if (newLabels.length > 20) {
            newLabels.shift()
            newData.shift()
          }

          requestRateData.value = {
            labels: newLabels,
            datasets: [{
              label: 'Requests/sec',
              data: newData,
              borderColor: 'rgb(99, 102, 241)',
              backgroundColor: 'rgba(99, 102, 241, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        }
      }

      // Update CPU Usage chart
      if (metrics.node_cpu_usage && metrics.node_cpu_usage.length > 0) {
        const latestValue = metrics.node_cpu_usage[metrics.node_cpu_usage.length - 1].value

        if (cpuUsageData.value.labels[0] === '...') {
          cpuUsageData.value = {
            labels: [timestamp],
            datasets: [{
              label: 'CPU %',
              data: [latestValue],
              borderColor: 'rgb(239, 68, 68)',
              backgroundColor: 'rgba(239, 68, 68, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        } else {
          const newLabels = [...cpuUsageData.value.labels, timestamp]
          const newData = [...cpuUsageData.value.datasets[0].data, latestValue]

          if (newLabels.length > 20) {
            newLabels.shift()
            newData.shift()
          }

          cpuUsageData.value = {
            labels: newLabels,
            datasets: [{
              label: 'CPU %',
              data: newData,
              borderColor: 'rgb(239, 68, 68)',
              backgroundColor: 'rgba(239, 68, 68, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        }
      }

      // Update Memory Usage chart
      if (metrics.node_memory_usage && metrics.node_memory_usage.length > 0) {
        const latestValue = metrics.node_memory_usage[metrics.node_memory_usage.length - 1].value

        if (memoryUsageData.value.labels[0] === '...') {
          memoryUsageData.value = {
            labels: [timestamp],
            datasets: [{
              label: 'Memory MB',
              data: [latestValue],
              borderColor: 'rgb(16, 185, 129)',
              backgroundColor: 'rgba(16, 185, 129, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        } else {
          const newLabels = [...memoryUsageData.value.labels, timestamp]
          const newData = [...memoryUsageData.value.datasets[0].data, latestValue]

          if (newLabels.length > 20) {
            newLabels.shift()
            newData.shift()
          }

          memoryUsageData.value = {
            labels: newLabels,
            datasets: [{
              label: 'Memory MB',
              data: newData,
              borderColor: 'rgb(16, 185, 129)',
              backgroundColor: 'rgba(16, 185, 129, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        }
      }

      // Update Response Time chart (using request_latency metric)
      if (metrics.request_latency && metrics.request_latency.length > 0) {
        const latestValue = metrics.request_latency[metrics.request_latency.length - 1].value

        if (responseTimeData.value.labels[0] === '...') {
          responseTimeData.value = {
            labels: [timestamp],
            datasets: [{
              label: 'Response Time (ms)',
              data: [latestValue],
              borderColor: 'rgb(245, 158, 11)',
              backgroundColor: 'rgba(245, 158, 11, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        } else {
          const newLabels = [...responseTimeData.value.labels, timestamp]
          const newData = [...responseTimeData.value.datasets[0].data, latestValue]

          if (newLabels.length > 20) {
            newLabels.shift()
            newData.shift()
          }

          responseTimeData.value = {
            labels: newLabels,
            datasets: [{
              label: 'Response Time (ms)',
              data: newData,
              borderColor: 'rgb(245, 158, 11)',
              backgroundColor: 'rgba(245, 158, 11, 0.1)',
              fill: true,
              tension: 0.4,
            }],
          }
        }
      }
    }
  } catch (error) {
    if (axios.isAxiosError(error) && error.code === 'ECONNABORTED') {
      console.warn('Metrics request timed out')
    } else {
      console.error('Failed to fetch metrics:', error)
    }
  }
}

onMounted(() => {
  fetchMetrics()
  // Refresh every 15 seconds
  setInterval(fetchMetrics, 15000)
})
</script>
