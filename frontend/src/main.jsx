import React from 'react'
import ReactDOM from 'react-dom/client'
import * as Sentry from '@sentry/react'
import App from './App'
import { despertarServidor } from './api/api'
import './styles.css'

// Sentry: cuenta los errores del navegador y mide cuánto tardan las páginas y los pedidos.
// El DSN va escrito acá A PROPÓSITO y no es un descuido: el DSN de un front no es un secreto.
// Solo sirve para ENVIAR datos a Sentry, nunca para leerlos, y cualquier aplicación web lo expone
// porque viaja en el JavaScript que baja el navegador. Es distinto de una clave de API, que sí
// permite leer o modificar. (El DSN del backend, en cambio, va por variable de entorno.)
Sentry.init({
  dsn: 'https://a9e20ef7b8a9b8cc12db8124b7b777e4@o4512312832819200.ingest.us.sentry.io/4512312872534016',
  environment: window.location.hostname, // de qué sitio vino: QA, producción o localhost
  integrations: [Sentry.browserTracingIntegration()],
  tracesSampleRate: 1.0 // medir todos los pedidos (alcanza para una aplicación de cursada)
})

// Despierta el backend (plan gratis) mientras la página termina de cargar.
despertarServidor()

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
