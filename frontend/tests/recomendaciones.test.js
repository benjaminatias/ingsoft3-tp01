import { describe, expect, it } from 'vitest'
import { decadaDe, etiquetaDePelicula, nivelDePuntuacion } from '../src/utils/recomendaciones'

describe('nivelDePuntuacion', () => {
  it.each([
    [null, 'sin-puntuar'],
    [undefined, 'sin-puntuar'],
    ['', 'sin-puntuar'],
    ['abc', 'invalida'],
    [0.5, 'invalida'],
    [11, 'invalida'],
    [9, 'excelente'],
    [10, 'excelente'],
    [7, 'buena'],
    [8.9, 'buena'],
    [5, 'regular'],
    [6.9, 'regular'],
    [1, 'mala'],
    [4.9, 'mala']
  ])('puntuación %s -> %s', (puntuacion, esperado) => {
    expect(nivelDePuntuacion(puntuacion)).toBe(esperado)
  })
})

describe('decadaDe', () => {
  it.each([
    [1994, 'años 1990'],
    [1900, 'años 1900'],
    [1999, 'años 1990'],
    [2000, 'años 2000'],
    [2014, 'años 2010'],
    [2025, 'años 2020']
  ])('año %i -> %s', (anio, esperado) => {
    expect(decadaDe(anio)).toBe(esperado)
  })

  it('rechaza años que no son enteros o anteriores al cine', () => {
    expect(decadaDe('abc')).toBe('sin-década')
    expect(decadaDe(1800)).toBe('sin-década')
    expect(decadaDe(1999.5)).toBe('sin-década')
  })
})

describe('etiquetaDePelicula', () => {
  it('devuelve vacío si no hay película', () => {
    expect(etiquetaDePelicula(null, 2026)).toBe('')
  })

  it('marca como clásico pendiente a una película pendiente de hace 20 años o más', () => {
    expect(etiquetaDePelicula({ estado: 'pendiente', anio: 2006 }, 2026)).toBe('clásico pendiente')
  })

  it('marca como por ver a una pendiente reciente', () => {
    expect(etiquetaDePelicula({ estado: 'pendiente', anio: 2007 }, 2026)).toBe('por ver')
  })

  it('marca como imperdible a una vista con puntuación excelente', () => {
    expect(etiquetaDePelicula({ estado: 'vista', anio: 2014, puntuacion: 9.5 }, 2026)).toBe(
      'imperdible (años 2010)'
    )
  })

  it('marca para olvidar a una vista con puntuación mala', () => {
    expect(etiquetaDePelicula({ estado: 'vista', anio: 1999, puntuacion: 3 }, 2026)).toBe('para olvidar')
  })

  it('usa el nivel y la década para una vista con puntuación intermedia', () => {
    expect(etiquetaDePelicula({ estado: 'vista', anio: 1994, puntuacion: 8 }, 2026)).toBe('buena (años 1990)')
  })
})
