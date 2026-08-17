import { useState } from 'react'

export const FILTROS_VACIOS = {
  titulo: '',
  estado: '',
  generoId: '',
  anio: '',
  puntuacionMin: ''
}

export default function FiltrosPeliculas({ generos, onFiltrar }) {
  const [filtros, setFiltros] = useState(FILTROS_VACIOS)

  function cambiar(campo, valor) {
    setFiltros((anterior) => ({ ...anterior, [campo]: valor }))
  }

  function filtrar(evento) {
    evento.preventDefault()
    onFiltrar(filtros)
  }

  function limpiar() {
    setFiltros(FILTROS_VACIOS)
    onFiltrar(FILTROS_VACIOS)
  }

  return (
    <section className="panel">
      <h2>Filtros</h2>

      <form className="formulario" onSubmit={filtrar}>
        <div className="campo">
          <label htmlFor="buscar">Buscar</label>
          <input
            id="buscar"
            type="text"
            value={filtros.titulo}
            onChange={(evento) => cambiar('titulo', evento.target.value)}
            placeholder="Título de la película"
          />
        </div>

        <div className="campo">
          <label htmlFor="filtro-estado">Estado</label>
          <select
            id="filtro-estado"
            value={filtros.estado}
            onChange={(evento) => cambiar('estado', evento.target.value)}
          >
            <option value="">Todos</option>
            <option value="pendiente">Pendiente</option>
            <option value="vista">Vista</option>
          </select>
        </div>

        <div className="campo">
          <label htmlFor="filtro-genero">Género</label>
          <select
            id="filtro-genero"
            value={filtros.generoId}
            onChange={(evento) => cambiar('generoId', evento.target.value)}
          >
            <option value="">Todos</option>
            {generos.map((genero) => (
              <option key={genero.id} value={genero.id}>
                {genero.nombre}
              </option>
            ))}
          </select>
        </div>

        <div className="campo">
          <label htmlFor="filtro-anio">Año</label>
          <input
            id="filtro-anio"
            type="number"
            value={filtros.anio}
            onChange={(evento) => cambiar('anio', evento.target.value)}
            placeholder="2014"
          />
        </div>

        <div className="campo">
          <label htmlFor="filtro-puntuacion">Puntuación mínima</label>
          <input
            id="filtro-puntuacion"
            type="number"
            step="0.5"
            min="1"
            max="10"
            value={filtros.puntuacionMin}
            onChange={(evento) => cambiar('puntuacionMin', evento.target.value)}
            placeholder="8"
          />
        </div>

        <div className="acciones-formulario">
          <button type="submit" className="boton boton-primario">
            Filtrar
          </button>
          <button type="button" className="boton" onClick={limpiar}>
            Limpiar
          </button>
        </div>
      </form>
    </section>
  )
}
