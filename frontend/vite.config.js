import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    // Durante el desarrollo /api se redirige al backend de Go.
    // De esta manera el frontend siempre utiliza fetch("/api/...").
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  test: {
    environment: 'node',
    include: ['tests/**/*.test.js'],
    coverage: {
      provider: 'v8',
      // Se mide la lógica que se testea sin DOM: utilidades y cliente de la API.
      // Los componentes React y el arranque (main.jsx, App.jsx) quedan afuera: son UI y
      // se verifican con otro tipo de test, no con unit tests.
      include: ['src/utils/**', 'src/api/**'],
      reporter: ['text', 'text-summary', 'lcov', 'html', 'json-summary'],
      reportsDirectory: './coverage',
      // Umbral que ROMPE el build: si baja de esto, vitest sale con código 1.
      thresholds: {
        lines: 80,
        branches: 85
      }
    }
  }
})
