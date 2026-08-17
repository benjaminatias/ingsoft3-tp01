import { formatearPromedio } from '../utils/formato'

export default function Estadisticas({ estadisticas, cargando }) {
  if (cargando && !estadisticas) {
    return <p className="mensaje">Cargando estadísticas...</p>
  }

  if (!estadisticas) {
    return null
  }

  const tarjetas = [
    { etiqueta: 'Total', valor: estadisticas.totalPeliculas },
    { etiqueta: 'Vistas', valor: estadisticas.vistas },
    { etiqueta: 'Pendientes', valor: estadisticas.pendientes },
    { etiqueta: 'Puntuación promedio', valor: formatearPromedio(estadisticas.puntuacionPromedio) }
  ]

  return (
    <section className="estadisticas" aria-label="Estadísticas de la colección">
      {tarjetas.map((tarjeta) => (
        <div className="tarjeta" key={tarjeta.etiqueta}>
          <span className="tarjeta-etiqueta">{tarjeta.etiqueta}</span>
          <span className="tarjeta-valor">{tarjeta.valor}</span>
        </div>
      ))}
    </section>
  )
}
