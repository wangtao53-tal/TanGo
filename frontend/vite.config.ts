import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')

  // 从环境变量读取配置，如果没有则使用默认值
  const frontendPort = parseInt(env.FRONTEND_PORT || '3000', 10)
  const backendPort = env.BACKEND_PORT || '8877'
  const backendHost = env.BACKEND_HOST || 'localhost'
  const backendUrl = `http://${backendHost}:${backendPort}`

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      port: frontendPort,
      proxy: {
        '/api': {
          target: backendUrl,
          changeOrigin: true,
        },
      },
    },
  }
})
