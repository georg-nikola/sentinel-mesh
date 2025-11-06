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

interface Node {
  name: string
  status: string
  cpu: number
  memory: number
  pods: number
}

interface Service {
  name: string
  status: string
  replicas: string
  port: number | string
  version: string
  deployment: string
}

const nodes = ref<Node[]>([])
const services = ref<Service[]>([])
const loading = ref(true)

// Fetch node information from Prometheus node_exporter metrics
const fetchNodes = async () => {
  try {
    // Query for node CPU usage
    const cpuResponse = await axios.get(
      `${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`,
      { timeout: API_CONFIG.TIMEOUT }
    )

    // Query for node memory usage
    const memResponse = await axios.get(
      `${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=100 * (1 - ((node_memory_MemAvailable_bytes) / (node_memory_MemTotal_bytes)))`,
      { timeout: API_CONFIG.TIMEOUT }
    )

    // Query for node status (up metric)
    const statusResponse = await axios.get(
      `${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=up{job="node-exporter"}`,
      { timeout: API_CONFIG.TIMEOUT }
    )

    // Query for pod count per node
    const podsResponse = await axios.get(
      `${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=count by(node) (kube_pod_info)`,
      { timeout: API_CONFIG.TIMEOUT }
    )

    const nodeMap = new Map<string, Partial<Node>>()

    // Process CPU data
    if (cpuResponse.data?.data?.result) {
      cpuResponse.data.data.result.forEach((metric: any) => {
        const instance = metric.metric.instance?.split(':')[0] || 'unknown'
        const cpu = Math.round(parseFloat(metric.value[1]))
        nodeMap.set(instance, { ...nodeMap.get(instance), name: instance, cpu })
      })
    }

    // Process memory data
    if (memResponse.data?.data?.result) {
      memResponse.data.data.result.forEach((metric: any) => {
        const instance = metric.metric.instance?.split(':')[0] || 'unknown'
        const memory = Math.round(parseFloat(metric.value[1]))
        const existing = nodeMap.get(instance) || {}
        nodeMap.set(instance, { ...existing, memory })
      })
    }

    // Process status data
    if (statusResponse.data?.data?.result) {
      statusResponse.data.data.result.forEach((metric: any) => {
        const instance = metric.metric.instance?.split(':')[0] || 'unknown'
        const isUp = metric.value[1] === '1'
        const existing = nodeMap.get(instance) || {}
        nodeMap.set(instance, { ...existing, status: isUp ? 'Ready' : 'NotReady' })
      })
    }

    // Process pod count data
    if (podsResponse.data?.data?.result) {
      podsResponse.data.data.result.forEach((metric: any) => {
        const nodeName = metric.metric.node
        const pods = parseInt(metric.value[1])
        // Try to find matching node by name
        for (const [instance, node] of nodeMap.entries()) {
          if (instance.includes(nodeName) || nodeName.includes(instance)) {
            nodeMap.set(instance, { ...node, pods })
            break
          }
        }
      })
    }

    nodes.value = Array.from(nodeMap.values()).map(node => ({
      name: node.name || 'unknown',
      status: node.status || 'Unknown',
      cpu: node.cpu || 0,
      memory: node.memory || 0,
      pods: node.pods || 0
    }))

    // If no nodes found from metrics, try kube_node_info
    if (nodes.value.length === 0) {
      const nodeInfoResponse = await axios.get(
        `${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=kube_node_info`,
        { timeout: API_CONFIG.TIMEOUT }
      )

      if (nodeInfoResponse.data?.data?.result) {
        nodes.value = nodeInfoResponse.data.data.result.map((metric: any) => ({
          name: metric.metric.node || 'unknown',
          status: 'Ready',
          cpu: 0,
          memory: 0,
          pods: 0
        }))
      }
    }
  } catch (error) {
    console.error('Failed to fetch node data:', error)
    // Fallback to basic node info
    nodes.value = [{ name: 'Node data unavailable', status: 'Unknown', cpu: 0, memory: 0, pods: 0 }]
  }
}

// Fetch service information from Prometheus up metrics
const fetchServices = async () => {
  try {
    const response = await axios.get(
      `${API_CONFIG.PROMETHEUS_URL}/api/v1/query?query=up`,
      { timeout: API_CONFIG.TIMEOUT }
    )

    if (response.data?.status === 'success') {
      const serviceMap = new Map<string, Service>()

      response.data.data.result.forEach((metric: any) => {
        const job = metric.metric.job
        if (!job || job === 'node-exporter') return // Skip node-exporter

        const instance = metric.metric.instance || ''
        const isUp = metric.value[1] === '1'
        const port = instance.includes(':') ? instance.split(':')[1] : '8080'

        // Format service name
        let serviceName = job
          .split('-')
          .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' ')

        // Get version from metric if available
        const version = metric.metric.version || 'v1.0.0'

        if (!serviceMap.has(job)) {
          serviceMap.set(job, {
            name: serviceName,
            status: isUp ? 'Running' : 'Down',
            replicas: '1/1',
            port: port,
            version: version,
            deployment: job
          })
        } else {
          // Update status if service is down
          const existing = serviceMap.get(job)!
          if (!isUp) {
            existing.status = 'Down'
          }
        }
      })

      services.value = Array.from(serviceMap.values()).sort((a, b) =>
        a.name.localeCompare(b.name)
      )
    }
  } catch (error) {
    console.error('Failed to fetch service data:', error)
    services.value = []
  }
}

const fetchInfrastructure = async () => {
  loading.value = true
  try {
    await Promise.all([fetchNodes(), fetchServices()])
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchInfrastructure()
  // Refresh every 30 seconds
  setInterval(fetchInfrastructure, 30000)
})
</script>
