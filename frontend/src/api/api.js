// Todas las llamadas HTTP de la aplicación están centralizadas aquí.
// Se utiliza fetch con rutas relativas: en desarrollo Vite redirige /api al
// backend y en Docker lo hace Nginx.

const BASE = '/api'
const CLAVE_TOKEN = 'gestorPeliculas.token'

// El token se guarda en localStorage para que la sesión sobreviva a un F5.
// El acceso va dentro de try/catch porque algunos navegadores lo bloquean.
export function obtenerToken() {
  try {
    return localStorage.getItem(CLAVE_TOKEN)
  } catch {
    return null
  }
}

export function guardarToken(token) {
  try {
    localStorage.setItem(CLAVE_TOKEN, token)
  } catch {
    // Si no se puede guardar, la sesión durará solamente lo que dure la página.
  }
}

export function borrarToken() {
  try {
    localStorage.removeItem(CLAVE_TOKEN)
  } catch {
    // Nada que hacer.
  }
}

// La aplicación registra aquí qué hacer cuando el backend rechaza el token.
let alExpirarSesion = null

export function configurarExpiracionSesion(callback) {
  alExpirarSesion = callback
}

async function peticion(ruta, opciones = {}) {
  const token = obtenerToken()
  const cabeceras = { 'Content-Type': 'application/json', ...(opciones.headers || {}) }
  if (token) {
    cabeceras.Authorization = `Bearer ${token}`
  }

  let respuesta
  try {
    respuesta = await fetch(ruta, { ...opciones, headers: cabeceras })
  } catch {
    throw new Error('No se pudo conectar con el servidor.')
  }

  // Si el token dejó de servir se cierra la sesión.
  // Un 401 al iniciar sesión (todavía sin token) no cuenta como expiración.
  if (respuesta.status === 401 && token) {
    borrarToken()
    if (alExpirarSesion) {
      alExpirarSesion()
    }
  }

  if (respuesta.status === 204) {
    return null
  }

  let datos = null
  try {
    datos = await respuesta.json()
  } catch {
    datos = null
  }

  if (!respuesta.ok) {
    const error = new Error((datos && datos.error) || 'Ocurrió un error inesperado.')
    error.status = respuesta.status
    throw error
  }

  return datos
}

// Autenticación

export async function registrarse(cuenta) {
  const sesion = await peticion(`${BASE}/auth/registro`, {
    method: 'POST',
    body: JSON.stringify(cuenta)
  })
  guardarToken(sesion.token)
  return sesion
}

export async function iniciarSesion(credenciales) {
  const sesion = await peticion(`${BASE}/auth/login`, {
    method: 'POST',
    body: JSON.stringify(credenciales)
  })
  guardarToken(sesion.token)
  return sesion
}

export function obtenerPerfil() {
  return peticion(`${BASE}/auth/perfil`)
}

export function cerrarSesion() {
  borrarToken()
}

// Arma la query string a partir de los filtros con valor.
function construirConsulta(filtros = {}) {
  const parametros = new URLSearchParams()
  Object.entries(filtros).forEach(([clave, valor]) => {
    if (valor !== '' && valor !== null && valor !== undefined) {
      parametros.append(clave, valor)
    }
  })
  const consulta = parametros.toString()
  return consulta ? `?${consulta}` : ''
}

// Películas

export function obtenerPeliculas(filtros) {
  return peticion(`${BASE}/peliculas${construirConsulta(filtros)}`)
}

export function obtenerPelicula(id) {
  return peticion(`${BASE}/peliculas/${id}`)
}

export function crearPelicula(pelicula) {
  return peticion(`${BASE}/peliculas`, {
    method: 'POST',
    body: JSON.stringify(pelicula)
  })
}

export function actualizarPelicula(id, pelicula) {
  return peticion(`${BASE}/peliculas/${id}`, {
    method: 'PUT',
    body: JSON.stringify(pelicula)
  })
}

export function eliminarPelicula(id) {
  return peticion(`${BASE}/peliculas/${id}`, { method: 'DELETE' })
}

export function cambiarEstado(id, estado) {
  return peticion(`${BASE}/peliculas/${id}/estado`, {
    method: 'PATCH',
    body: JSON.stringify({ estado })
  })
}

export function puntuarPelicula(id, puntuacion) {
  return peticion(`${BASE}/peliculas/${id}/puntuacion`, {
    method: 'PATCH',
    body: JSON.stringify({ puntuacion })
  })
}

// Géneros

export function obtenerGeneros() {
  return peticion(`${BASE}/generos`)
}

export function crearGenero(genero) {
  return peticion(`${BASE}/generos`, {
    method: 'POST',
    body: JSON.stringify(genero)
  })
}

export function actualizarGenero(id, genero) {
  return peticion(`${BASE}/generos/${id}`, {
    method: 'PUT',
    body: JSON.stringify(genero)
  })
}

export function eliminarGenero(id) {
  return peticion(`${BASE}/generos/${id}`, { method: 'DELETE' })
}

// Estadísticas

export function obtenerEstadisticas() {
  return peticion(`${BASE}/estadisticas`)
}
