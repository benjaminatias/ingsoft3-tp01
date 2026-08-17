import { useEffect, useState } from 'react'
import { peliculaAFormulario, prepararPelicula } from '../utils/formato'
import { validarPelicula } from '../utils/validacion'

const FORMULARIO_VACIO = {
  titulo: '',
  anio: '',
  generoId: '',
  estado: 'pendiente',
  puntuacion: ''
}

// Formulario utilizado tanto para agregar como para editar una película.
export default function FormularioPelicula({ generos, peliculaEditando, onGuardar, onCancelar }) {
  const [formulario, setFormulario] = useState(FORMULARIO_VACIO)
  const [error, setError] = useState('')
  const [enviando, setEnviando] = useState(false)

  useEffect(() => {
    if (peliculaEditando) {
      setFormulario(peliculaAFormulario(peliculaEditando))
    } else {
      setFormulario(FORMULARIO_VACIO)
    }
    setError('')
  }, [peliculaEditando])

  function cambiar(campo, valor) {
    setFormulario((anterior) => {
      const siguiente = { ...anterior, [campo]: valor }
      // Una película pendiente nunca tiene puntuación.
      if (campo === 'estado' && valor === 'pendiente') {
        siguiente.puntuacion = ''
      }
      return siguiente
    })
  }

  async function enviar(evento) {
    evento.preventDefault()

    const pelicula = prepararPelicula(formulario)
    const mensaje = validarPelicula(pelicula)
    if (mensaje) {
      setError(mensaje)
      return
    }

    setError('')
    setEnviando(true)
    try {
      await onGuardar(pelicula)
      if (!peliculaEditando) {
        setFormulario(FORMULARIO_VACIO)
      }
    } catch (fallo) {
      setError(fallo.message)
    } finally {
      setEnviando(false)
    }
  }

  return (
    <section className="panel">
      <h2>{peliculaEditando ? 'Editar película' : 'Agregar película'}</h2>

      <form className="formulario" onSubmit={enviar}>
        <div className="campo">
          <label htmlFor="titulo">Título</label>
          <input
            id="titulo"
            type="text"
            maxLength={200}
            value={formulario.titulo}
            onChange={(evento) => cambiar('titulo', evento.target.value)}
            placeholder="Interstellar"
          />
        </div>

        <div className="campo">
          <label htmlFor="anio">Año</label>
          <input
            id="anio"
            type="number"
            value={formulario.anio}
            onChange={(evento) => cambiar('anio', evento.target.value)}
            placeholder="2014"
          />
        </div>

        <div className="campo">
          <label htmlFor="generoId">Género</label>
          <select
            id="generoId"
            value={formulario.generoId}
            onChange={(evento) => cambiar('generoId', evento.target.value)}
          >
            <option value="">Seleccioná un género</option>
            {generos.map((genero) => (
              <option key={genero.id} value={genero.id}>
                {genero.nombre}
              </option>
            ))}
          </select>
        </div>

        <div className="campo">
          <label htmlFor="estado">Estado</label>
          <select
            id="estado"
            value={formulario.estado}
            onChange={(evento) => cambiar('estado', evento.target.value)}
          >
            <option value="pendiente">Pendiente</option>
            <option value="vista">Vista</option>
          </select>
        </div>

        <div className="campo">
          <label htmlFor="puntuacion">Puntuación</label>
          <input
            id="puntuacion"
            type="number"
            step="0.5"
            min="1"
            max="10"
            value={formulario.puntuacion}
            onChange={(evento) => cambiar('puntuacion', evento.target.value)}
            disabled={formulario.estado !== 'vista'}
            placeholder={formulario.estado === 'vista' ? '8.5' : 'Solo para vistas'}
          />
        </div>

        {error && <p className="mensaje mensaje-error">{error}</p>}

        <div className="acciones-formulario">
          <button type="submit" className="boton boton-primario" disabled={enviando}>
            {peliculaEditando ? 'Guardar cambios' : 'Agregar película'}
          </button>
          {peliculaEditando && (
            <button type="button" className="boton" onClick={onCancelar}>
              Cancelar
            </button>
          )}
        </div>
      </form>
    </section>
  )
}
