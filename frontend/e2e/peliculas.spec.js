import { expect, test } from '@playwright/test'
import { CUENTA, tokenDePrueba } from './ayudas'

// Entra por la pantalla de acceso, como lo haría una persona. La cuenta de pruebas se asegura por
// la api antes (a través del front, que la proxea): sin eso, el primer login de QA no tendría cuenta.
async function entrar (page, request) {
  await tokenDePrueba(request)
  await page.goto('/')
  await page.getByLabel('Email', { exact: true }).fill(CUENTA.email)
  await page.getByLabel('Contraseña', { exact: true }).fill(CUENTA.password)
  // Hay dos botones «Iniciar sesión» (la pestaña y el envío): el de envío es el que está dentro del <form>.
  await page.locator('form').getByRole('button', { name: 'Iniciar sesión' }).click()
  await expect(page.getByText(CUENTA.email)).toBeVisible()      // la sesión quedó abierta, sin sleeps
}

// El formulario de «Agregar película»: sus campos comparten nombre con los de los filtros (Año, Género…).
function formularioDeAlta (page) {
  return page.locator('section').filter({ has: page.getByRole('heading', { name: 'Agregar película' }) })
}

test('crear una película la muestra en la lista, y eliminarla la saca', async ({ page, request }) => {
  await entrar(page, request)
  const titulo = `e2e ${Date.now()}`                          // único: no choca con otra corrida

  const alta = formularioDeAlta(page)
  await alta.getByLabel('Título').fill(titulo)
  await alta.getByLabel('Año').fill('2014')
  await alta.getByLabel('Género').selectOption({ index: 1 })   // el primer género real (el 0 es «Seleccioná…»)
  await alta.getByRole('button', { name: 'Agregar película' }).click()

  const celda = page.getByRole('cell', { name: titulo, exact: true })
  await expect(celda).toBeVisible()                            // el dato que ESTE test creó

  await page.getByRole('button', { name: `Eliminar ${titulo}` }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Eliminar', exact: true }).click()
  await expect(celda).toHaveCount(0)                           // y el borrado, comprobado
})

test('un año imposible muestra el error y no crea nada', async ({ page, request }) => {
  await entrar(page, request)
  const titulo = `e2e ${Date.now()} invalida`

  const alta = formularioDeAlta(page)
  await alta.getByLabel('Título').fill(titulo)
  await alta.getByLabel('Año').fill('1500')                    // antes del primer cine (1888)
  await alta.getByLabel('Género').selectOption({ index: 1 })
  await alta.getByRole('button', { name: 'Agregar película' }).click()

  await expect(alta.getByRole('alert')).toContainText('El año debe estar entre 1888')   // el usuario ve el error…

  await page.reload()                                          // …y al recargar, la base no tiene nada
  await expect(page.getByText(CUENTA.email)).toBeVisible()
  await expect(page.getByRole('cell', { name: titulo, exact: true })).toHaveCount(0)
})

test('iniciar sesión muestra la colección de la cuenta, y cerrar sesión vuelve al acceso', async ({ page, request }) => {
  await entrar(page, request)
  await expect(page.getByRole('heading', { name: 'Agregar película' })).toBeVisible()   // la colección cargó

  await page.getByRole('button', { name: 'Cerrar sesión' }).click()

  await expect(page.getByText(CUENTA.email)).toHaveCount(0)    // la cuenta desapareció de la pantalla
  await expect(page.getByText('Iniciá sesión para acceder a tu colección.')).toBeVisible()
})
