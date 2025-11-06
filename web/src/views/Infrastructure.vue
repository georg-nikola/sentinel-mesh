<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">
      Infrastructure Overview
    </h1>

    <!-- Nodes Table -->
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden mb-8">
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">
          Cluster Nodes
        </h3>
      </div>
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-900">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Node Name
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              CPU
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Memory
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Pods
            </th>
          </tr>
        </thead>
        <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
          <tr
            v-for="node in nodes"
            :key="node.name"
          >
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
              {{ node.name }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="{
                  'bg-green-100 text-green-800': node.status === 'Ready',
                  'bg-red-100 text-red-800': node.status === 'NotReady',
                }"
              >
                {{ node.status }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
              {{ node.cpu }}%
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
              {{ node.memory }}%
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
              {{ node.pods }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Services Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="service in services"
        :key="service.name"
        class="bg-white dark:bg-gray-800 shadow rounded-lg p-6"
      >
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">
            {{ service.name }}
          </h3>
          <span
            class="px-2 py-1 text-xs font-semibold rounded-full"
            :class="{
              'bg-green-100 text-green-800': service.status === 'Running',
              'bg-red-100 text-red-800': service.status === 'Down',
            }"
          >
            {{ service.status }}
          </span>
        </div>
        <dl class="space-y-2">
          <div class="flex justify-between text-sm">
            <dt class="text-gray-500 dark:text-gray-400">
              Replicas
            </dt>
            <dd class="text-gray-900 dark:text-white font-medium">
              {{ service.replicas }}
            </dd>
          </div>
          <div class="flex justify-between text-sm">
            <dt class="text-gray-500 dark:text-gray-400">
              Port
            </dt>
            <dd class="text-gray-900 dark:text-white font-medium">
              {{ service.port }}
            </dd>
          </div>
          <div class="flex justify-between text-sm">
            <dt class="text-gray-500 dark:text-gray-400">
              Version
            </dt>
            <dd class="text-gray-900 dark:text-white font-medium">
              {{ service.version }}
            </dd>
          </div>
        </dl>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { API_CONFIG } from '@/config/api'

const nodes = ref([
  { name: 'orbstack', status: 'Ready', cpu: 35, memory: 48, pods: 10 },
])

const services = ref([
  { name: 'API', status: 'Running', replicas: '2/2', port: 8080, version: 'v1.0.0', deployment: 'api' },
  { name: 'Collector', status: 'Running', replicas: '2/2', port: 8080, version: 'v1.0.0', deployment: 'collector' },
  { name: 'ML Service', status: 'Running', replicas: '1/1', port: 8000, version: 'v1.0.0', deployment: 'ml-service' },
  { name: 'Prometheus', status: 'Running', replicas: '1/1', port: 9090, version: 'v2.47.0', deployment: 'prometheus' },
  { name: 'Grafana', status: 'Running', replicas: '1/1', port: 3000, version: 'v10.1.0', deployment: 'grafana' },
  { name: 'InfluxDB', status: 'Running', replicas: '1/1', port: 8086, version: 'v2.7', deployment: 'influxdb' },
  { name: 'Redis', status: 'Running', replicas: '1/1', port: 6379, version: 'v7', deployment: 'redis' },
])

const fetchPrometheusMetrics = async () => {
  try {
    const response = await axios.get(`${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=up`, {
      timeout: API_CONFIG.TIMEOUT,
    })
    if (response.data.status === 'success') {
      console.log('Prometheus metrics:', response.data.data.result)
      // Update service status based on prometheus metrics
      response.data.data.result.forEach((metric: any) => {
        const serviceName = metric.metric.job
        const isUp = metric.value[1] === '1'
        const service = services.value.find(s =>
          s.name.toLowerCase().includes(serviceName.toLowerCase())
        )
        if (service) {
          service.status = isUp ? 'Running' : 'Down'
        }
      })
    }
  } catch (error) {
    if (axios.isAxiosError(error) && error.code === 'ECONNABORTED') {
      console.warn('Prometheus request timed out')
    } else {
      console.error('Failed to fetch Prometheus metrics:', error)
    }
  }
}

onMounted(() => {
  fetchPrometheusMetrics()
  // Refresh every 15 seconds
  setInterval(fetchPrometheusMetrics, 15000)
})
</script>
