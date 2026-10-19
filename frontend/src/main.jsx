import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import { despertarServidor } from './api/api'
import './styles.css'

// Despierta el backend (plan gratis) mientras la página termina de cargar.
despertarServidor()

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
