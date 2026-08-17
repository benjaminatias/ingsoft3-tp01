import { useState } from 'react'
import { validarNombreGenero } from '../utils/validacion'

// Sección pequeña para administrar los géneros de la colección.
export default function Generos({ generos, mensajeError, onCrear, onActualizar, onEliminar }) {
  const [abierto, setAbierto] = useState(false)
  const [nuevoNombre, setNuevoNombre] = useState('')
  const [editandoId, setEditandoId] = useState(null)
  const [nombreEditado, setNombreEditado] = useState('')
  const [error, setError] = useState('')

  async function agregar(evento) {
    evento.preventDefault()

    const mensaje = validarNombreGenero(nuevoNombre)
    if (mensaje) {
      setError(mensaje)
      return
    }

    try {
      await onCrear({ nombre: nuevoNombre.trim() })
      setNuevoNombre('')
      setError('')
    } catch (fallo) {
      setError(fallo.message)
    }
  }

  function empezarEdicion(genero) {
    setEditandoId(genero.id)
    setNombreEditado(genero.nombre)
    setError('')
  }

  async function guardarEdicion(evento) {
    evento.preventDefault()

    const mensaje = validarNombreGenero(nombreEditado)
    if (mensaje) {
      setError(mensaje)
      return
    }

    try {
      await onActualizar(editandoId, { nombre: nombreEditado.trim() })
      setEditandoId(null)
      setError('')
    } catch (fallo) {
      setError(fallo.message)
    }
  }

  async function eliminar(genero) {
    try {
      await onEliminar(genero)
      setError('')
    } catch (fallo) {
      setError(fallo.message)
    }
  }

  return (
    <section className="panel">
      <div className="panel-encabezado">
        <h2>Administrar géneros</h2>
        <button type="button" className="boton boton-pequeno" onClick={() => setAbierto(!abierto)}>
          {abierto ? 'Ocultar' : 'Mostrar'}
        </button>
      </div>

      {abierto && (
        <div>
          <form className="formulario-linea" onSubmit={agregar}>
            <input
              type="text"
              value={nuevoNombre}
              maxLength={50}
              onChange={(evento) => setNuevoNombre(evento.target.value)}
              placeholder="Nuevo género"
              aria-label="Nombre del nuevo género"
            />
            <button type="submit" className="boton boton-primario boton-pequeno">
              Agregar
            </button>
          </form>

          {(error || mensajeError) && (
            <p className="mensaje mensaje-error">{error || mensajeError}</p>
          )}

          {generos.length === 0 ? (
            <p className="mensaje">No hay géneros registrados.</p>
          ) : (
            <ul className="lista-generos">
              {generos.map((genero) => (
                <li key={genero.id}>
                  {editandoId === genero.id ? (
                    <form className="formulario-linea" onSubmit={guardarEdicion}>
                      <input
                        type="text"
                        value={nombreEditado}
                        maxLength={50}
                        onChange={(evento) => setNombreEditado(evento.target.value)}
                        aria-label={`Nuevo nombre para ${genero.nombre}`}
                      />
                      <button type="submit" className="boton boton-pequeno boton-primario">
                        Guardar
                      </button>
                      <button type="button" className="boton boton-pequeno" onClick={() => setEditandoId(null)}>
                        Cancelar
                      </button>
                    </form>
                  ) : (
                    <>
                      <span>{genero.nombre}</span>
                      <span className="acciones-genero">
                        <button type="button" className="boton boton-pequeno" onClick={() => empezarEdicion(genero)}>
                          Editar
                        </button>
                        <button
                          type="button"
                          className="boton boton-pequeno boton-peligro"
                          onClick={() => eliminar(genero)}
                        >
                          Eliminar
                        </button>
                      </span>
                    </>
                  )}
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </section>
  )
}
