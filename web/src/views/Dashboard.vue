<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">Dashboard Overview</h1>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8">
      <div v-for="stat in stats" :key="stat.name" class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
        <div class="p-5">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <component :is="stat.icon" class="h-6 w-6 text-gray-400" />
            </div>
            <div class="ml-5 w-0 flex-1">
              <dl>
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">{{ stat.name }}</dt>
                <dd class="flex items-baseline">
                  <div class="text-2xl font-semibold text-gray-900 dark:text-white">{{ stat.value }}</div>
                  <div
                    class="ml-2 flex items-baseline text-sm font-semibold"
                    :class="stat.change >= 0 ? 'text-green-600' : 'text-red-600'"
                  >
                    {{ stat.change >= 0 ? '+' : '' }}{{ stat.change }}%
                  </div>
                </dd>
              </dl>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
      <!-- CPU Usage Chart -->
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">CPU Usage</h3>
        <div class="h-64">
          <Line :data="cpuChartData" :options="chartOptions" />
        </div>
      </div>

      <!-- Memory Usage Chart -->
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Memory Usage</h3>
        <div class="h-64">
          <Line :data="memoryChartData" :options="chartOptions" />
        </div>
      </div>
    </div>

    <!-- Recent Events -->
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg">
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">Recent Events</h3>
      </div>
      <ul class="divide-y divide-gray-200 dark:divide-gray-700">
        <li v-for="event in recentEvents" :key="event.id" class="px-6 py-4 hover:bg-gray-50 dark:hover:bg-gray-700">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <div
                class="flex-shrink-0 w-2 h-2 rounded-full"
                :class="{
                  'bg-green-400': event.type === 'success',
                  'bg-yellow-400': event.type === 'warning',
                  'bg-red-400': event.type === 'error',
                  'bg-blue-400': event.type === 'info',
                }"
              ></div>
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ event.message }}</p>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ event.source }}</p>
              </div>
            </div>
            <span class="text-sm text-gray-500 dark:text-gray-400">{{ formatTime(event.timestamp) }}</span>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, shallowRef } from 'vue'
import { CpuChipIcon, ServerIcon, CircleStackIcon, CloudIcon } from '@heroicons/vue/24/outline'
import { formatDistanceToNow } from 'date-fns'
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

const stats = ref([
  { name: 'Active Nodes', value: '1', change: 0, icon: ServerIcon },
  { name: 'Services Up', value: '0', change: 0, icon: CpuChipIcon },
  { name: 'Services Down', value: '0', change: 0, icon: CircleStackIcon },
  { name: 'Active Pods', value: '10', change: 0, icon: CloudIcon },
])

const recentEvents = ref([
  {
    id: 1,
    message: 'API service running with 2 replicas',
    source: 'Kubernetes',
    type: 'success',
    timestamp: new Date(Date.now() - 2 * 60 * 1000),
  },
  {
    id: 2,
    message: 'Collector in CrashLoopBackOff - needs K8s config',
    source: 'Monitoring',
    type: 'error',
    timestamp: new Date(Date.now() - 5 * 60 * 1000),
  },
  {
    id: 3,
    message: 'ML Service anomaly detection active',
    source: 'ML Service',
    type: 'info',
    timestamp: new Date(Date.now() - 10 * 60 * 1000),
  },
  {
    id: 4,
    message: 'Prometheus scraping 4 targets',
    source: 'Prometheus',
    type: 'success',
    timestamp: new Date(Date.now() - 15 * 60 * 1000),
  },
])

// Chart data - Initialize with empty point to ensure chart renders
const cpuChartData = shallowRef({
  labels: ['...'] as string[],
  datasets: [
    {
      label: 'CPU Usage (%)',
      data: [0] as number[],
      borderColor: 'rgb(59, 130, 246)',
      backgroundColor: 'rgba(59, 130, 246, 0.1)',
      fill: true,
      tension: 0.4,
    },
  ],
})

