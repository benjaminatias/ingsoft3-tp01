# decisiones.md

## TP1 — Git colaborativo

### 1. Por qué Git no pudo resolver el conflicto solo

Git no puede resolver este conflicto solo ya que dos ramas escribieron la misma línea (existen dos
versiones), por lo que no puede decidir cuál es la "correcta"; entonces se debe intervenir para
resolver con criterio qué se debe hacer.

Para evitar este conflicto, dos ramas pueden modificar la misma línea pero no al mismo tiempo:
primero debe mergearse una rama para que luego la otra modifique esa misma línea sin generar ningún
conflicto.

### 2. Problemas encontrados y cómo los solucioné

Tuve confusiones con las ramas y algunos comandos de Git Bash, no recordaba exactamente su
funcionamiento. Para esto utilicé la IA para orientarme y entender correctamente cada paso.

Algunos inconvenientes surgieron al crear la carpeta `img` con las capturas, lo cual hice de forma
manual utilizando la consola — pero todo fue resuelto correctamente y entendido.

### 3. Declaración de uso de IA

Utilicé la IA para que me explique detalladamente algunas partes del paso a paso, para entender bien
lo que estaba haciendo. También la usé para generar comandos de Git de la consola, los cuales utilicé
y logré entender gracias a la IA.

Para verificar los resultados que me dio la IA, revisé y comparé constantemente si cada comando
ejecutado hacía correctamente lo que le pedí, chequeando siempre en el repositorio de GitHub los
cambios y detalles.

---

## TP2 — Contenedores

### 1. App elegida y por qué

Elegí un **gestor de películas** (backend en Go + Gin, frontend en React + Vite, base de datos
PostgreSQL con GORM), contra los criterios de la guía (§3.3):

- **Backend + Frontend + BD**: cumple los tres — API REST en Go, SPA en React, PostgreSQL relacional.
- **Corre local hoy, sin magia**: probado end-to-end con `docker compose up -d --build`, sin
  servicios externos ni configuración especial más allá del `.env`.
- **La entiendo / puedo modificarla**: entiendo bien el flujo general de la app y toda la parte de
  contenerización (Docker, compose, red, persistencia), que trabajé paso a paso probando cada
  comando. Las partes de autenticación (JWT, bcrypt) y algunos detalles internos del backend en Go
  todavía las estoy repasando — ver la declaración de uso de IA, sección 5.
- **Tamaño**: CRUD de películas + géneros + estadísticas + autenticación. Es un poco más grande que
  el "CRUD + 2-3 pantallas" que sugiere la guía como ideal, pero se mantiene manejable.

La idea concreta de "gestor de películas" surgió de una conversación con una IA, a la que le
consulté qué tipo de aplicación convenía armar para cumplir los requisitos de la materia. No partí
de una app propia preexistente.

### 2. Decisiones de contenerización

- **Backend (Go)**: Dockerfile multi-stage — etapa 1 con `golang:1.25-alpine` (compilador completo),
  etapa 2 con `alpine:3.20` (solo el binario compilado + certificados). Elegí Alpine para la imagen
  final por su tamaño mínimo. Con `CGO_ENABLED=0` genero un binario estático que no depende de
  librerías del sistema, lo cual es justamente lo que permite que la imagen final sea tan liviana.
- **Frontend (React/Vite)**: Dockerfile multi-stage — etapa 1 con `node:20-alpine` (build con
  `npm ci` + `npm run build`), etapa 2 con `nginx:1.27-alpine` sirviendo los estáticos. El
  `nginx.conf` resuelve el problema de la SPA (§2.6 de la guía): el frontend llama a rutas relativas
  (`/api/...`), y es nginx quien las reenvía al backend por nombre de servicio dentro de la red de
  compose — así evito CORS y no dejo ninguna URL de entorno escrita en el código del front.
- **Qué persiste y qué no**: solo los datos de PostgreSQL, en el volumen nombrado `db_data`, montado
  en `/var/lib/postgresql/data`. Todo lo demás (contenedores de backend y frontend) es efímero por
  diseño — se pueden recrear sin pérdida de información porque no guardan estado propio.
- **Secretos**: `POSTGRES_PASSWORD`, `DB_PASSWORD` y `JWT_SECRET` viven en un `.env` que no se
  commitea (está en `.gitignore`), con un `.env.example` commiteado como plantilla sin secretos
  reales.

### 3. Sobre la autenticación (JWT) — una decisión fuera del mínimo pedido

Mi proyecto incluye autenticación con JWT y contraseñas hasheadas con bcrypt, que no es un requisito
del TP2 (la consigna solo pide backend + frontend + BD). La agregué a propósito: se lo pedí
explícitamente a la IA porque quería que el proyecto tuviera login, más allá de lo mínimo pedido.

Entiendo que esto suma superficie a defender: sé explicar por qué las contraseñas se guardan
hasheadas con bcrypt y no en texto plano, qué es un JWT y para qué sirve, y por qué la función
`Conectar()` de `database.go` reintenta la conexión varias veces en vez de fallar directo (por la
diferencia entre "el contenedor arrancó" y "el servicio está listo", que también se ve con
`depends_on` + `healthcheck` en el compose).

### 4. Problemas encontrados y cómo los resolví

- **`docker push` rechazado con `denied`**: la sesión de `docker login` se había invalidado. Se
  resolvió repitiendo el login contra `ghcr.io` antes de reintentar el push.
- **`docker pull` con `unauthorized` al verificar la publicación**: los packages de ghcr.io nacen
  **privados** por defecto. Hasta no cambiar la visibilidad a *Public* desde la web de GitHub
  (Package settings → Change visibility), ningún `pull` sin sesión iniciada funciona, aunque el
  `push` haya sido exitoso.
- **Backend en contenedor no se conectaba a la base**: usar `localhost` en la connection string
  apunta al propio contenedor, no a la máquina host. Se resolvió reemplazándolo por el nombre del
  servicio (`db`) dentro de la red de compose.
- **`docker images` no mostraba las imágenes base usadas en el build**: BuildKit las usa
  internamente sin registrarlas siempre con su tag completo en el listado visible. Se resolvió
  haciendo un `docker pull` explícito de esas imágenes antes de listarlas.

### 5. Declaración de uso de IA

Usé IA (Claude) de forma intensiva durante el desarrollo de este proyecto:

- La **idea de la aplicación** (gestor de películas) surgió de consultarle a la IA qué tipo de app
  convenía para cumplir los requisitos de la materia.
- **La mayor parte del código** (backend en Go, frontend en React, Dockerfiles, `docker-compose.yml`,
  `nginx.conf`) fue generada con asistencia de IA; yo revisé y aprendí.
- **Verificación**: probé todo manualmente — registro y login desde el navegador y también por
  `curl`/`Invoke-RestMethod` (incluyendo el caso de error, con credenciales inválidas después de un
  `docker compose down -v`), creación y listado de películas autenticado con el token JWT, y la
  prueba de persistencia completa (`down` conserva datos, `down -v` los borra). No corrí todavía la
  suite de tests automatizados (`go test`, Vitest) que trae el proyecto — queda pendiente para cuando
  lo repase antes de la defensa.
- **Estado actual de comprensión**: entiendo bien el flujo general (Docker multi-stage, compose,
  persistencia, red de servicios) porque lo trabajé paso a paso. Las partes de autenticación (JWT,
  bcrypt) y algunos detalles del backend en Go todavía los estoy repasando para poder explicarlos con
  seguridad en la defensa oral.
