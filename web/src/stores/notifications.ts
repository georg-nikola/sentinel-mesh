import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Notification {
  id: string
  title: string
  message: string
  severity: 'info' | 'warning' | 'error' | 'critical'
  timestamp: Date
  read: boolean
}

export const useNotificationStore = defineStore('notifications', () => {
  const notifications = ref<Notification[]>([
    {
      id: '1',
      title: 'System Alert',
      message: 'High CPU usage detected on node-1',
      severity: 'warning',
      timestamp: new Date(Date.now() - 5 * 60 * 1000),
      read: false,
    },
    {
      id: '2',
      title: 'Infrastructure Update',
      message: 'Kafka cluster scaled to 3 replicas',
      severity: 'info',
      timestamp: new Date(Date.now() - 15 * 60 * 1000),
      read: false,
    },
  ])

  const unreadCount = computed(() =>
    notifications.value.filter(n => !n.read).length
  )

  const dismissNotification = (id: string) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  const markAsRead = (id: string) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.read = true
    }
  }

  return {
    notifications,
    unreadCount,
    dismissNotification,
    markAsRead,
  }
})
