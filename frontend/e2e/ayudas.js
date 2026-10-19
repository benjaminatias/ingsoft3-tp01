import { expect } from '@playwright/test'

// Cuenta de pruebas de QA. No es un secreto: QA es un entorno descartable y esta cuenta solo
// existe ahí. Es UNA sola (no una por corrida) para no llenar la base de usuarios que no se pueden borrar.
export const CUENTA = {
  nombre: 'Bot E2E',
  email: 'bot-e2e@ejemplo.test',
  password: 'clave-e2e-1234'
}

// Devuelve el token de la cuenta de pruebas: la crea si no existe, o inicia sesión si ya existía.
// `base` es la dirección de la api ('' usa el baseURL del test, o sea el front, que la proxea).
export async function tokenDePrueba(request, base = '') {
  const registro = await request.post(`${base}/api/auth/registro`, { data: CUENTA })
  if (registro.status() === 201) {
    return (await registro.json()).token
  }
  expect(registro.status(), 'el registro de la cuenta de pruebas solo puede dar 201 o 409').toBe(409)

  const login = await request.post(`${base}/api/auth/login`, {
    data: { email: CUENTA.email, password: CUENTA.password }
  })
  expect(login.status()).toBe(200)
  return (await login.json()).token
}
