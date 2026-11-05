// Resumen de una colección de películas para mostrar en la cabecera.

// Cuenta cuántas películas hay por estado y calcula el porcentaje vistas.
export function resumirColeccion(peliculas) {
  if (!Array.isArray(peliculas) || peliculas.length === 0) {
    return { total: 0, vistas: 0, pendientes: 0, porcentajeVistas: 0, mensaje: 'Colección vacía' }
  }

  let vistas = 0
  let pendientes = 0
  for (const pelicula of peliculas) {
    if (pelicula.estado === 'vista') {
      vistas += 1
    } else if (pelicula.estado === 'pendiente') {
      pendientes += 1
    }
  }

  const total = peliculas.length
  const porcentajeVistas = Math.round((vistas / total) * 100)

  let mensaje
  if (porcentajeVistas === 100) {
    mensaje = 'Colección completa'
  } else if (porcentajeVistas >= 50) {
    mensaje = 'Más de la mitad vista'
  } else if (porcentajeVistas > 0) {
    mensaje = 'Queda mucho por ver'
  } else {
    mensaje = 'Todavía no viste ninguna'
  }

  return { total, vistas, pendientes, porcentajeVistas, mensaje }
}

// Devuelve el título de la película mejor puntuada, o null si ninguna tiene puntuación.
export function mejorPuntuada(peliculas) {
  if (!Array.isArray(peliculas)) {
    return null
  }
  let mejor = null
  for (const pelicula of peliculas) {
    if (pelicula.puntuacion === null || pelicula.puntuacion === undefined) {
      continue
    }
    if (mejor === null || pelicula.puntuacion > mejor.puntuacion) {
      mejor = pelicula
    }
  }
  return mejor ? mejor.titulo : null
}
