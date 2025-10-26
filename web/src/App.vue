<template>
  <div id="app" class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Navigation -->
    <nav class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <!-- Logo -->
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-gray-900 dark:text-white">
                <span class="text-blue-600 dark:text-blue-400">Sentinel</span> Mesh
              </h1>
            </div>

            <!-- Navigation Links -->
            <div class="hidden md:ml-6 md:flex md:space-x-8">
              <router-link
                v-for="item in navigation"
                :key="item.name"
                :to="item.href"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium border-b-2"
                :class="
                  $route.path === item.href
                    ? 'border-blue-500 text-gray-900 dark:text-white'
                    : 'border-transparent text-gray-500 dark:text-gray-400 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-300'
                "
              >
                <component :is="item.icon" class="w-4 h-4 mr-2" />
                {{ item.name }}
              </router-link>
            </div>
          </div>

          <!-- Right side -->
          <div class="flex items-center space-x-4">
            <!-- Connection Status - hidden on small screens -->
            <div class="hidden sm:flex items-center space-x-2">
              <div
                class="w-2 h-2 rounded-full"
                :class="connectionStatus === 'connected' ? 'bg-green-400' : 'bg-red-400'"
              ></div>
              <span class="text-sm text-gray-500 dark:text-gray-400">
                {{ connectionStatus === 'connected' ? 'Connected' : 'Disconnected' }}
              </span>
            </div>

            <!-- Theme Toggle -->
            <button
              @click="toggleTheme"
              class="p-2 text-gray-400 dark:text-gray-300 hover:text-gray-500 dark:hover:text-gray-200"
              title="Toggle theme"
            >
              <SunIcon v-if="theme === 'dark'" class="w-5 h-5" />
              <MoonIcon v-else class="w-5 h-5" />
            </button>

            <!-- Notifications -->
            <button
              @click="showNotifications = !showNotifications"
              class="relative p-2 text-gray-400 dark:text-gray-300 hover:text-gray-500 dark:hover:text-gray-200"
            >
              <BellIcon class="w-5 h-5" />
              <span
                v-if="unreadNotifications > 0"
                class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center"
              >
                {{ unreadNotifications }}
              </span>
            </button>

            <!-- User Menu -->
            <div class="relative">
              <button
                @click="showUserMenu = !showUserMenu"
                class="flex items-center text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
                  <span class="text-white text-sm font-medium">U</span>
                </div>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Mobile Navigation -->
      <div class="md:hidden" v-if="showMobileMenu">
        <div class="pt-2 pb-3 space-y-1">
          <router-link
            v-for="item in navigation"
            :key="item.name"
            :to="item.href"
            class="block pl-3 pr-4 py-2 text-base font-medium"
            :class="
              $route.path === item.href
                ? 'bg-blue-50 border-blue-500 text-blue-700'
                : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700'
            "
          >
            {{ item.name }}
          </router-link>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="flex-1">
      <router-view />
    </main>

    <!-- Global Notifications -->
    <div
      v-if="showNotifications"
      class="fixed inset-0 z-50 overflow-hidden"
      @click="showNotifications = false"
    >
      <div class="absolute inset-0 bg-black bg-opacity-25 dark:bg-opacity-50"></div>
      <div class="absolute right-0 top-16 w-96 bg-white dark:bg-gray-800 shadow-lg rounded-lg m-4">
        <div class="p-4">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Notifications</h3>
          <div class="space-y-2">
            <div
              v-for="notification in notifications"
              :key="notification.id"
              class="p-3 rounded-lg border"
              :class="getNotificationClass(notification.severity)"
            >
              <div class="flex justify-between items-start">
                <div>
                  <p class="text-sm font-medium dark:text-white">{{ notification.title }}</p>
                  <p class="text-sm text-gray-600 dark:text-gray-300">{{ notification.message }}</p>
                  <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ formatTime(notification.timestamp) }}</p>
                </div>
                <button
                  @click.stop="dismissNotification(notification.id)"
                  class="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div
      v-if="loading"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 flex items-center space-x-3">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
        <span class="text-gray-700">Loading...</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  BellIcon,
  XMarkIcon,
  HomeIcon,
  ChartBarIcon,
  CpuChipIcon,
  ShieldCheckIcon,
  Cog6ToothIcon,
  MoonIcon,
  SunIcon
} from '@heroicons/vue/24/outline'
import { useWebSocketStore } from '@/stores/websocket'
import { useNotificationStore } from '@/stores/notifications'
import { useTheme } from '@/composables/useTheme'
import { formatDistanceToNow } from 'date-fns'

// Stores & Composables
const websocketStore = useWebSocketStore()
const notificationStore = useNotificationStore()
const { theme, toggleTheme, initTheme } = useTheme()

// Router
const router = useRouter()

// Reactive state
const showMobileMenu = ref(false)
const showNotifications = ref(false)
const showUserMenu = ref(false)
const loading = ref(false)

// Navigation items
const navigation = [
  { name: 'Dashboard', href: '/', icon: HomeIcon },
  { name: 'Metrics', href: '/metrics', icon: ChartBarIcon },
  { name: 'Infrastructure', href: '/infrastructure', icon: CpuChipIcon },
  { name: 'Security', href: '/security', icon: ShieldCheckIcon },
  { name: 'Settings', href: '/settings', icon: Cog6ToothIcon },
]

// Computed properties
const connectionStatus = computed(() => websocketStore.connectionStatus)
const notifications = computed(() => notificationStore.notifications)
const unreadNotifications = computed(() => notificationStore.unreadCount)

// Methods
const formatTime = (timestamp: Date) => {
  return formatDistanceToNow(timestamp, { addSuffix: true })
}

const getNotificationClass = (severity: string) => {
  const classes = {
    info: 'border-blue-200 dark:border-blue-700 bg-blue-50 dark:bg-blue-900/30',
    warning: 'border-yellow-200 dark:border-yellow-700 bg-yellow-50 dark:bg-yellow-900/30',
    error: 'border-red-200 dark:border-red-700 bg-red-50 dark:bg-red-900/30',
    critical: 'border-red-300 dark:border-red-600 bg-red-100 dark:bg-red-900/40',
  }
  return classes[severity as keyof typeof classes] || classes.info
}

const dismissNotification = (id: string) => {
  notificationStore.dismissNotification(id)
}

// Lifecycle
onMounted(() => {
  // Initialize theme
  initTheme()

  // Initialize WebSocket connection
  websocketStore.connect()

  // Close dropdowns when clicking outside
  document.addEventListener('click', (event) => {
    const target = event.target as HTMLElement
    if (!target.closest('.relative')) {
      showUserMenu.value = false
    }
  })
})

onUnmounted(() => {
  websocketStore.disconnect()
})
</script>

<style scoped>
/* Custom styles for the app */
.router-link-active {
  @apply border-blue-500 text-gray-900;
}

/* Loading animation */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>