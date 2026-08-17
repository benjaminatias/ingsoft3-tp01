import { describe, expect, it } from 'vitest'
import {
  anioMaximo,
  validarNombreGenero,
  validarPelicula,
  validarPuntuacion
} from '../src/utils/validacion'

const PELICULA_VALIDA = {
  titulo: 'Interstellar',
  anio: 2014,
  generoId: 4,
  estado: 'vista',
  puntuacion: 9.5
}

describe('validarPelicula', () => {
  it('acepta una película vista con puntuación válida', () => {
    expect(validarPelicula(PELICULA_VALIDA)).toBeNull()
  })

  it('acepta una película pendiente sin puntuación', () => {
    expect(validarPelicula({ ...PELICULA_VALIDA, estado: 'pendiente', puntuacion: null })).toBeNull()
  })

  it('rechaza un título vacío', () => {
    expect(validarPelicula({ ...PELICULA_VALIDA, titulo: '   ' })).toBe('El título es obligatorio.')
  })

  it('rechaza un año fuera de rango', () => {
    expect(validarPelicula({ ...PELICULA_VALIDA, anio: 1800 })).toContain('El año debe estar entre')
    expect(validarPelicula({ ...PELICULA_VALIDA, anio: anioMaximo() + 1 })).toContain('El año debe estar entre')
  })

  it('rechaza una película sin género', () => {
    expect(validarPelicula({ ...PELICULA_VALIDA, generoId: 0 })).toBe('El género es obligatorio.')
  })

  it('rechaza un estado inválido', () => {
    expect(validarPelicula({ ...PELICULA_VALIDA, estado: 'viendo' })).toContain('El estado solamente puede ser')
  })

  it('rechaza una película pendiente con puntuación', () => {
    expect(validarPelicula({ ...PELICULA_VALIDA, estado: 'pendiente' })).toBe(
      'Una película pendiente no puede tener puntuación.'
    )
  })
})

describe('validarPuntuacion', () => {
  it('acepta puntuaciones entre 1 y 10 con un decimal', () => {
    expect(validarPuntuacion(1)).toBeNull()
    expect(validarPuntuacion(7.5)).toBeNull()
    expect(validarPuntuacion(10)).toBeNull()
  })

  it('rechaza una puntuación menor a 1', () => {
    expect(validarPuntuacion(0.5)).toContain('La puntuación debe estar entre')
  })

  it('rechaza una puntuación mayor a 10', () => {
    expect(validarPuntuacion(11)).toContain('La puntuación debe estar entre')
  })

  it('rechaza más de un decimal', () => {
    expect(validarPuntuacion(8.55)).toBe('La puntuación admite como máximo un decimal.')
  })

  it('rechaza un valor vacío', () => {
    expect(validarPuntuacion('')).toBe('La puntuación no es un número válido.')
  })
})

describe('validarNombreGenero', () => {
  it('acepta un nombre válido', () => {
    expect(validarNombreGenero('Western')).toBeNull()
  })

  it('rechaza un nombre vacío', () => {
    expect(validarNombreGenero('  ')).toBe('El nombre del género es obligatorio.')
  })

  it('rechaza un nombre demasiado corto o demasiado largo', () => {
    expect(validarNombreGenero('W')).toContain('entre 2 y 50 caracteres')
    expect(validarNombreGenero('a'.repeat(51))).toContain('entre 2 y 50 caracteres')
  })
})
