import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useWebSocketStore = defineStore('websocket', () => {
  const connectionStatus = ref<'connected' | 'disconnected'>('disconnected')
  const socket = ref<WebSocket | null>(null)

  const connect = () => {
    // Placeholder WebSocket connection
    // In production, this would connect to the actual WebSocket endpoint
    connectionStatus.value = 'connected'
    console.log('WebSocket connected (mock)')
  }

  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
    }
    connectionStatus.value = 'disconnected'
    console.log('WebSocket disconnected')
  }

  return {
    connectionStatus,
    connect,
    disconnect,
  }
})
