// Clasificaciones para mostrar en la colección: nivel de la puntuación,
// década de estreno y una etiqueta que resume una película.

const ANIO_PRIMERA_PELICULA = 1888

// Traduce una puntuación (1 a 10) a un nivel legible.
export function nivelDePuntuacion(puntuacion) {
  if (puntuacion === null || puntuacion === undefined || puntuacion === '') {
    return 'sin-puntuar'
  }
  const numero = Number(puntuacion)
  if (Number.isNaN(numero) || numero < 1 || numero > 10) {
    return 'invalida'
  }
  if (numero >= 9) {
    return 'excelente'
  }
  if (numero >= 7) {
    return 'buena'
  }
  if (numero >= 5) {
    return 'regular'
  }
  return 'mala'
}

// Devuelve la década de un año, por ejemplo 1994 -> "años 1990".
export function decadaDe(anio) {
  const numero = Number(anio)
  if (!Number.isInteger(numero) || numero < ANIO_PRIMERA_PELICULA) {
    return 'sin-década'
  }
  const siglo = numero >= 2000 ? 2000 : 1900
  const decada = Math.floor((numero - siglo) / 10) * 10
  return `años ${siglo + decada}`
}

// Arma la etiqueta que se muestra junto a una película.
export function etiquetaDePelicula(pelicula, anioActual) {
  if (!pelicula) {
    return ''
  }
  if (pelicula.estado === 'pendiente') {
    const espera = anioActual - Number(pelicula.anio)
    if (espera >= 20) {
      return 'clásico pendiente'
    }
    return 'por ver'
  }
  const nivel = nivelDePuntuacion(pelicula.puntuacion)
  if (nivel === 'excelente') {
    return `imperdible (${decadaDe(pelicula.anio)})`
  }
  if (nivel === 'mala') {
    return 'para olvidar'
  }
  return `${nivel} (${decadaDe(pelicula.anio)})`
}
