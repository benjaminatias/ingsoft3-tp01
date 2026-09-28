package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"gestor-peliculas/internal/auth"
	"gestor-peliculas/internal/handlers"
	"gestor-peliculas/internal/models"
)

// repoFalso es el doble del repositorio de usuarios: no toca ninguna base de datos.
// Es un mock porque además de devolver respuestas armadas, registra cómo lo llamaron
// para que el test pueda verificar esa interacción.
type repoFalso struct {
	existe          bool
	usuario         models.Usuario
	errBuscar       error
	errExiste       error
	llamadasCrear   int
	usuarioCreado   models.Usuario
	emailConsultado string
}

func (r *repoFalso) ExisteEmail(email string) (bool, error) {
	r.emailConsultado = email
	return r.existe, r.errExiste
}

func (r *repoFalso) Crear(u *models.Usuario) error {
	r.llamadasCrear++
	r.usuarioCreado = *u
	return nil
}

func (r *repoFalso) BuscarPorEmail(email string) (models.Usuario, error) {
	r.emailConsultado = email
	return r.usuario, r.errBuscar
}

func (r *repoFalso) BuscarPorID(id uint) (models.Usuario, error) {
	return r.usuario, r.errBuscar
}

// Arrange: el email ya está registrado. Act: se intenta registrar de nuevo.
// Assert: 409 y, sobre todo, que NUNCA se haya intentado crear la cuenta.
func TestRegistroConEmailExistenteDevuelve409YNoCrea(t *testing.T) {
	repo := &repoFalso{existe: true}
	router := handlers.NuevoRouterCon(handlers.NuevoConUsuarios(nil, repo))

	respuesta := ejecutar(t, router, http.MethodPost, "/api/auth/registro",
		`{"nombre":"Benja","email":"  Repetido@Ejemplo.COM ","password":"12345678"}`)

	if respuesta.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	if repo.emailConsultado != "repetido@ejemplo.com" {
		t.Fatalf("el handler debe consultar el email normalizado, consultó %q", repo.emailConsultado)
	}
	if repo.llamadasCrear != 0 {
		t.Fatalf("no debía crearse la cuenta, Crear se llamó %d veces", repo.llamadasCrear)
	}
}

// Registro exitoso: la contraseña llega al repositorio hasheada, nunca en texto plano.
func TestRegistroGuardaElHashYNoLaContrasena(t *testing.T) {
	repo := &repoFalso{}
	router := handlers.NuevoRouterCon(handlers.NuevoConUsuarios(nil, repo))

	respuesta := ejecutar(t, router, http.MethodPost, "/api/auth/registro",
		`{"nombre":"Benja","email":"nuevo@ejemplo.com","password":"contraseña-segura"}`)

	if respuesta.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	if repo.llamadasCrear != 1 {
		t.Fatalf("Crear debía llamarse una vez, se llamó %d", repo.llamadasCrear)
	}
	if repo.usuarioCreado.PasswordHash == "contraseña-segura" || repo.usuarioCreado.PasswordHash == "" {
		t.Fatal("el repositorio debe recibir el hash, no la contraseña en texto plano")
	}
	if !auth.VerificarPassword(repo.usuarioCreado.PasswordHash, "contraseña-segura") {
		t.Fatal("el hash guardado debe corresponder a la contraseña ingresada")
	}
}

// Caso de error: login de una cuenta inexistente responde 401 genérico, sin token.
func TestLoginConCuentaInexistenteDevuelve401(t *testing.T) {
	repo := &repoFalso{errBuscar: handlers.ErrUsuarioNoEncontrado}
	router := handlers.NuevoRouterCon(handlers.NuevoConUsuarios(nil, repo))

	respuesta := ejecutar(t, router, http.MethodPost, "/api/auth/login",
		`{"email":"nadie@ejemplo.com","password":"12345678"}`)

	if respuesta.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d", respuesta.Code)
	}
	var cuerpo map[string]any
	_ = json.Unmarshal(respuesta.Body.Bytes(), &cuerpo)
	if _, hayToken := cuerpo["token"]; hayToken {
		t.Fatal("una cuenta inexistente no debe recibir token")
	}
}

// Caso de error: si el repositorio falla, es un 500 y no un 401 (no es culpa del usuario).
func TestLoginConFalloDelRepositorioDevuelve500(t *testing.T) {
	repo := &repoFalso{errBuscar: errors.New("conexión perdida")}
	router := handlers.NuevoRouterCon(handlers.NuevoConUsuarios(nil, repo))

	respuesta := ejecutar(t, router, http.MethodPost, "/api/auth/login",
		`{"email":"alguien@ejemplo.com","password":"12345678"}`)

	if respuesta.Code != http.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d", respuesta.Code)
	}
}

// Login feliz: con el hash correcto en el doble se emite un token válido.
func TestLoginConCredencialesCorrectasEmiteToken(t *testing.T) {
	hash, err := auth.HashearPassword("contraseña-segura")
	if err != nil {
		t.Fatalf("no se pudo hashear: %v", err)
	}
	repo := &repoFalso{usuario: models.Usuario{ID: 7, Nombre: "Benja", Email: "b@ejemplo.com", PasswordHash: hash}}
	router := handlers.NuevoRouterCon(handlers.NuevoConUsuarios(nil, repo))

	respuesta := ejecutar(t, router, http.MethodPost, "/api/auth/login",
		`{"email":"b@ejemplo.com","password":"contraseña-segura"}`)

	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	var cuerpo struct{ Token string }
	_ = json.Unmarshal(respuesta.Body.Bytes(), &cuerpo)
	credenciales, err := auth.ValidarToken(cuerpo.Token)
	if err != nil || credenciales.UsuarioID != 7 {
		t.Fatalf("el token debe ser válido y del usuario 7: %v %+v", err, credenciales)
	}
}
