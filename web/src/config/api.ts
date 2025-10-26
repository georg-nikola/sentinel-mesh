// API configuration based on environment
const isDevelopment = import.meta.env.MODE === 'development'

export const API_CONFIG = {
  // In development, use localhost ports
  // In production (K8s), use relative URLs that nginx will proxy
  PROMETHEUS_URL: isDevelopment
    ? 'http://localhost:30090'
    : '/api/prometheus',

  ML_SERVICE_URL: isDevelopment
    ? 'http://localhost:8002'
    : '/api/ml',

  // API timeout in milliseconds
  TIMEOUT: 5000,
}
