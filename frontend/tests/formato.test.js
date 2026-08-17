import { describe, expect, it } from 'vitest'
import {
  formatearDecimal,
  formatearEstado,
  formatearPromedio,
  formatearPuntuacion,
  nombreGenero,
  peliculaAFormulario,
  prepararPelicula
} from '../src/utils/formato'

describe('formatearPuntuacion', () => {
  it('muestra la puntuación sobre 10', () => {
    expect(formatearPuntuacion(8.5)).toBe('8,5 / 10')
    expect(formatearPuntuacion(9)).toBe('9 / 10')
    expect(formatearPuntuacion(10)).toBe('10 / 10')
  })

  it('muestra un guion cuando no hay puntuación', () => {
    expect(formatearPuntuacion(null)).toBe('-')
    expect(formatearPuntuacion(undefined)).toBe('-')
    expect(formatearPuntuacion('')).toBe('-')
  })
})

describe('formatearPromedio', () => {
  it('utiliza coma decimal', () => {
    expect(formatearPromedio(8.2)).toBe('8,2')
    expect(formatearPromedio(7)).toBe('7')
  })

  it('muestra un guion cuando todavía no hay películas puntuadas', () => {
    expect(formatearPromedio(null)).toBe('-')
  })
})

describe('formatearDecimal', () => {
  it('convierte el punto en coma', () => {
    expect(formatearDecimal(9.5)).toBe('9,5')
    expect(formatearDecimal(3)).toBe('3')
  })
})

describe('formatearEstado', () => {
  it('traduce los estados de la API', () => {
    expect(formatearEstado('vista')).toBe('Vista')
    expect(formatearEstado('pendiente')).toBe('Pendiente')
    expect(formatearEstado('otro')).toBe('-')
  })
})

describe('nombreGenero', () => {
  it('devuelve el nombre del género incluido en la película', () => {
    expect(nombreGenero({ genero: { id: 4, nombre: 'Ciencia ficción' } })).toBe('Ciencia ficción')
  })

  it('devuelve un guion cuando la película no trae género', () => {
    expect(nombreGenero({})).toBe('-')
    expect(nombreGenero(null)).toBe('-')
  })
})

describe('prepararPelicula', () => {
  it('limpia el título y convierte los números', () => {
    const resultado = prepararPelicula({
      titulo: '  Interstellar  ',
      anio: '2014',
      generoId: '4',
      estado: 'vista',
      puntuacion: '9.5'
    })

    expect(resultado).toEqual({
      titulo: 'Interstellar',
      anio: 2014,
      generoId: 4,
      estado: 'vista',
      puntuacion: 9.5
    })
  })

  it('envía la puntuación en null cuando la película está pendiente', () => {
    const resultado = prepararPelicula({
      titulo: 'Dune: Part Two',
      anio: '2024',
      generoId: '4',
      estado: 'pendiente',
      puntuacion: '8.5'
    })

    expect(resultado.puntuacion).toBeNull()
  })

  it('acepta la coma decimal', () => {
    const resultado = prepararPelicula({
      titulo: 'Interstellar',
      anio: '2014',
      generoId: '4',
      estado: 'vista',
      puntuacion: '8,5'
    })

    expect(resultado.puntuacion).toBe(8.5)
  })
})

describe('peliculaAFormulario', () => {
  it('transforma una película de la API en valores de formulario', () => {
    const formulario = peliculaAFormulario({
      id: 1,
      titulo: 'Interstellar',
      anio: 2014,
      generoId: 4,
      estado: 'vista',
      puntuacion: 9.5
    })

    expect(formulario).toEqual({
      titulo: 'Interstellar',
      anio: '2014',
      generoId: '4',
      estado: 'vista',
      puntuacion: '9.5'
    })
  })

  it('deja la puntuación vacía cuando es null', () => {
    const formulario = peliculaAFormulario({
      titulo: 'Dune: Part Two',
      anio: 2024,
      generoId: 4,
      estado: 'pendiente',
      puntuacion: null
    })

    expect(formulario.puntuacion).toBe('')
  })
})
