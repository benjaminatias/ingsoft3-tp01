// Funciones de formato utilizadas por la interfaz.

// Convierte un número a texto con coma decimal (formato en español).
export function formatearDecimal(valor) {
  if (valor === null || valor === undefined || Number.isNaN(Number(valor))) {
    return '-'
  }
  return String(Number(valor)).replace('.', ',')
}

// Muestra la puntuación como "8,5 / 10".
// Las películas sin puntuación (pendientes) muestran "-".
export function formatearPuntuacion(puntuacion) {
  if (puntuacion === null || puntuacion === undefined || puntuacion === '') {
    return '-'
  }
  const numero = Number(puntuacion)
  if (Number.isNaN(numero)) {
    return '-'
  }
  return `${formatearDecimal(numero)} / 10`
}

// Muestra el promedio de la colección.
export function formatearPromedio(promedio) {
  if (promedio === null || promedio === undefined) {
    return '-'
  }
  return formatearDecimal(promedio)
}

// Muestra el estado con la primera letra en mayúscula.
export function formatearEstado(estado) {
  if (estado === 'vista') return 'Vista'
  if (estado === 'pendiente') return 'Pendiente'
  return '-'
}

// Devuelve el nombre del género de una película.
export function nombreGenero(pelicula) {
  if (!pelicula || !pelicula.genero || !pelicula.genero.nombre) {
    return '-'
  }
  return pelicula.genero.nombre
}

// Prepara los datos del formulario para enviarlos al backend.
// El título se limpia y la puntuación se envía como número o null.
export function prepararPelicula(formulario) {
  const estado = formulario.estado
  const puntuacionTexto = String(formulario.puntuacion ?? '').trim()
  const tienePuntuacion = estado === 'vista' && puntuacionTexto !== ''

  return {
    titulo: String(formulario.titulo ?? '').trim(),
    anio: Number(formulario.anio),
    generoId: Number(formulario.generoId),
    estado,
    puntuacion: tienePuntuacion ? Number(puntuacionTexto.replace(',', '.')) : null
  }
}

// Convierte una película recibida del backend en los valores del formulario.
export function peliculaAFormulario(pelicula) {
  return {
    titulo: pelicula.titulo ?? '',
    anio: pelicula.anio ? String(pelicula.anio) : '',
    generoId: pelicula.generoId ? String(pelicula.generoId) : '',
    estado: pelicula.estado ?? 'pendiente',
    puntuacion: pelicula.puntuacion === null || pelicula.puntuacion === undefined
      ? ''
      : String(pelicula.puntuacion)
  }
}
