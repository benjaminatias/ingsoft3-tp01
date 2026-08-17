import { useState } from 'react'
import { formatearEstado, formatearPuntuacion, nombreGenero } from '../utils/formato'
import { validarPuntuacion } from '../utils/validacion'

export default function PeliculaItem({ pelicula, onCambiarEstado, onPuntuar, onEditar, onEliminar }) {
  const [puntuando, setPuntuando] = useState(false)
  const [puntuacion, setPuntuacion] = useState('')
  const [error, setError] = useState('')

  const estaVista = pelicula.estado === 'vista'

  function abrirPuntuacion() {
    setPuntuacion(pelicula.puntuacion === null || pelicula.puntuacion === undefined ? '' : String(pelicula.puntuacion))
    setError('')
    setPuntuando(true)
  }

  async function guardarPuntuacion(evento) {
    evento.preventDefault()

    const mensaje = validarPuntuacion(puntuacion)
    if (mensaje) {
      setError(mensaje)
      return
    }

    try {
      await onPuntuar(pelicula.id, Number(String(puntuacion).replace(',', '.')))
      setPuntuando(false)
    } catch (fallo) {
      setError(fallo.message)
    }
  }

  return (
    <tr>
      <td data-etiqueta="Título">{pelicula.titulo}</td>
      <td data-etiqueta="Año">{pelicula.anio}</td>
      <td data-etiqueta="Género">{nombreGenero(pelicula)}</td>
      <td data-etiqueta="Estado">
        <span className={estaVista ? 'etiqueta etiqueta-vista' : 'etiqueta etiqueta-pendiente'}>
          {formatearEstado(pelicula.estado)}
        </span>
      </td>
      <td data-etiqueta="Puntuación">
        {puntuando ? (
          <form className="puntuar" onSubmit={guardarPuntuacion}>
            <input
              type="number"
              step="0.5"
              min="1"
              max="10"
              value={puntuacion}
              onChange={(evento) => setPuntuacion(evento.target.value)}
              aria-label={`Puntuación de ${pelicula.titulo}`}
            />
            <button type="submit" className="boton boton-pequeno boton-primario">
              Guardar
            </button>
            <button type="button" className="boton boton-pequeno" onClick={() => setPuntuando(false)}>
              Cancelar
            </button>
            {error && <span className="mensaje-error">{error}</span>}
          </form>
        ) : (
          formatearPuntuacion(pelicula.puntuacion)
        )}
      </td>
      <td data-etiqueta="Acciones" className="celda-acciones">
        <button
          type="button"
          className="boton boton-pequeno"
          onClick={() => onCambiarEstado(pelicula.id, estaVista ? 'pendiente' : 'vista')}
        >
          {estaVista ? 'Marcar pendiente' : 'Marcar vista'}
        </button>
        {estaVista && !puntuando && (
          <button type="button" className="boton boton-pequeno" onClick={abrirPuntuacion}>
            Puntuar
          </button>
        )}
        <button type="button" className="boton boton-pequeno" onClick={() => onEditar(pelicula)}>
          Editar
        </button>
        <button type="button" className="boton boton-pequeno boton-peligro" onClick={() => onEliminar(pelicula)}>
          Eliminar
        </button>
      </td>
    </tr>
  )
}
