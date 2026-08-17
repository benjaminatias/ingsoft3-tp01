import { useState } from 'react'
import * as api from '../api/api'
import { normalizarEmail, validarLogin, validarRegistro } from '../utils/validacion'

const FORMULARIO_VACIO = {
  nombre: '',
  email: '',
  password: '',
  repetirPassword: ''
}

// Panel de acceso: permite iniciar sesión o crear una cuenta nueva.
export default function Login({ onAutenticado, aviso }) {
  const [modo, setModo] = useState('login')
  const [formulario, setFormulario] = useState(FORMULARIO_VACIO)
  const [error, setError] = useState('')
  const [enviando, setEnviando] = useState(false)

  const esRegistro = modo === 'registro'

  function cambiar(campo, valor) {
    setFormulario((anterior) => ({ ...anterior, [campo]: valor }))
  }

  function cambiarModo(nuevoModo) {
    setModo(nuevoModo)
    setFormulario(FORMULARIO_VACIO)
    setError('')
  }

  async function enviar(evento) {
    evento.preventDefault()

    const mensaje = esRegistro ? validarRegistro(formulario) : validarLogin(formulario)
    if (mensaje) {
      setError(mensaje)
      return
    }

    setError('')
    setEnviando(true)
    try {
      const sesion = esRegistro
        ? await api.registrarse({
            nombre: formulario.nombre.trim(),
            email: normalizarEmail(formulario.email),
            password: formulario.password
          })
        : await api.iniciarSesion({
            email: normalizarEmail(formulario.email),
            password: formulario.password
          })

      onAutenticado(sesion)
    } catch (fallo) {
      setError(fallo.message)
    } finally {
      setEnviando(false)
    }
  }

  return (
    <div className="pantalla-acceso">
      <div className="panel panel-acceso">
        <h1>Mis Películas</h1>
        <p className="mensaje">
          {esRegistro
            ? 'Creá una cuenta para administrar tu colección.'
            : 'Iniciá sesión para acceder a tu colección.'}
        </p>

        {aviso && <p className="mensaje mensaje-error">{aviso}</p>}

        <div className="pestanas">
          <button
            type="button"
            className={!esRegistro ? 'pestana pestana-activa' : 'pestana'}
            onClick={() => cambiarModo('login')}
          >
            Iniciar sesión
          </button>
          <button
            type="button"
            className={esRegistro ? 'pestana pestana-activa' : 'pestana'}
            onClick={() => cambiarModo('registro')}
          >
            Crear cuenta
          </button>
        </div>

        <form className="formulario-acceso" onSubmit={enviar}>
          {esRegistro && (
            <div className="campo">
              <label htmlFor="nombre">Nombre</label>
              <input
                id="nombre"
                type="text"
                maxLength={50}
                autoComplete="name"
                value={formulario.nombre}
                onChange={(evento) => cambiar('nombre', evento.target.value)}
                placeholder="Tu nombre"
              />
            </div>
          )}

          <div className="campo">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              maxLength={120}
              autoComplete="email"
              value={formulario.email}
              onChange={(evento) => cambiar('email', evento.target.value)}
              placeholder="nombre@ejemplo.com"
            />
          </div>

          <div className="campo">
            <label htmlFor="password">Contraseña</label>
            <input
              id="password"
              type="password"
              autoComplete={esRegistro ? 'new-password' : 'current-password'}
              value={formulario.password}
              onChange={(evento) => cambiar('password', evento.target.value)}
              placeholder={esRegistro ? 'Mínimo 8 caracteres' : 'Tu contraseña'}
            />
          </div>

          {esRegistro && (
            <div className="campo">
              <label htmlFor="repetirPassword">Repetir contraseña</label>
              <input
                id="repetirPassword"
                type="password"
                autoComplete="new-password"
                value={formulario.repetirPassword}
                onChange={(evento) => cambiar('repetirPassword', evento.target.value)}
                placeholder="Repetí la contraseña"
              />
            </div>
          )}

          {error && <p className="mensaje mensaje-error">{error}</p>}

          <button type="submit" className="boton boton-primario boton-ancho" disabled={enviando}>
            {enviando
              ? 'Enviando...'
              : esRegistro
                ? 'Crear cuenta'
                : 'Iniciar sesión'}
          </button>
        </form>
      </div>
    </div>
  )
}