const memoryChartData = shallowRef({
  labels: ['...'] as string[],
  datasets: [
    {
      label: 'Memory Usage (MB)',
      data: [0] as number[],
      borderColor: 'rgb(16, 185, 129)',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true,
      tension: 0.4,
    },
  ],
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  animation: {
    duration: 750,
  },
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

const fetchPrometheusData = async () => {
  try {
    // Fetch service status with timeout
    const upQuery = await axios.get(`${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=up`, {
      timeout: API_CONFIG.TIMEOUT,
    })
    if (upQuery.data.status === 'success') {
      const results = upQuery.data.data.result
      const servicesUp = results.filter((r: any) => r.value[1] === '1').length
      const servicesDown = results.filter((r: any) => r.value[1] === '0').length

      stats.value[1].value = servicesUp.toString()
      stats.value[2].value = servicesDown.toString()

      // Add events for service status changes
      results.forEach((metric: any) => {
        const serviceName = metric.metric.job
        const isUp = metric.value[1] === '1'
        const existingEvent = recentEvents.value.find(e =>
          e.message.includes(serviceName)
        )

        if (!existingEvent) {
          recentEvents.value.unshift({
            id: Date.now() + Math.random(),
            message: `${serviceName} service is ${isUp ? 'up' : 'down'}`,
            source: 'Prometheus',
            type: isUp ? 'success' : 'error',
            timestamp: new Date(),
          })
          // Keep only last 10 events
          if (recentEvents.value.length > 10) {
            recentEvents.value = recentEvents.value.slice(0, 10)
          }
        }
      })
    }

    // Fetch Go memory stats
    const memQuery = await axios.get(`${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=go_memstats_alloc_bytes`, {
      timeout: API_CONFIG.TIMEOUT,
    })
    if (memQuery.data.status === 'success') {
      const totalMem = memQuery.data.data.result.reduce((acc: number, r: any) => {
        return acc + parseFloat(r.value[1])
      }, 0)
      const memMB = (totalMem / 1024 / 1024).toFixed(0)
      console.log(`Total memory usage: ${memMB}MB`)

      // Update memory chart
      const timestamp = new Date().toLocaleTimeString()

      // Replace placeholder data on first update
      if (memoryChartData.value.labels[0] === '...') {
        memoryChartData.value = {
          labels: [timestamp],
          datasets: [{
            label: 'Memory Usage (MB)',
            data: [parseFloat(memMB)],
            borderColor: 'rgb(16, 185, 129)',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            fill: true,
            tension: 0.4,
          }],
        }
      } else {
        const newLabels = [...memoryChartData.value.labels, timestamp]
        const newData = [...memoryChartData.value.datasets[0].data, parseFloat(memMB)]

        // Keep only last 20 data points
        if (newLabels.length > 20) {
          newLabels.shift()
          newData.shift()
        }

        memoryChartData.value = {
          labels: newLabels,
          datasets: [{
            label: 'Memory Usage (MB)',
            data: newData,
            borderColor: 'rgb(16, 185, 129)',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            fill: true,
            tension: 0.4,
          }],
        }
      }
    }

    // Fetch CPU stats (using goroutines as a proxy for activity)
    const cpuQuery = await axios.get(`${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=go_goroutines`, {
      timeout: API_CONFIG.TIMEOUT,
    })
    if (cpuQuery.data.status === 'success') {
      const totalGoroutines = cpuQuery.data.data.result.reduce((acc: number, r: any) => {
        return acc + parseFloat(r.value[1])
      }, 0)

      // Update CPU chart (normalized goroutine count as percentage)
      const timestamp = new Date().toLocaleTimeString()
      const cpuProxy = Math.min((totalGoroutines / 10), 100) // Normalize to 0-100

      // Replace placeholder data on first update
      if (cpuChartData.value.labels[0] === '...') {
        cpuChartData.value = {
          labels: [timestamp],
          datasets: [{
            label: 'CPU Usage (%)',
            data: [cpuProxy],
            borderColor: 'rgb(59, 130, 246)',
            backgroundColor: 'rgba(59, 130, 246, 0.1)',
            fill: true,
            tension: 0.4,
          }],
        }
      } else {
        const newLabels = [...cpuChartData.value.labels, timestamp]
        const newData = [...cpuChartData.value.datasets[0].data, cpuProxy]

        // Keep only last 20 data points
        if (newLabels.length > 20) {
          newLabels.shift()
          newData.shift()
        }

        cpuChartData.value = {
          labels: newLabels,
          datasets: [{
            label: 'CPU Usage (%)',
            data: newData,
            borderColor: 'rgb(59, 130, 246)',
            backgroundColor: 'rgba(59, 130, 246, 0.1)',
            fill: true,
            tension: 0.4,
          }],
        }
      }
    }
  } catch (error) {
    if (axios.isAxiosError(error) && error.code === 'ECONNABORTED') {
      console.warn('Prometheus request timed out')
    } else {
      console.error('Failed to fetch Prometheus data:', error)
    }
  }
}

onMounted(() => {
  fetchPrometheusData()
  // Refresh every 15 seconds
  setInterval(fetchPrometheusData, 15000)
})

const formatTime = (timestamp: Date) => {
  return formatDistanceToNow(timestamp, { addSuffix: true })
}
</script>
