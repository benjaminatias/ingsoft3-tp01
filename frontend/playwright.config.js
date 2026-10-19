import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',                // recoge los DOS archivos; cada job del pipeline corre el suyo por nombre
  timeout: 60_000,                 // tope de CADA test: generoso para las e2e, de sobra para los pedidos a la api
  expect: { timeout: 15_000 },     // tope de CADA aserción: no lo hereda del de arriba (el default es 5 s)
  use: {
    actionTimeout: 15_000,           // tope de CADA click/fill: si el botón no existe, el test falla en 15 s y no en un minuto
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',   // el FRONT. La api tiene la suya: API_BASE_URL (api.spec.js)
    trace: 'on-first-retry',         // la traza se graba en el retry de un fallo (los verdes no la generan)
    screenshot: 'only-on-failure',   // screenshot de cada fallo en el reporte
  },
  retries: 1,                      // 1 retry: absorbe una demora suelta; si un test pasa recién ahí, sale «flaky»
  reporter: [
    ['html', { open: 'never' }],
    ['list'],
    ['json', { outputFile: 'playwright-results.json' }],   // para el resumen de la corrida (cuántos pasaron)
  ],
})
