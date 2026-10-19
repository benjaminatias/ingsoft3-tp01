import { expect, test } from '@playwright/test'
import { tokenDePrueba } from './ayudas'

// La api tiene su propia dirección: NO es la del front. Sin barra al final.
const API = process.env.API_BASE_URL || 'http://localhost:8080'

// Cabeceras con la sesión de la cuenta de pruebas.
async function cabeceras (request) {
  const token = await tokenDePrueba(request, API)
  return { Authorization: `Bearer ${token}` }
}

// Pide las películas (opcionalmente filtradas por título) y las devuelve ya convertidas.
async function listar (request, headers, titulo = '') {
  const consulta = titulo ? `?titulo=${encodeURIComponent(titulo)}` : ''
  const r = await request.get(`${API}/api/peliculas${consulta}`, { headers })
  expect(r.status()).toBe(200)
  return await r.json()
}

// Un género que ya existe en la base (los crea la propia api al arrancar).
async function primerGenero (request, headers) {
  const r = await request.get(`${API}/api/generos`, { headers })
  expect(r.status()).toBe(200)
  const generos = await r.json()
  expect(generos.length).toBeGreaterThan(0)
  return generos[0].id
}

test('el alta guarda en la base de verdad y el borrado la saca', async ({ request }) => {
  const headers = await cabeceras(request)
  const generoId = await primerGenero(request, headers)
  const titulo = `api ${Date.now()}`                       // único: no choca con otra corrida

  const alta = await request.post(`${API}/api/peliculas`, {
    headers, data: { titulo, anio: 2014, generoId, estado: 'pendiente' }
  })
  expect(alta.status()).toBe(201)                          // la api dice «creada»
  const creada = await alta.json()

  let borrada = false
  try {
    const lista = await listar(request, headers, titulo)   // y la BASE la devuelve: no es un doble
    expect(lista.some(p => p.id === creada.id && p.titulo === titulo)).toBe(true)

    const una = await request.get(`${API}/api/peliculas/${creada.id}`, { headers })
    expect(una.status()).toBe(200)

    const borrado = await request.delete(`${API}/api/peliculas/${creada.id}`, { headers })
    expect(borrado.status()).toBe(204)                     // «borrada, no hay nada que devolver»
    borrada = true

    const despues = await listar(request, headers, titulo) // y el borrado, comprobado
    expect(despues.some(p => p.id === creada.id)).toBe(false)
    const ya = await request.get(`${API}/api/peliculas/${creada.id}`, { headers })
    expect(ya.status()).toBe(404)
  } finally {
    // si una aserción de arriba falló, igual se limpia lo que se creó
    if (!borrada) {
      await request.delete(`${API}/api/peliculas/${creada.id}`, { headers })
    }
  }
})

test('un título vacío lo rechaza la api con 400 y no crea nada', async ({ request }) => {
  const headers = await cabeceras(request)
  const generoId = await primerGenero(request, headers)

  const alta = await request.post(`${API}/api/peliculas`, {
    headers, data: { titulo: '', anio: 2014, generoId, estado: 'pendiente' }
  })
  expect(alta.status()).toBe(400)                          // el contrato del error
  expect((await alta.json()).error).toBe('El título es obligatorio.')

  // Y de verdad no se creó nada: no hay ninguna película con el título vacío.
  // (No comparo la cantidad antes y después: QA lo comparten otras corridas y la cuenta podría moverse.)
  const despues = await listar(request, headers)
  expect(despues.some(p => p.titulo.trim() === '')).toBe(false)
})

test('un género con películas no se puede borrar: la base lo impide (409) y sigue existiendo', async ({ request }) => {
  const headers = await cabeceras(request)
  const nombre = `api genero ${Date.now()}`                // único
  const tituloPelicula = `api pelicula ${Date.now()}`

  const altaGenero = await request.post(`${API}/api/generos`, { headers, data: { nombre } })
  expect(altaGenero.status()).toBe(201)
  const genero = await altaGenero.json()

  let pelicula = null
  let generoBorrado = false
  try {
    const altaPelicula = await request.post(`${API}/api/peliculas`, {
      headers, data: { titulo: tituloPelicula, anio: 2014, generoId: genero.id, estado: 'pendiente' }
    })
    expect(altaPelicula.status()).toBe(201)
    pelicula = await altaPelicula.json()

    // La regla vive en una consulta a la base: un doble no la ve.
    const intento = await request.delete(`${API}/api/generos/${genero.id}`, { headers })
    expect(intento.status()).toBe(409)

    const sigue = await request.get(`${API}/api/generos/${genero.id}`, { headers })
    expect(sigue.status()).toBe(200)

    // Sin películas, ahora sí se puede borrar: primero la película, después el género.
    const borraPelicula = await request.delete(`${API}/api/peliculas/${pelicula.id}`, { headers })
    expect(borraPelicula.status()).toBe(204)
    pelicula = null

    const borraGenero = await request.delete(`${API}/api/generos/${genero.id}`, { headers })
    expect(borraGenero.status()).toBe(204)
    generoBorrado = true

    const ya = await request.get(`${API}/api/generos/${genero.id}`, { headers })
    expect(ya.status()).toBe(404)
  } finally {
    // limpieza aunque algo haya fallado: la película antes que el género
    if (pelicula) {
      await request.delete(`${API}/api/peliculas/${pelicula.id}`, { headers })
    }
    if (!generoBorrado) {
      await request.delete(`${API}/api/generos/${genero.id}`, { headers })
    }
  }
})
