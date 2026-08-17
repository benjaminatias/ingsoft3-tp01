// Validaciones del formulario. Repiten las reglas del backend para avisar
// al usuario antes de enviar la petición. El backend siempre vuelve a validar.

export const ANIO_MINIMO = 1888
export const TITULO_MAX = 200
export const PUNTUACION_MINIMA = 1
export const PUNTUACION_MAXIMA = 10

// El año máximo permitido nunca se hardcodea.
export function anioMaximo() {
  return new Date().getFullYear() + 5
}

// Valida los datos ya preparados con prepararPelicula().
// Devuelve un mensaje de error o null si todo es correcto.
export function validarPelicula(pelicula) {
  const titulo = String(pelicula.titulo ?? '').trim()
  if (titulo === '') {
    return 'El título es obligatorio.'
  }
  if (titulo.length > TITULO_MAX) {
    return `El título no puede superar los ${TITULO_MAX} caracteres.`
  }

  const anio = Number(pelicula.anio)
  if (!Number.isInteger(anio) || anio < ANIO_MINIMO || anio > anioMaximo()) {
    return `El año debe estar entre ${ANIO_MINIMO} y ${anioMaximo()}.`
  }

  if (!pelicula.generoId) {
    return 'El género es obligatorio.'
  }

  if (pelicula.estado !== 'pendiente' && pelicula.estado !== 'vista') {
    return 'El estado solamente puede ser "pendiente" o "vista".'
  }

  if (pelicula.puntuacion !== null && pelicula.puntuacion !== undefined) {
    if (pelicula.estado !== 'vista') {
      return 'Una película pendiente no puede tener puntuación.'
    }
    return validarPuntuacion(pelicula.puntuacion)
  }

  return null
}

// Valida una puntuación suelta (rango 1 a 10 con un decimal como máximo).
export function validarPuntuacion(puntuacion) {
  const numero = Number(puntuacion)
  if (puntuacion === '' || puntuacion === null || Number.isNaN(numero)) {
    return 'La puntuación no es un número válido.'
  }
  if (numero < PUNTUACION_MINIMA || numero > PUNTUACION_MAXIMA) {
    return `La puntuación debe estar entre ${PUNTUACION_MINIMA} y ${PUNTUACION_MAXIMA}.`
  }
  if (Math.abs(numero * 10 - Math.round(numero * 10)) > 1e-9) {
    return 'La puntuación admite como máximo un decimal.'
  }
  return null
}

// --- Cuentas de usuario ---

export const PASSWORD_MIN = 8
export const PASSWORD_MAX_BYTES = 72 // límite de bcrypt en el backend
export const EMAIL_MAX = 120

// Normaliza el email igual que el backend: sin espacios y en minúsculas.
export function normalizarEmail(email) {
  return String(email ?? '').trim().toLowerCase()
}

// Valida el email. Devuelve un mensaje de error o null.
export function validarEmail(email) {
  const limpio = normalizarEmail(email)
  if (limpio === '') {
    return 'El email es obligatorio.'
  }
  if (limpio.length > EMAIL_MAX) {
    return `El email no puede superar los ${EMAIL_MAX} caracteres.`
  }
  // Comprobación simple: algo@algo.algo, sin espacios.
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(limpio)) {
    return 'El email no tiene un formato válido.'
  }
  return null
}

// Valida la contraseña con las mismas reglas que el backend.
export function validarPassword(password) {
  const valor = String(password ?? '')
  if (valor === '') {
    return 'La contraseña es obligatoria.'
  }
  if (valor.length < PASSWORD_MIN) {
    return `La contraseña debe tener al menos ${PASSWORD_MIN} caracteres.`
  }
  if (new TextEncoder().encode(valor).length > PASSWORD_MAX_BYTES) {
    return `La contraseña no puede superar los ${PASSWORD_MAX_BYTES} bytes.`
  }
  return null
}

// Valida el formulario de inicio de sesión.
export function validarLogin({ email, password }) {
  const errorEmail = validarEmail(email)
  if (errorEmail) {
    return errorEmail
  }
  if (String(password ?? '') === '') {
    return 'La contraseña es obligatoria.'
  }
  return null
}

// Valida el formulario de registro, incluyendo la repetición de la contraseña.
export function validarRegistro({ nombre, email, password, repetirPassword }) {
  const nombreLimpio = String(nombre ?? '').trim()
  if (nombreLimpio === '') {
    return 'El nombre es obligatorio.'
  }
  if (nombreLimpio.length < 2 || nombreLimpio.length > 50) {
    return 'El nombre debe tener entre 2 y 50 caracteres.'
  }

  const errorEmail = validarEmail(email)
  if (errorEmail) {
    return errorEmail
  }

  const errorPassword = validarPassword(password)
  if (errorPassword) {
    return errorPassword
  }

  if (repetirPassword !== undefined && password !== repetirPassword) {
    return 'Las contraseñas no coinciden.'
  }

  return null
}

// Valida el nombre de un género (entre 2 y 50 caracteres).
export function validarNombreGenero(nombre) {
  const limpio = String(nombre ?? '').trim()
  if (limpio === '') {
    return 'El nombre del género es obligatorio.'
  }
  if (limpio.length < 2 || limpio.length > 50) {
    return 'El nombre del género debe tener entre 2 y 50 caracteres.'
  }
  return null
}
