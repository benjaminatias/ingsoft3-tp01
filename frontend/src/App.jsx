import { useEffect, useState } from 'react'
import * as api from './api/api'
import Header from './components/Header'
import Estadisticas from './components/Estadisticas'
import FormularioPelicula from './components/FormularioPelicula'
import FiltrosPeliculas, { FILTROS_VACIOS } from './components/FiltrosPeliculas'
import ListaPeliculas from './components/ListaPeliculas'
import Generos from './components/Generos'
import ModalConfirmacion from './components/ModalConfirmacion'
import Login from './components/Login'

export default function App() {
  const [usuario, setUsuario] = useState(null)
  const [verificandoSesion, setVerificandoSesion] = useState(true)

  const [peliculas, setPeliculas] = useState([])
  const [generos, setGeneros] = useState([])
  const [estadisticas, setEstadisticas] = useState(null)
  const [filtros, setFiltros] = useState(FILTROS_VACIOS)

  const [cargando, setCargando] = useState(true)
  const [errorPeliculas, setErrorPeliculas] = useState('')
  const [mensaje, setMensaje] = useState(null)
  const [errorGeneros, setErrorGeneros] = useState('')

  const [peliculaEditando, setPeliculaEditando] = useState(null)
  const [confirmacion, setConfirmacion] = useState(null)

  // Al abrir la aplicación se comprueba si el token guardado sigue siendo válido.
  useEffect(() => {
    api.configurarExpiracionSesion(() => {
      setUsuario(null)
      setMensaje({ texto: 'Tu sesión expiró. Volvé a iniciar sesión.', tipo: 'error' })
    })

    async function recuperarSesion() {
      if (!api.obtenerToken()) {
        setVerificandoSesion(false)
        return
      }
      try {
        const perfil = await api.obtenerPerfil()
        setUsuario(perfil.usuario)
      } catch {
        api.cerrarSesion()
        setUsuario(null)
      } finally {
        setVerificandoSesion(false)
      }
    }

    recuperarSesion()
  }, [])

  // Los datos se cargan recién cuando hay sesión iniciada.
  useEffect(() => {
    if (!usuario) {
      return
    }
    cargarGeneros()
    cargarPeliculas(FILTROS_VACIOS)
    cargarEstadisticas()
  }, [usuario])

  function avisar(texto, tipo = 'exito') {
    setMensaje({ texto, tipo })
  }

  function autenticado(sesion) {
    setMensaje(null)
    setFiltros(FILTROS_VACIOS)
    setUsuario(sesion.usuario)
  }

  function cerrarSesion() {
    api.cerrarSesion()
    setUsuario(null)
    setPeliculas([])
    setGeneros([])
    setEstadisticas(null)
    setPeliculaEditando(null)
    setConfirmacion(null)
    setMensaje(null)
  }

  async function cargarPeliculas(filtrosAplicados) {
    setCargando(true)
    setErrorPeliculas('')
    try {
      const datos = await api.obtenerPeliculas(filtrosAplicados)
      setPeliculas(datos || [])
    } catch (fallo) {
      setPeliculas([])
      setErrorPeliculas(fallo.message || 'No se pudieron obtener las películas.')
    } finally {
      setCargando(false)
    }
  }

  async function cargarGeneros() {
    try {
      const datos = await api.obtenerGeneros()
      setGeneros(datos || [])
    } catch {
      setErrorGeneros('No se pudieron obtener los géneros.')
    }
  }

  async function cargarEstadisticas() {
    try {
      const datos = await api.obtenerEstadisticas()
      setEstadisticas(datos)
    } catch {
      setEstadisticas(null)
    }
  }

  // Después de cada operación se vuelve a consultar el backend,
  // de manera que la interfaz refleje siempre el estado confirmado.
  async function refrescar() {
    await Promise.all([cargarPeliculas(filtros), cargarEstadisticas()])
  }

  async function guardarPelicula(pelicula) {
    if (peliculaEditando) {
      await api.actualizarPelicula(peliculaEditando.id, pelicula)
      setPeliculaEditando(null)
      avisar('Película actualizada correctamente.')
    } else {
      await api.crearPelicula(pelicula)
      avisar('Película agregada correctamente.')
    }
    await refrescar()
  }

  async function cambiarEstado(id, estado) {
    try {
      await api.cambiarEstado(id, estado)
      await refrescar()
      avisar(estado === 'vista' ? 'Película marcada como vista.' : 'Película marcada como pendiente.')
    } catch (fallo) {
      avisar(fallo.message, 'error')
    }
  }

  async function puntuar(id, puntuacion) {
    await api.puntuarPelicula(id, puntuacion)
    await refrescar()
    avisar('Puntuación guardada correctamente.')
  }

  function pedirConfirmacionPelicula(pelicula) {
    setConfirmacion({
      mensaje: `¿Seguro que desea eliminar "${pelicula.titulo}"?`,
      accion: async () => {
        try {
          await api.eliminarPelicula(pelicula.id)
          if (peliculaEditando && peliculaEditando.id === pelicula.id) {
            setPeliculaEditando(null)
          }
          await refrescar()
          avisar('Película eliminada correctamente.')
        } catch (fallo) {
          avisar(fallo.message, 'error')
        }
      }
    })
  }

  function pedirConfirmacionGenero(genero) {
    setConfirmacion({
      mensaje: `¿Seguro que desea eliminar el género "${genero.nombre}"?`,
      accion: async () => {
        try {
          await api.eliminarGenero(genero.id)
          await cargarGeneros()
          setErrorGeneros('')
          avisar('Género eliminado correctamente.')
        } catch (fallo) {
          setErrorGeneros(fallo.message)
        }
      }
    })
  }

  async function crearGenero(genero) {
    await api.crearGenero(genero)
    await cargarGeneros()
    avisar('Género agregado correctamente.')
  }

  async function actualizarGenero(id, genero) {
    await api.actualizarGenero(id, genero)
    await cargarGeneros()
    await refrescar()
    avisar('Género actualizado correctamente.')
  }

  function aplicarFiltros(nuevosFiltros) {
    setFiltros(nuevosFiltros)
    cargarPeliculas(nuevosFiltros)
  }

  async function confirmar() {
    const accion = confirmacion.accion
    setConfirmacion(null)
    await accion()
  }

  if (verificandoSesion) {
    return (
      <div className="contenedor">
        <p className="mensaje">Cargando...</p>
      </div>
    )
  }

  if (!usuario) {
    return <Login onAutenticado={autenticado} aviso={mensaje && mensaje.tipo === 'error' ? mensaje.texto : ''} />
  }

  return (
    <div className="contenedor">
      <Header usuario={usuario} onCerrarSesion={cerrarSesion} />

      <Estadisticas estadisticas={estadisticas} cargando={cargando} />

      {mensaje && (
        <p className={mensaje.tipo === 'error' ? 'mensaje mensaje-error' : 'mensaje mensaje-exito'}>
          {mensaje.texto}
        </p>
      )}

      <FormularioPelicula
        generos={generos}
        peliculaEditando={peliculaEditando}
        onGuardar={guardarPelicula}
        onCancelar={() => setPeliculaEditando(null)}
      />

      <FiltrosPeliculas generos={generos} onFiltrar={aplicarFiltros} />

      <ListaPeliculas
        peliculas={peliculas}
        cargando={cargando}
        error={errorPeliculas}
        onCambiarEstado={cambiarEstado}
        onPuntuar={puntuar}
        onEditar={setPeliculaEditando}
        onEliminar={pedirConfirmacionPelicula}
      />

      <Generos
        generos={generos}
        mensajeError={errorGeneros}
        onCrear={crearGenero}
        onActualizar={actualizarGenero}
        onEliminar={pedirConfirmacionGenero}
      />

      {confirmacion && (
        <ModalConfirmacion
          mensaje={confirmacion.mensaje}
          onConfirmar={confirmar}
          onCancelar={() => setConfirmacion(null)}
        />
      )}
    </div>
  )
}
