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

// Fetch node and service information from API
const fetchInfrastructureData = async () => {
  try {
    const response = await axios.get(
      `${API_CONFIG.BASE_URL}/api/v1/infrastructure`,
      { timeout: API_CONFIG.TIMEOUT }
    )

    if (response.data) {
      // Set nodes from API response
      nodes.value = response.data.nodes || []

      // Set services from API response
      services.value = (response.data.services || []).map((service: any) => ({
        name: service.name,
        status: service.status,
        replicas: service.replicas,
        port: service.port,
        version: service.version,
        deployment: service.deployment
      }))
    }
  } catch (error) {
    console.error('Failed to fetch infrastructure data:', error)
    // Fallback to empty state
    nodes.value = []
    services.value = []
  }
}

const fetchInfrastructure = async () => {
  loading.value = true
  try {
    await fetchInfrastructureData()
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
