import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  borrarToken,
  configurarExpiracionSesion,
  guardarToken,
  iniciarSesion,
  obtenerPeliculas
} from '../src/api/api'

// Respuesta falsa con la forma mínima que usa api.js (status, ok y json()).
function respuestaFalsa(status, cuerpo) {
  return {
    status,
    ok: status >= 200 && status < 300,
    json: async () => cuerpo
  }
}

// El entorno de vitest es "node": no existe localStorage. Se reemplaza por un doble en memoria.
function instalarAlmacenFalso() {
  const datos = new Map()
  vi.stubGlobal('localStorage', {
    getItem: (clave) => (datos.has(clave) ? datos.get(clave) : null),
    setItem: (clave, valor) => datos.set(clave, String(valor)),
    removeItem: (clave) => datos.delete(clave)
  })
}

let fetchFalso

beforeEach(() => {
  instalarAlmacenFalso()
  fetchFalso = vi.fn()
  vi.stubGlobal('fetch', fetchFalso)
  configurarExpiracionSesion(null)
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('api con fetch reemplazado por un mock', () => {
  it('descarta los filtros vacíos al armar la URL', async () => {
    // Arrange
    fetchFalso.mockResolvedValue(respuestaFalsa(200, []))

    // Act
    await obtenerPeliculas({ genero: '', anio: 2014, estado: 'vista', busqueda: null })

    // Assert: el mock permite verificar A QUÉ URL se llamó, sin salir a la red
    expect(fetchFalso).toHaveBeenCalledTimes(1)
    expect(fetchFalso.mock.calls[0][0]).toBe('/api/peliculas?anio=2014&estado=vista')
  })

  it('envía el token guardado en la cabecera Authorization', async () => {
    guardarToken('abc123')
    fetchFalso.mockResolvedValue(respuestaFalsa(200, []))

    await obtenerPeliculas()

    const opciones = fetchFalso.mock.calls[0][1]
    expect(opciones.headers.Authorization).toBe('Bearer abc123')
  })

  it('no envía Authorization cuando no hay sesión', async () => {
    fetchFalso.mockResolvedValue(respuestaFalsa(200, []))

    await obtenerPeliculas()

    expect(fetchFalso.mock.calls[0][1].headers.Authorization).toBeUndefined()
  })

  it('guarda el token que devuelve el login', async () => {
    fetchFalso.mockResolvedValue(respuestaFalsa(200, { token: 'nuevo-token' }))

    await iniciarSesion({ email: 'a@b.com', password: '12345678' })
    fetchFalso.mockResolvedValue(respuestaFalsa(200, []))
    await obtenerPeliculas()

    expect(fetchFalso.mock.calls[1][1].headers.Authorization).toBe('Bearer nuevo-token')
  })

  it('devuelve null ante un 204 sin cuerpo', async () => {
    fetchFalso.mockResolvedValue(respuestaFalsa(204, null))

    await expect(obtenerPeliculas()).resolves.toBeNull()
  })
})

describe('errores de la API', () => {
  // Test parametrizado: la misma regla (un status de error se convierte en Error con mensaje y status).
  it.each([
    [400, 'El título es obligatorio.'],
    [404, 'La película no existe.'],
    [409, 'Ya existe una cuenta con ese email.'],
    [500, 'No se pudo completar la operación.']
  ])('convierte el status %i en un Error con el mensaje del backend', async (status, mensaje) => {
    fetchFalso.mockResolvedValue(respuestaFalsa(status, { error: mensaje }))

    await expect(obtenerPeliculas()).rejects.toMatchObject({ message: mensaje, status })
  })

  it('usa un mensaje genérico si el error no trae cuerpo', async () => {
    fetchFalso.mockResolvedValue({ status: 500, ok: false, json: async () => { throw new Error('no es json') } })

    await expect(obtenerPeliculas()).rejects.toThrow('Ocurrió un error inesperado.')
  })

  it('avisa que no hay conexión cuando fetch falla', async () => {
    fetchFalso.mockRejectedValue(new TypeError('Failed to fetch'))

    await expect(obtenerPeliculas()).rejects.toThrow('No se pudo conectar con el servidor.')
  })
})

describe('expiración de la sesión', () => {
  it('un 401 con token cierra la sesión y avisa a la aplicación', async () => {
    const alExpirar = vi.fn()
    configurarExpiracionSesion(alExpirar)
    guardarToken('vencido')
    fetchFalso.mockResolvedValue(respuestaFalsa(401, { error: 'La sesión expiró.' }))

    await expect(obtenerPeliculas()).rejects.toThrow('La sesión expiró.')

    expect(alExpirar).toHaveBeenCalledTimes(1)
    // Con la sesión cerrada, la siguiente petición ya no lleva token.
    fetchFalso.mockResolvedValue(respuestaFalsa(200, []))
    await obtenerPeliculas()
    expect(fetchFalso.mock.calls[1][1].headers.Authorization).toBeUndefined()
  })

  it('un 401 sin token (login fallido) NO cuenta como expiración', async () => {
    const alExpirar = vi.fn()
    configurarExpiracionSesion(alExpirar)
    borrarToken()
    fetchFalso.mockResolvedValue(respuestaFalsa(401, { error: 'Email o contraseña incorrectos.' }))

    await expect(iniciarSesion({ email: 'a@b.com', password: 'mal' })).rejects.toThrow(
      'Email o contraseña incorrectos.'
    )

    expect(alExpirar).not.toHaveBeenCalled()
  })
})
