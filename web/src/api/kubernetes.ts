import axios from 'axios'

const k8sApi = axios.create({
  baseURL: 'http://localhost:30080/api',
  timeout: 5000,
})

export interface PodInfo {
  name: string
  namespace: string
  status: string
  cpu: string
  memory: string
  restarts: number
  age: string
}

export interface NodeInfo {
  name: string
  status: string
  cpu: number
  memory: number
  pods: number
}

export interface ServiceInfo {
  name: string
  status: string
  replicas: string
  port: number
  version: string
}

export const fetchPods = async (): Promise<PodInfo[]> => {
  try {
    const response = await k8sApi.get('/pods')
    return response.data
  } catch (error) {
    console.error('Failed to fetch pods:', error)
    return []
  }
}

export const fetchNodes = async (): Promise<NodeInfo[]> => {
  try {
    const response = await k8sApi.get('/nodes')
    return response.data
  } catch (error) {
    console.error('Failed to fetch nodes:', error)
    return []
  }
}

export const fetchServices = async (): Promise<ServiceInfo[]> => {
  try {
    const response = await k8sApi.get('/services')
    return response.data
  } catch (error) {
    console.error('Failed to fetch services:', error)
    return []
  }
}

export const fetchMetrics = async (metric: string, duration: string = '1h') => {
  try {
    const response = await k8sApi.get(`/metrics/${metric}`, {
      params: { duration },
    })
    return response.data
  } catch (error) {
    console.error(`Failed to fetch ${metric} metrics:`, error)
    return { data: [] }
  }
}
