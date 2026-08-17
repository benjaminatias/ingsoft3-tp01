import { describe, expect, it } from 'vitest'
import {
  normalizarEmail,
  validarEmail,
  validarLogin,
  validarPassword,
  validarRegistro
} from '../src/utils/validacion'

describe('normalizarEmail', () => {
  it('quita espacios y pasa a minúsculas', () => {
    expect(normalizarEmail('  Usuario@Ejemplo.COM ')).toBe('usuario@ejemplo.com')
  })
})

describe('validarEmail', () => {
  it('acepta un email correcto', () => {
    expect(validarEmail('usuario@ejemplo.com')).toBeNull()
    expect(validarEmail('  Usuario@Ejemplo.COM  ')).toBeNull()
  })

  it('rechaza emails inválidos', () => {
    expect(validarEmail('')).toBe('El email es obligatorio.')
    expect(validarEmail('sin-arroba')).toBe('El email no tiene un formato válido.')
    expect(validarEmail('usuario@sinpunto')).toBe('El email no tiene un formato válido.')
    expect(validarEmail('con espacio@ejemplo.com')).toBe('El email no tiene un formato válido.')
  })

  it('rechaza un email demasiado largo', () => {
    expect(validarEmail(`${'a'.repeat(115)}@ejemplo.com`)).toContain('no puede superar')
  })
})

describe('validarPassword', () => {
  it('acepta una contraseña de al menos 8 caracteres', () => {
    expect(validarPassword('12345678')).toBeNull()
    expect(validarPassword('contraseña-larga')).toBeNull()
  })

  it('rechaza una contraseña vacía o corta', () => {
    expect(validarPassword('')).toBe('La contraseña es obligatoria.')
    expect(validarPassword('corta')).toContain('al menos 8 caracteres')
  })

  it('rechaza una contraseña de más de 72 bytes (límite de bcrypt)', () => {
    expect(validarPassword('a'.repeat(73))).toContain('72 bytes')
    // Los caracteres acentuados ocupan más de un byte.
    expect(validarPassword('á'.repeat(37))).toContain('72 bytes')
  })
})

describe('validarLogin', () => {
  it('acepta credenciales completas', () => {
    expect(validarLogin({ email: 'usuario@ejemplo.com', password: 'cualquiera' })).toBeNull()
  })

  it('exige email y contraseña', () => {
    expect(validarLogin({ email: '', password: 'cualquiera' })).toBe('El email es obligatorio.')
    expect(validarLogin({ email: 'usuario@ejemplo.com', password: '' })).toBe('La contraseña es obligatoria.')
  })
})

describe('validarRegistro', () => {
  const CUENTA = {
    nombre: 'Benja',
    email: 'benja@ejemplo.com',
    password: '12345678',
    repetirPassword: '12345678'
  }

  it('acepta una cuenta válida', () => {
    expect(validarRegistro(CUENTA)).toBeNull()
  })

  it('exige un nombre de 2 a 50 caracteres', () => {
    expect(validarRegistro({ ...CUENTA, nombre: '  ' })).toBe('El nombre es obligatorio.')
    expect(validarRegistro({ ...CUENTA, nombre: 'B' })).toContain('entre 2 y 50 caracteres')
    expect(validarRegistro({ ...CUENTA, nombre: 'a'.repeat(51) })).toContain('entre 2 y 50 caracteres')
  })

  it('valida el email y la contraseña', () => {
    expect(validarRegistro({ ...CUENTA, email: 'mal' })).toBe('El email no tiene un formato válido.')
    expect(validarRegistro({ ...CUENTA, password: 'corta', repetirPassword: 'corta' })).toContain(
      'al menos 8 caracteres'
    )
  })

  it('exige que las dos contraseñas coincidan', () => {
    expect(validarRegistro({ ...CUENTA, repetirPassword: 'otra-distinta' })).toBe(
      'Las contraseñas no coinciden.'
    )
  })
})
