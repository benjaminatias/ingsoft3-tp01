import PeliculaItem from './PeliculaItem'

export default function ListaPeliculas({ peliculas, cargando, error, onCambiarEstado, onPuntuar, onEditar, onEliminar }) {
  return (
    <section className="panel">
      <h2>Películas</h2>

      {cargando && <p className="mensaje">Cargando películas...</p>}

      {!cargando && error && <p className="mensaje mensaje-error">{error}</p>}

      {!cargando && !error && peliculas.length === 0 && (
        <p className="mensaje">No hay películas registradas.</p>
      )}

      {!cargando && !error && peliculas.length > 0 && (
        <div className="tabla-contenedor">
          <table className="tabla">
            <thead>
              <tr>
                <th>Título</th>
                <th>Año</th>
                <th>Género</th>
                <th>Estado</th>
                <th>Puntuación</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {peliculas.map((pelicula) => (
                <PeliculaItem
                  key={pelicula.id}
                  pelicula={pelicula}
                  onCambiarEstado={onCambiarEstado}
                  onPuntuar={onPuntuar}
                  onEditar={onEditar}
                  onEliminar={onEliminar}
                />
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}
