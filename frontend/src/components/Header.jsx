export default function Header({ usuario, onCerrarSesion }) {
  return (
    <header className="encabezado">
      <div className="encabezado-titulo">
        <h1>Mis Películas</h1>
        <p>Colección personal: registrá, puntuá y filtrá tus películas.</p>
      </div>

      {usuario && (
        <div className="sesion">
          <span className="sesion-usuario">
            {usuario.nombre}
            <span className="sesion-email">{usuario.email}</span>
          </span>
          <button type="button" className="boton boton-pequeno" onClick={onCerrarSesion}>
            Cerrar sesión
          </button>
        </div>
      )}
    </header>
  )
}
