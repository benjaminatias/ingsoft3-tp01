# decisiones.md

## Enlaces del TP7

**Paquetes públicos** (los mismos del TP6, ahora con consumidor: los cuatro servicios de Render ejecutan estas imágenes; se bajan con `docker pull` sin credenciales y cada imagen lleva una sola etiqueta, la de su commit):
- Backend: https://github.com/users/benjaminatias/packages/container/package/ingsoft3-tp01-backend
- Frontend: https://github.com/users/benjaminatias/packages/container/package/ingsoft3-tp01-frontend
- Etiqueta de la imagen que rompió la app (corrida roja): `sha-431e81a8444a68b6b0075ca83ea9552f22fe32ce`
- Etiqueta de la imagen del arreglo, la que llegó a PROD (corrida verde): `sha-da5e68c3a68c9b343a18b8a8f138851da9c5ddf3`

**La corrida donde la integración quedó VERDE y la e2e ROJA** (frenó la promoción; `deploy-prod` ni arrancó):
- Commit que rompió la app (cambia el texto del botón de alta en el front, sin tocar `e2e/`): https://github.com/benjaminatias/ingsoft3-tp01/commit/431e81a8444a68b6b0075ca83ea9552f22fe32ce
- Corrida: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39332337445 (la e2e en rojo: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39332337445/job/118334101292)
- Reporte de integración (verde): https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39332337445/artifacts/12224504300
- Reporte e2e (rojo, con captura y traza): https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39332337445/artifacts/12225237058

**La corrida completa en verde hasta PROD**, después del arreglo (`build → deploy-qa → integracion → e2e → aprobación → deploy-prod`):
- Corrida: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39334390495 (el deploy a PROD: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39334390495/job/118341046850)
- Commit del arreglo: https://github.com/benjaminatias/ingsoft3-tp01/commit/da5e68c3a68c9b343a18b8a8f138851da9c5ddf3

**Entornos vivos hasta la defensa** (los *Events* de cada servicio de Render no se pueden linkear: los muestro en vivo):
- QA: https://peliculas-front-qa.onrender.com (api: https://peliculas-api-qa.onrender.com)
- PROD: https://peliculas-front-prod.onrender.com (api: https://peliculas-api-prod.onrender.com)

## Enlaces de este TP (TP6)

**Paquetes públicos** (se bajan con `docker pull` sin credenciales; los dos etiquetados con el commit del merge `17fa23dc4b74368d2f5c3a97ee8e108d767618f8`):
- Backend: https://github.com/users/benjaminatias/packages/container/package/ingsoft3-tp01-backend
- Frontend: https://github.com/users/benjaminatias/packages/container/package/ingsoft3-tp01-frontend

**La cadena que prueba que el registry es confiable** (corridas del PR «Publicacion de imagen»):
- Job del PR con los tests en verde, donde «Entrar al registry» aparece salteado: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39061871417/job/117481315673
- Lista de pasos de la corrida de `main`, donde «Construir y publicar la imagen del backend» es el último paso: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39062370300/job/117482789211

**Entornos vivos hasta la defensa:**
- QA: https://peliculas-front-qa.onrender.com (api: https://peliculas-api-qa.onrender.com)
- PROD: https://peliculas-front-prod.onrender.com (api: https://peliculas-api-prod.onrender.com)

Nota: la app tiene dos Dockerfiles (backend y frontend), así que se publican **dos paquetes**.

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

## TP3 — Planificación y trazabilidad

### 1. Duración del sprint

Elegí un sprint de **2 semanas**. Es una duración estándar en la industria: da margen suficiente
para completar una historia con sus tareas sin la presión de un sprint de 1 semana (que en un
proyecto individual, con tiempo limitado entre clases y otras materias, sería demasiado ajustado),
pero sigue siendo lo bastante corto como para poder ajustar el rumbo rápido si algo no está
funcionando, en vez de esperar un mes entero para replanificar.

### 2. Límite de trabajo en progreso

Elegí un límite de **2** en la columna *In Progress*. Trabajando individualmente, la guía sugiere
"cantidad de personas + 1" — en mi caso, 1 + 1 = 2. El "+1" me deja una válvula para cuando algo
queda esperando (por ejemplo, una tarea bloqueada por una duda o una revisión) y necesito poder
avanzar en otra cosa sin quedarme frenado. Pasarme de ese número haría que el límite deje de
cumplir su función: evitar que tenga muchas cosas a medio hacer al mismo tiempo, en vez de terminar
una por una.

### 3. Diagnóstico de la historia mal escrita

La historia de ejemplo: *"Como desarrollador quiero crear la tabla usuarios para guardar los
datos."*

**Por qué está mal escrita:** es una tarea técnica disfrazada de historia de usuario. Nadie —ni un
cliente ni un usuario final— "quiere" que exista una tabla en la base de datos; eso es un medio,
no un fin que le aporte valor a alguien. El "para guardar los datos" no es un beneficio real, es
casi una repetición circular del "qué". Además, no tiene criterios de aceptación, así que no hay
forma de verificar cuándo está "hecha".

**Cómo la reescribiría:** separando el beneficio real del paso técnico. Por ejemplo: *"Como usuario
quiero registrarme con mi email y contraseña para poder guardar mis películas favoritas de forma
personal"* como historia, y "crear la tabla `usuarios`" pasa a ser una de sus tareas técnicas, no
la historia en sí.

### 4. Problemas encontrados y cómo los resolví

- **`gh auth status` mostró que no estaba logueado con la GitHub CLI**, a pesar de que ya venía
  usando Git normalmente en los TPs anteriores. Se resolvió con `gh auth login`, autenticando por
  navegador con el código de 8 caracteres que la CLI generó — la autenticación de `gh` es
  independiente de la de Git/GitHub en el navegador, algo que no tenía claro antes de este TP.

### 5. Declaración de uso de IA

Usé IA (Claude) para guiarme paso a paso durante todo el TP3: crear el Project con la
configuración correcta (auto-import, visibilidad pública), las etiquetas y los issues por
terminal, vincular la jerarquía con sub-issues, configurar el sprint y el límite de WIP, y armar
el Pull Request con `Closes #16` para probar la trazabilidad. Entendí el proceso a medida que lo
hacía —no fue copiar y pegar sin más—: verifiqué cada paso en la web de GitHub antes de seguir al
siguiente (por ejemplo, confirmando visualmente que la jerarquía quedó bien anidada, que el issue
se cerró solo al mergear el PR, y que el tablero se movió sin intervención manual).

## TP4 — CI: Pipelines as Code

### 1. Estructura elegida del pipeline

El workflow tiene **dos jobs en paralelo**, `build-backend` y `build-frontend`, en vez de uno solo
o de pasos secuenciales dentro del mismo job. La razón es que mi app tiene **dos Dockerfiles
independientes** (backend en Go, frontend en React/Vite), y construir uno no depende del resultado
del otro: no tiene sentido esperar a que termine el backend para recién empezar el frontend. Al
separarlos en jobs, cada uno corre en su propio runner limpio y en simultáneo, así que el tiempo
total del pipeline es el del job más lento, no la suma de los dos.

La contrapartida de esta separación es que los jobs **no comparten filesystem**: si alguna vez
necesitara pasarle algo de un job a otro (por ejemplo, un artefacto), tendría que declararlo
explícitamente. Para este TP no hace falta, porque cada job solo necesita su propio código fuente
(que trae con su propio `actions/checkout`) y no depende de nada que produzca el otro.

### 2. Qué cachea el pipeline

Se cachean las **capas de Docker** de cada imagen, vía `docker/setup-buildx-action` +
`cache-from`/`cache-to: type=gha`, con un `scope` distinto por job (`backend` y `frontend`). Sin
ese `scope` separado, los dos jobs comparten el mismo estante de cache por default y se pisan entre
sí — lo comprobé leyendo la documentación de Docker antes de escribirlo, no me pasó en la práctica
porque lo puse bien desde el principio.

Lo que efectivamente se reutiliza son las capas que no cambiaron entre corridas: en el backend, por
ejemplo, `go mod download` se cachea completo mientras no toque `go.mod`/`go.sum`, y solo se
rehacen las capas de `COPY` del código y la compilación. Verifiqué esto en la pestaña Actions: la
corrida que subió el cache por primera vez tardó 1m41s; la siguiente, que lo reutilizó, tardó 29s,
con la palabra `CACHED` en las capas de dependencias de ambos jobs.

**Qué pasa si el cache desaparece:** nada grave, solo se pierde velocidad. GitHub puede desalojarlo
en cualquier momento (tiene límite de tamaño y política de expiración), así que el pipeline tiene
que poder reconstruir todo desde cero sin el cache — que es justamente lo que pasó en la primera
corrida, antes de que existiera cualquier capa guardada. Si el pipeline fallara sin cache, no sería
un cache: sería una dependencia escondida que se estaba colando sin que yo la notara.

### 3. Por qué el pipeline construye con mi Dockerfile en vez de compilar por su cuenta

Si el workflow compilara por su lado (por ejemplo, corriendo `go build` y `npm run build`
directamente en el runner, sin pasar por Docker), tendría **dos definiciones de build distintas**:
la que usa el pipeline para verificar, y la que uso yo para efectivamente correr la app en
contenedores (TP2). Esas dos definiciones divergen tarde o temprano — una diferencia de versión de
Go, una variable de entorno que solo está seteada en un lado, una dependencia del sistema que el
Dockerfile instala y el runner no tiene — y en ese momento estaría verificando algo que no es lo
que después despliego. Usar el mismo Dockerfile como única fuente de verdad evita ese problema:
lo que el pipeline confirma que compila es exactamente lo que se empaqueta y se corre.

### 4. Problemas encontrados y cómo los resolví

- **`docker build` fallaba con un error de conexión al motor** (`open //./pipe/dockerDesktopLinuxEngine`).
  No era un problema del Dockerfile ni del comando: Docker Desktop no estaba levantado en ese
  momento. Se resolvió simplemente abriéndolo y esperando a que el motor terminara de arrancar.

- **Al romper el build a propósito agregando un import inexistente en Go, VS Code me lo borraba
  solo al guardar.** La extensión de Go tiene activado el organizador de imports (`goimports`) al
  guardar, que detecta imports no usados y los elimina automáticamente — y como el paquete falso no
  se usaba en ningún lado del código, lo sacaba antes de que llegara a compilarse. Lo resolví
  escribiéndolo como **import en blanco** (`_ "paquete/que/no/existe"`), que Go trata como
  intencional (un import por efectos secundarios) y que `goimports` no toca aunque no se referencie
  ningún identificador del paquete.

- **El badge del README quedó mal la primera vez**: pegué por error el nombre de mi rama
  (`docs/badge-readme`) en vez del Markdown del badge que copié de GitHub. Lo detecté revisando el
  archivo antes de hacer commit, y lo corregí pegando la línea correcta
  (`[![CI](...badge.svg)](...)`), con los dos niveles de corchetes: uno para la imagen y otro para
  el link de destino, así al clickear el badge lleva al historial de corridas y no a un SVG suelto.

- **Los tags `v2.0.0` y `v3.0.0` no se habían creado en su momento**, solo `v1.0.0`. Los agregué
  retroactivamente sobre los commits donde efectivamente cerré cada práctico (identificados
  revisando el historial de commits y el contenido de cada uno), sin tocar ningún tag existente —
  un mismo commit puede tener más de un tag, así que no hizo falta mover ni rehacer nada.

### 5. Declaración de uso de IA

Usé IA (Claude) para guiarme paso a paso durante todo el TP4: entender la teoría de CI antes de
tocar el YAML, escribir el workflow con los dos jobs y el cache de capas, configurar los required
status checks en Settings → Branches, y diseñar la forma de romper el build a propósito (elegí
un import inexistente en Go, adaptado a mi stack, siguiendo el criterio de la guía según el tipo de
lenguaje). Verifiqué cada paso contra la evidencia real antes de seguir: confirmé el error exacto
en la terminal antes de asumir que el build estaba roto como quería, leí los logs de Actions
buscando la palabra `CACHED` en vez de asumir que el cache funcionaba, y confirmé en la propia
página del PR que los checks aparecían como *Required* y que el botón de merge quedaba bloqueado
antes de dar el paso por cumplido.

## TP5 — Calidad automatizada: tests, coverage y el umbral que frena un merge

### 1. Qué lógica elegí testear y por qué esa

Elegí testear donde un bug me dolería más en mi app de películas. Primero la **validación de
películas**: la regla más delicada es la puntuación según el estado, porque una película pendiente
no puede tener puntuación y una vista tiene que tenerla entre 1 y 10 con un solo decimal. Si esa
regla se rompe, la base se llena de datos inconsistentes y después no hay forma fácil de
arreglarlos. Segundo, las **cuentas y la autenticación** (bcrypt, JWT, email repetido), donde un
bug es un problema de seguridad y no de estética. Tercero, los **géneros**: que no se pueda borrar
uno que tiene películas asociadas y que no haya dos con el mismo nombre aunque cambien las
mayúsculas. Y en el frontend, el **cliente de la API** (`api.js`), porque es el único lugar que
decide qué hacer cuando el token vence: un 401 con sesión abierta cierra la sesión, pero un 401 en
el login no.

### 2. Las técnicas de test y las herramientas que usé

Mi stack es Go + vitest, no .NET, así que averigüé qué usar para cada cosa que pide la consigna:

- **Test parametrizado:** en Go, *table-driven tests* con `t.Run` (por ejemplo
  `TestRegistroYLoginInvalidos` y `TestIntegracionGenerosErrores`, donde cada fila es una regla
  distinta). En el frontend, `it.each` sobre los status 400, 404, 409 y 500 en `api.test.js`.
- **Caso de error:** en el backend, un login de una cuenta inexistente devuelve 401 sin token, y si
  el repositorio falla devuelve 500 y no 401, porque no es culpa del usuario. En el frontend, la
  falla de conexión y el 401 con y sin token.
- **Que la dependencia entre desde afuera:** una interfaz (`RepositorioUsuarios`) que el `Handler`
  recibe como campo, explicada en el punto 3.
- **Fabricar el doble:** un struct hecho a mano (`repoFalso`) que implementa la interfaz, sin
  librerías de mocks. En el frontend, `vi.fn()` en lugar de `fetch`.
- **Medir la cobertura:** `go test -coverprofile` con `-coverpkg=./internal/...` en el backend, y
  `vitest run --coverage` (motor v8) en el frontend.
- **Un umbral que rompe el build:** en el frontend, `coverage.thresholds` de vitest. En Go no
  existe una bandera para eso, así que lo resolví con un script (`backend/scripts/cobertura.sh`)
  que lee el total de `go tool cover` y termina con error si está por debajo del número.
- **Qué entra en la cuenta:** el `-coverpkg` en Go y el `include` de vitest (punto 4).
- **Que las herramientas entren a la etapa de tests del Dockerfile:** la etapa `test` copia
  `tests/` y el script (en Go no hay dependencias de desarrollo aparte), y en el frontend le saqué
  `tests` al `.dockerignore`, que si no los dejaba afuera.

Como sin base de datos los tests de integración se omiten y la cobertura queda muy baja (46,9 %),
el job `build-backend` levanta un PostgreSQL descartable como servicio y los tests corren con
`--network host`. Todo esto entró en el PR del pipeline:
https://github.com/benjaminatias/ingsoft3-tp01/pull/28

### 3. El refactor para poder mockear

Antes, los handlers de `Registro`, `Login` y `Perfil` hablaban directo con `*gorm.DB`, o sea que
**no podía ejecutar un handler de cuentas sin tener PostgreSQL corriendo**: no había dónde meter un
doble. Lo que cambié fue agregar la interfaz `RepositorioUsuarios` (`ExisteEmail`, `Crear`,
`BuscarPorEmail`, `BuscarPorID`) con una implementación real sobre GORM, y hacer que el `Handler`
la reciba como campo. `Nuevo(db)` la arma con GORM, así que el resto de la app no cambió. Además
el repositorio devuelve un error propio (`ErrUsuarioNoEncontrado`) para que los handlers no
dependan del error interno de GORM.

Mi test `TestRegistroConEmailExistenteDevuelve409YNoCrea` **reemplaza la dependencia con la base de
datos** por `repoFalso`, y el assert verifica dos cosas: que la respuesta sea 409 y que `Crear`
nunca se haya llamado. Eso último es lo que un mock permite y un stub no: el stub solo devuelve
respuestas armadas, el mock además registra cómo lo llamaron para que el test lo compruebe. Si
invierto la regla del `if` en `Registro`, ese test se pone en rojo.

### 4. Mi umbral de coverage y qué dejé afuera de la cuenta

En el **backend** puse **70 % de sentencias**, y al cerrar el TP5 daba **77,3 %**. Go no tiene métrica de ramas,
solo sentencias, así que no puedo reportar el número de rama de ese lado. En el **frontend** puse
**80 % de líneas y 85 % de ramas**: medí 82,4 % y 87,0 %, y dejé unos 2 puntos de margen, chico a
propósito para frenar una caída real sin fallar por ruido. Puse umbral en las dos métricas porque
con vitest 2.1.9 una función que nadie llama suma líneas sin cubrir pero casi no mueve las ramas
(en mi demostración las ramas quedaron en 86,4 %, arriba del umbral, y frenó por líneas). Con los tests de
`recomendaciones.js` el frontend quedó en 84,7 % de líneas y 89,3 % de ramas (y con los tests que sumé
en los TPs siguientes hoy mide 85,2 % y 89,6 %). En mi Windows el mismo comando da 85,1 % de líneas,
porque v8 cuenta distinto según el sistema: el número que manda es el del pipeline, que corre en Linux.

Con el 70 % me pasó algo que vale la pena contar: lo elegí mirando mi medición local (74,4 %),
pero la primera corrida limpia en GitHub dio **69,9 %** y el paso se puso rojo por una décima.
Podía bajar el umbral a 65, pero eso era ajustar el número al código. Preferí escribir los tests
que faltaban: el punto más flojo era `generos.go` (33 %), así que agregué tests de integración de
géneros (ciclo completo, errores, duplicados y que no se borre uno con películas). Subió a 78,9 %
y el total quedó con más de 7 puntos de margen. Para subir el umbral haría falta testear
`peliculas.go` (66,7 %) y `estadisticas.go` (55,6 %) en el backend, y en el frontend las
funciones CRUD de `api.js` (70,6 % de líneas y solo 36 % de funciones). Si mañana lo subo 10
puntos, el backend fallaría hoy.

Sobre línea contra rama: la de línea puede mentir más, porque un `if` en una sola línea cuenta
como cubierto aunque solo se recorra uno de sus dos caminos. La de rama exige recorrer los dos.

**Qué dejé afuera de la cuenta.** En el backend, `cmd/api` (el arranque, `main.go`), que queda
fuera por estar fuera de `internal/`, y los modelos, que son structs sin sentencias ejecutables.
`database.go` sí cuenta (73,3 %) porque lo ejercitan los tests de integración. En el frontend mido
solamente `src/utils` y `src/api`: los componentes React, `App.jsx` y `main.jsx` quedan afuera
porque son UI y necesitan tests con DOM, que este TP no pide. Lo armé como lista de lo que
**incluyo**, así un archivo nuevo en esas carpetas entra a la cuenta solo aunque nadie lo pruebe
(es justo lo que usé en la demostración).

Pruebas: la corrida verde de `main` con el resumen de cobertura del backend (77,3 %) y el reporte
descargable:
https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/36480097250
Y la primera corrida roja del backend por umbral (69,9 % contra 70 %):
https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/36478044700

### 5. Por qué coverage alto no garantiza calidad

La cobertura mide qué líneas **se ejecutaron**, no si alguien **verificó** algo. Con mi propio
código: si escribo un test que llame a `nivelDePuntuacion(9)` sin ningún `expect`, esa línea sube
la cobertura igual y no comprueba nada. Un 77 % no quiere decir que el 77 % de mi código está
verificado, quiere decir que esas sentencias corrieron durante los tests. Para comprobar que mis
tests sí verifican, se hicieron pruebas de mutación a mano (con ayuda de Claude, en una copia del
código): se invirtió una regla, por ejemplo `if existe` por `if !existe` en el registro o `>= 9`
por `> 9`, y en ambos casos algún test se puso en rojo.

### 6. El Pull Request bloqueado

Agregué `frontend/src/utils/recomendaciones.js` (tres funciones con varios caminos) **sin ningún
test**. Compilaba perfecto y los 54 tests seguían en verde, pero el check `build-frontend` se puso
rojo en la métrica **líneas**: `ERROR: Coverage for lines (71.67%) does not meet global threshold
(80%)`. Las 54 sentencias nuevas sin cubrir bajaron el total del 82,4 % al 71,7 %. Para
arreglarlo escribí `recomendaciones.test.js`, con un test por cada camino que las funciones
declaran (27 casos entre `it.each` y tests sueltos), el check pasó a verde (84,7 %) y mergeé.

Este freno es distinto del del TP4: aquel frena cuando el código **no compila**, este frena cuando
compila y pasa todo, por un criterio de calidad que fijé yo. Y lo que deja pasar igual son los
tests que ejecutan el código sin verificarlo, los errores de lógica que ningún test contempla y
todo lo de la interfaz, porque los componentes no entran a la cuenta.

Hice **dos Pull Requests a propósito**: el primero cuenta la historia (rojo, tests, verde, merge)
y el segundo queda **abierto y en rojo hasta la defensa**, con el mismo problema sin arreglar,
para poder comprobar el freno sin depender de una captura.
- PR 1, el mergeado, con la secuencia completa: https://github.com/benjaminatias/ingsoft3-tp01/pull/29
- Corrida roja por umbral de ese PR, con el número en el log:
  https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/36481374573
- PR 2, el freno vigente (abierto y en rojo): https://github.com/benjaminatias/ingsoft3-tp01/pull/30

### 7. El ejercicio del camino sin cubrir

- **Qué línea es:** `frontend/src/utils/formato.js`, el `return '-'` de `formatearPuntuacion`
  cuando `Number.isNaN(numero)` (línea 19). Algo parecido pasa en `formatearDecimal`, donde el
  `return '-'` de la línea 6 tampoco lo recorre ningún test.
- **Qué entrada la recorrería:** un valor concreto, `'abc'`: `formatearPuntuacion('abc')` entraría a
  ese camino y devolvería `'-'`.
- **Qué decidí hacer: no lo agregué.** Revisé el código y la aplicación solo llama a
  `formatearPuntuacion` con `pelicula.puntuacion` (en `PeliculaItem.jsx`), que la API siempre manda
  como número o `null`, y el `null` y el vacío ya se resuelven en el `if` de arriba. Ese camino es
  defensa extra frente a un dato que el backend no manda, por eso lo dejé sin test.

### 8. Problemas encontrados y cómo los resolví

- **El umbral del backend falló en la primera corrida (69,9 % contra 70 %)** aunque en mi máquina
  daba 74,4 %, porque la corrida limpia de GitHub mide distinto. Lo resolví agregando tests de
  géneros (punto 4) en vez de bajar el número.

- **Conflicto en `frontend/package.json`** al abrir el PR del pipeline: mi rama había salido de un
  `main` que todavía no tenía el PR anterior del frontend, y los dos tocaban las mismas líneas de
  los scripts. Lo resolví en el editor web de GitHub dejando la versión que tenía las dos cosas.

- **El `git push` fue rechazado** después de resolver ese conflicto desde la web, porque GitHub
  había agregado un commit de merge que yo no tenía. Lo resolví con `git pull` y volví a pushear.

- **Saltos de línea en Windows:** Git convierte LF a CRLF y eso puede romper un script `.sh` dentro
  del contenedor Linux. Agregué un `.gitattributes` que fuerza LF en los `.sh`.

### 9. Declaración de uso de IA

Usé IA (Claude) para todo el TP5: el refactor de `RepositorioUsuarios`, los tests con mock del
backend, los tests de integración de géneros, los tests del frontend (`api.test.js` y
`recomendaciones.test.js`), la etapa de tests de los dos Dockerfiles, el script `cobertura.sh` con
el umbral, los cambios del `ci.yml` y el código sin tests de la demostración del freno. Lo verifiqué
así: corrí `go test ./...` y `npm run test:coverage` en mi máquina, revisé las corridas en Actions
para ver el resumen de cobertura y los artefactos, y comprobé en los Pull Requests que los checks
aparecieran como *Required* con el merge bloqueado cuando correspondía. Además, Claude comprobó
que los tests detectan cambios invirtiendo reglas en una copia del código (un test se puso en rojo
en cada caso). El caso que **no** está cubierto es el del punto 7.


## TP6 — CD: environments, aprobaciones y deployment patterns

### 1. Por qué el artefacto se publica solo con la verificación en verde

La imagen se publica en `ghcr.io`, etiquetada con el commit que la produjo (`sha-<commit>`), y solo cuando el cambio entra a `main`. Lo que garantiza que lo publicado es lo que pasó la verificación son **tres eslabones encadenados**, y ninguno es un control aparte que yo tenga que acordarme de mirar: nada entra a `main` sin los checks en verde (son *required*), solo `main` publica (el `push` del paso de publicar depende de `github.ref == 'refs/heads/main'`, y el login al registry solo corre en `push`), y publicar es **el último paso** del job que testea, sin `always()` ni `!cancelled()`, así que si un paso anterior falla (tests rojos, cobertura bajo el umbral) el job se corta antes de llegar. Los dos primeros eslabones se ven en la corrida de un PR, donde «Entrar al registry» aparece salteado aunque todo pasó. El tercero se ve en la lista de pasos de una corrida de `main` (links arriba).

Si se publicara igual, el registry dejaría de significar «esto pasó la verificación» y pasaría a ser un depósito de lo que alguien construyó: tendría que mirar cada imagen y confiar en que alguien probó. Lo que la cadena **no** garantiza: un `docker push` hecho a mano con mis credenciales entra igual, y un tag como `sha-<commit>` es mutable (se puede volver a empujar con otro contenido). Lo que lo cerraría es desplegar por **digest** (`@sha256:...`), que identifica el contenido y no el nombre.

### 2. Continuous Delivery o Continuous Deployment

Implementé **Continuous Delivery**: QA se despliega solo con cada merge a `main`, pero PROD requiere una aprobación humana. Para Continuous Deployment (PROD también automático) me faltaría poder confiar en que mi pipeline detecta un cambio malo sin que yo lo mire: tests de interfaz, observabilidad (tasa de errores, métricas) y un rollback automático. Hoy mi propio caso lo muestra: el cambio con una errata en el título pasó tests y smoke en verde y solo lo detectó una persona mirando QA. Para este contexto (un equipo de una persona, una app chica, sin usuarios reales) la aprobación humana cuesta un minuto y compra una última mirada; en un equipo con muchos deploys por día y buena observabilidad, esa aprobación sería un cuello de botella y no agregaría valor.

### 3. El diseño de la cadena y el alcance de cada secret

La cadena completa es: PR verde → merge a `main` → `build-backend` y `build-frontend` en verde → `deploy-qa` → smoke de QA verde → aprobación humana → `deploy-prod` → smoke de PROD.

- `deploy-qa` tiene `needs: [build-backend, build-frontend]` (sin CI verde no hay deploy) y un `if` que lo limita a `push` sobre `main`; los PRs verifican pero nunca despliegan. Su environment `qa` no tiene reglas: es automático a propósito.
- `deploy-prod` tiene `needs: deploy-qa` (PROD solo después de QA vivo) y el environment `production`, que tiene como *required reviewer* a mí mismo, con *Prevent self-review* desactivado (soy el único del equipo). El workflow **se pausa** ahí hasta que alguien aprueba. No repite el `if`: lo hereda por `needs`. Tiene un `concurrency` con `cancel-in-progress: false`, para que dos deploys a PROD no se pisen.
- Los jobs de deploy disparan los Deploy Hooks de Render con **`&ref=$GITHUB_SHA`**. Sin ese parámetro, un hook despliega la punta de la rama y no el commit que el pipeline acaba de verificar: con dos merges seguidos, la corrida del primero desplegaría el segundo, y entre que una corrida queda esperando aprobación y que alguien aprueba, `main` se puede mover. Con `&ref`, se aprueba un commit y se despliega ese commit.
- **Secrets:** los de QA (`RENDER_HOOK_API_QA`, `RENDER_HOOK_FRONT_QA`) viven en el environment `qa`, y los de PROD (`RENDER_HOOK_API_PROD`, `RENDER_HOOK_FRONT_PROD`) viven **dentro del environment `production`**, no en *Repository secrets*. Eso importa porque un secret de environment solo se entrega a un job que **pasó la regla de ese environment**: si los de PROD estuvieran en el repo, cualquier job de cualquier workflow (por ejemplo uno de un PR) los podría leer sin que nadie aprobara nada, y el gate no protegería nada. Para entender qué pasa entre el merge y PROD: el merge dispara el CI, el CI en verde habilita `deploy-qa`, QA verde habilita `deploy-prod`, y ahí GitHub frena hasta mi aprobación; recién al aprobar se le entregan los secrets de PROD al job.

### 4. Qué mira el aprobador antes de aprobar

Antes de apretar *Approve* miro: que `deploy-qa` esté en verde y que el smoke haya pasado, **abrir QA y ver el cambio con mis ojos** (no solo el título: el detalle que distingue la versión nueva de la anterior), y que el commit de la corrida sea el que espero. Lo que **no** puedo ver: los tests y el smoke no miran textos de pantalla, ni qué versión está corriendo en realidad (ver punto 7), ni el estado de los datos de PROD, ni qué errores ve un usuario real. Me falta observabilidad para decidir mejor: un `/health` que informe el commit, una tasa de errores y logs centralizados.

Mi caso real: aprobé **sin querer** el cambio con la errata (apreté *Approve* en vez de *Reject*), así que el error llegó a PROD. Eso muestra el límite de un gate humano: la decisión depende de un click y no tiene una confirmación. Lo resolví con el rollback del punto 8.

**El rechazo documentado:** rechacé después ese mismo cambio (el PR #35) en el gate de `production`, con el motivo que escribí al ver QA con el título mal escrito: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39099242979/job/117596336792

**El flujo completo** (QA verde, *Waiting for review*, aprobación, PROD desplegado y smoke verde) es la corrida del arreglo (PR #36): https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/39101448910

### 5. La letra chica del free tier y cómo la maneja el pipeline

En el plan gratis de Render, una instancia **se duerme por inactividad** y el primer pedido puede tardar 50 segundos o más (arranque en frío), y el build tarda lo que tarde la caché. Neon también suspende su cómputo cuando no se usa. Para esto, el smoke test **reintenta**: 30 intentos cada 20 segundos (hasta unos 10 minutos), en lugar de un `curl` seco que daría falsos rojos. En la práctica, mis rollbacks y deploys con caché tardaron menos de un minuto, pero el primer registro de usuario en un entorno dormido falló con el mensaje genérico del front hasta que la api despertó.

### 6. Qué garantía pierdo porque Render reconstruye

Render **no usa la imagen que publiqué en ghcr**: al recibir el hook, clona el repo en ese commit y construye su propia imagen. La imagen que corre en QA y PROD sale del mismo commit y del mismo Dockerfile que verifiqué, pero es **otro build**: una dependencia o una imagen base que cambie entre un build y otro podría producir un resultado distinto sin que el CI lo note. Es decir, la cadena de confianza del registry se corta en el último tramo. Lo cerraría desplegando en Render la imagen ya publicada por digest, de modo que lo que corre sea byte a byte lo que se verificó.

### 7. Qué prueba mi smoke test y qué no

El script `scripts/smoke.sh` verifica cuatro cosas: que `/health` de la api responda (el endpoint hace ping a PostgreSQL), que un login con credenciales falsas devuelva **401** (prueba que la app consulta la tabla de usuarios; si la base estuviera inalcanzable daría 500), que el front sirva la página, y que `/health` **a través del front** responda (prueba que nginx le llega a la api de su entorno). Falla el job si no responde en el tiempo de reintentos. Lo probé con un servidor de prueba: pasa cuando todo contesta, y falla si el login da algo distinto de 401 o si nadie escucha.

Lo que **no** dice es qué versión corre: si el hook se disparó pero el build nuevo todavía no terminó, el entorno viejo sigue contestando verde y el smoke da verde contra la versión anterior. Mi `deploy-qa` terminó en 36 segundos, lo que es una pista de ese riesgo. Lo resolvería haciendo que `/health` devuelva el commit y que el smoke lo compare con `$GITHUB_SHA`.

### 8. Deployment pattern para una producción real y plan de rollback

**Patrón elegido: blue-green.** Para esta app, con usuarios reales, usaría blue-green: dos entornos completos (el actual y el nuevo) y un router que cambia de uno al otro. Lo elijo porque el rollback es instantáneo (volver el switch) y porque mi app es chica y los contenedores son baratos, así que el costo de tener el doble de infraestructura es bajo. El riesgo es la **base de datos compartida**: las dos versiones tendrían que ser compatibles con el mismo esquema, y con `AutoMigrate` de GORM una migración que borra una columna rompería a la versión vieja. Haría falta hacer los cambios de esquema en dos pasos (agregar, después quitar). Canary sería mejor para ir validando con una parte del tráfico, pero hoy me faltan las métricas para decidir si el canary está sano (tasa de errores, latencia) y un router con tráfico ponderado; sin eso, un canary sería a ciegas.

**Mi plan de rollback actual**, el que ejecuté de verdad cuando aprobé la errata, y su número medido:
1. Identificar el último commit bueno, en la lista *Deploys* del servicio de PROD en Render (acá `d7629e6`).
2. Cargar los Deploy Hooks de PROD en la terminal sin dejarlos a la vista (`read -rs`) y disparar los dos con `&ref=<commit bueno>`.
3. Verificar en *Deploys* que el deploy de ese commit figure como *Live* en los dos servicios, y recargar PROD.

**Tiempo medido:** disparé el rollback a las 02:21:27; Render empezó los dos deploys a las 02:21:28; la api quedó *Live* en 17,8 s (≈ 02:21:46) y el front a las 02:21:56. El rollback tardó **unos 29 segundos**, medido en el servicio más lento. Es poco porque el build estaba cacheado (sin caché serían varios minutos) y no incluye el arranque en frío del plan gratis. La lista de *Deploys* de `peliculas-front-prod` muestra las dos filas: el deploy de la errata (`1b63182`, PR #35) y debajo el de vuelta a `d7629e6`.

**Lo que el rollback de código no deshace son los datos**: la base de Neon no se toca, así que lo que se haya escrito sigue escrito, y una migración que ya cambió el esquema sigue aplicada. Eso es lo que haría falta cubrir (backups, migraciones reversibles) para que el rollback fuera completo. La métrica DORA que ejercito es el **tiempo de recuperación** (*time to restore*).

### 9. Qué quedó por variable y qué adentro de la imagen

La dirección del backend **no está adentro de la imagen del front**. El nginx del front usa una plantilla (`frontend/default.conf.template`) que la imagen de nginx completa al arrancar con dos variables de entorno: `BACKEND_URL` (la api de ese entorno, sin barra al final) y `DNS_RESOLVER`. Los valores por defecto son los de Docker Compose, así que `docker compose up` sigue andando igual. Cada front apunta a la api de su entorno (QA con QA, PROD con PROD). Lo que queda adentro de la imagen es la aplicación compilada de React, la configuración de nginx (la estructura, no la dirección) y los puertos. Esto es lo que hace posible el TP7: el mismo build, configurado distinto por entorno. Como Render reconstruye cada servicio (punto 6), «la misma imagen» significa el mismo Dockerfile y el mismo commit configurados con variables distintas, no un único artefacto byte a byte.

Las bases son distintas de verdad: creé un dato de prueba solo en `app_prod` y comprobé que en `app_qa` no aparece (1 fila en `app_prod`, 0 en `app_qa`), y después lo borré.

### 10. Problemas encontrados y cómo los resolví

- **Auto-Deploy estaba prendido en los cuatro servicios.** Render los crea con *On Commit*, así que el merge del PR #34 desplegó PROD por su cuenta, sin pasar por mi pipeline ni mi aprobación. Lo descubrí en la lista *Deploys* de `peliculas-front-prod` (una fila con trigger «Auto-Deploy») y lo puse en **Off** en los cuatro.
- **El environment `production` no tenía los reviewers activados** en la primera corrida: el job corrió sin pausarse y no esperó aprobación. Activé *Required reviewers* y volví a correr.
- **Un hook de PROD no llegó al job:** `curl` falló con `URL rejected: Malformed input`, porque la variable estaba vacía (el secret estaba mal cargado). Lo cargué de nuevo con el nombre exacto.
- **Aprobé sin querer** el cambio con la errata: lo recuperé con el rollback del punto 8 y rechacé después con motivo (punto 4).
- **Un servicio de Render con la URL equivocada:** creé el front de QA con un nombre que no cambié, y Render fija la URL al crear el servicio (`ingsoft3-tp01-1pct.onrender.com`), así que renombrarlo no sirve. Lo borré y lo volví a crear con el nombre correcto.
- **Mi `main` local quedó con commits viejos** después de un *squash merge* (Git los ve como «adelantados» aunque el contenido es el mismo). Lo dejé igual a GitHub con `git reset --hard origin/main`.
- **El registro de usuario falló en un entorno recién despertado** con el mensaje genérico del front: era el arranque en frío de la api (punto 5).

### 11. Release

La llegada a PROD se etiqueta como `v6.0.0` y se publica como Release de GitHub con notas; el tag marca el commit que efectivamente está corriendo en PROD, no la punta de la rama.

### 12. Declaración de uso de IA

Usé IA (Claude) para todo el TP6: los dos jobs de deploy del `ci.yml` y su orden, el script `smoke.sh`, la plantilla de nginx y los cambios de los Dockerfiles, la publicación de imágenes a `ghcr.io`, y la guía paso a paso de la configuración de Neon y Render y de los environments de GitHub. También los dos cambios de demostración del gate (el que lleva la errata en el título y el que la corrige) los armó Claude, y el comentario del rechazo lo redacté a partir de un borrador suyo. Lo verifiqué así: ejecuté cada paso yo mismo, comprobé en Render las tres cosas que se revisan (Auto-Deploy en Off en los cuatro servicios, que cada hook dispare su servicio, y que las bases sean distintas con el dato de prueba), probé el gate con un rechazo y una aprobación reales, e hice un rollback real y cronometrado. Puedo explicar qué pasa entre el merge y PROD: cada compuerta está descripta en el punto 3.


## TP7 — Contenedores en el pipeline + pruebas de integración y e2e

### 1. Build once, deploy many: qué problema del TP6 resuelve

En el TP6, mi pipeline promovía un **commit**, y cada entorno lo **reconstruía**: Render clonaba el repo y volvía a compilar. Lo que corría en QA no era lo que habían probado mis tests, sino otra construcción del mismo código, y la de PROD era una tercera, hecha a otra hora. Con mi propio código el riesgo es concreto: la imagen final del backend parte de `alpine:3.20` y hace `apk add ca-certificates`, y el front usa `nginx:1.27-alpine`. Esas etiquetas de imagen base y los paquetes de Alpine se resuelven **en el momento del build**, así que dos builds del mismo commit hechos con días de diferencia pueden no ser idénticos, y yo no tendría cómo demostrarlo.

La solución es la imagen. El pipeline la construye **una sola vez**, ya con todo adentro, la publica en el registry, y los entornos **ejecutan esa misma imagen**: QA y PROD corren el mismo binario. Promover pasa a ser decirle a PROD «ejecutá la imagen `X`», la misma `X` que QA y las pruebas verificaron, y volver atrás es ejecutar la anterior, que sigue intacta en el registry.

### 2. Mi estrategia de etiquetas

- **En el registry:** cada imagen lleva una sola etiqueta, `sha-<commit>`, con los 40 caracteres del commit de `main` que la produjo. Garantiza **trazabilidad exacta**: de qué commit salió. No garantiza inmutabilidad: una etiqueta se puede volver a publicar con otro contenido (por ejemplo, volver a correr la corrida de un commit republica el mismo `sha-…`). Lo único que nadie puede mover es el **digest** (`@sha256:…`), la huella del contenido. Los vi al hacer el `docker pull`; el hook de Render también acepta un digest en lugar del tag, y esa sería la garantía completa.
- **En git:** `v7.0.0` es el **nombre** de una versión que sé que anduvo en PROD. Es el catálogo de versiones: de ahí llego al commit y del commit a la imagen.
- **Por qué mi pipeline no publica `latest`:** `latest` es un puntero que se mueve con cada merge y no dice qué contenido hay detrás. Si algún deploy lo nombrara, ejecutaría «lo último que se publicó» haya pasado las pruebas o no, y después no podría decir qué desplegué. Con una sola etiqueta por imagen, la de su commit, no hay forma de desplegar algo sin decir exactamente qué es.
- **Con qué imagen configuré los servicios de Render, y por qué no es ésa la que corre:** a los cuatro servicios les puse como fuente *Existing Image* con la imagen de un merge viejo (`sha-437f49d…`, el último de `main` cuando los cambié). Render pide la dirección de una imagen para cambiar la fuente y esa era una que ya había comprobado que existía y era pública. Lo que corre en cada entorno lo decide después cada deploy de mi pipeline, con su `imgURL`: *Settings → Image* queda clavado en la primera y no se actualiza. Por eso *Manual Deploy* es peligroso: despliega la imagen de *Settings*, no la que está corriendo, y puede llevar PROD para atrás sin avisar.
- **Servicios nuevos o viejos:** no creé servicios nuevos; **le cambié la fuente a los cuatro existentes**. Así conservé las URLs, las variables de entorno, los deploy hooks y el historial. No hay servicios viejos que borrar.

### 3. De la release a la imagen, en un paso

Del tag `v7.0.0` al commit: `git rev-list -n1 v7.0.0`. Ese commit es el de la release, y mi pipeline ya publicó las dos imágenes de ese commit con su nombre, `sha-<commit>`: la del backend y la del frontend llevan la misma etiqueta. Basta abrir la release en GitHub, ver a qué commit apunta, y buscar esa etiqueta en los dos paquetes. Volver a esa versión sería mandarle esa imagen al hook de cada servicio.

### 4. Cómo compruebo, desde afuera, que un entorno ejecuta mi imagen

En los *Events* de cada servicio de Render, el deploy dice «Triggered via Deploy Hook» y nombra el `sha-<commit>` de esa corrida, y esa misma etiqueta existe en mis dos paquetes. Lo comprobé en los **cuatro** servicios (api y front, en QA y en PROD). Un deploy hecho a mano desde el panel no diría «Triggered via Deploy Hook» y no habría corrida que lo respalde. Además, el `docker pull` de cada imagen por su etiqueta funciona desde cualquier máquina sin credenciales.

**Lo que mi smoke no alcanza a probar:** solo pregunta si el entorno **responde**. Si el deploy falla, el servicio sigue sirviendo la versión anterior y contesta igual: el smoke da verde contra la versión vieja. En mi corrida roja, `deploy-qa` terminó en 12 segundos, una pista de que el smoke contestó mientras la imagen nueva todavía se estaba desplegando. Se cerraría haciendo que el `/health` de la api devuelva el commit con el que se construyó la imagen (un `ARG` en el Dockerfile) y que el smoke lo compare con el commit de la corrida. Es el paso siguiente; en este TP no lo implementé.

### 5. Qué flujos elegí, y qué NO puse en cada suite

Mi app tiene login, y las películas y géneros requieren sesión. Las dos suites usan una **cuenta de pruebas única** en QA (`bot-e2e@ejemplo.test`), que el propio test crea por la api si no existe, en vez de una cuenta por corrida, porque los usuarios no se pueden borrar y se acumularían.

**Integración** (`frontend/e2e/api.spec.js`, sin navegador, contra la api de QA, con la base de verdad):
- *El alta guarda y el borrado saca.* Crea una película con título único (con la hora adentro), comprueba que la base la devuelve, la borra y comprueba que ya no está. Si falla, el alta o el borrado no llegan a la base.
- *Título vacío:* la api responde 400 con su mensaje y no se creó nada. Si falla, la api aceptó un dato inválido.
- *La tercera, elegida por mí: un género con películas no se puede borrar.* Responde 409 y el género sigue existiendo, y recién sin películas se puede borrar. La elegí porque la regla vive en una **consulta a la base** que un doble no ve, y porque si se rompe, la gente pierde la relación entre películas y géneros sin enterarse. Para comprobar que atrapa de verdad, se rompió a propósito esa regla en una copia del backend (con ayuda de Claude) y la prueba se puso roja.
- Cada prueba borra lo que crea, comprobándolo, aunque falle una aserción.

**e2e** (`frontend/e2e/peliculas.spec.js`, Playwright con un navegador real, contra el front de QA):
- *Crear una película la muestra en la lista, y eliminarla la saca.* Llena el formulario, la ve en la tabla con su título único, la elimina con la confirmación y comprueba que desapareció.
- *Un año imposible muestra el error y no crea nada.* Escribe el año 1500, ve el mensaje y, tras recargar la página, la base no tiene esa película.
- *Iniciar sesión muestra la colección y cerrar sesión vuelve al acceso.* Es el flujo que mi usuario hace todos los días: si mañana no anda, me escriben. Toca la base (el login consulta el usuario).

**Qué NO puse:** en las e2e no repetí las reglas de validación campo por campo (eso lo cubren los unitarios del TP5 y la integración), ni los flujos de editar, filtrar o puntuar: son caros, pocos y se pueden romper por motivos que no son del código. En la integración no puse listados con filtros, errores de autenticación ni estadísticas: cada capa verifica lo suyo. Los flujos que ninguna de las dos recorre quedan sin red final.

### 6. Qué probás en integración y qué en e2e, y el par verde/rojo

La integración le habla a la api directo y comprueba que **la api y su base se entienden**: lo que la base guarda, rechaza y devuelve. La e2e comprueba que **una persona puede usar la app** completa, y por eso ve el front. No son lo mismo: una e2e verde no me dice que la api esté bien por dentro, y una integración verde no me dice que la pantalla use bien esa api.

**El diagnóstico de mi corrida roja.** En la corrida del commit que rompió la app (cambié el texto del botón de alta en el front), el smoke test dio verde, la integración dio verde y la e2e dio roja, y `deploy-prod` no arrancó: no llegó a pedir aprobación.

Ese par de colores dice quién se rompió sin abrir el código. La integración le habla a la api de QA directo, con los nombres correctos de los campos, y la base guardó y borró sin problemas: la api y la base están sanas. El smoke también dio verde porque solo pregunta si la api, la base y el front responden. Lo único que falló es un navegador usando la app como una persona, y eso apunta al front: algo de lo que ve el usuario cambió. Es la fila de la tabla de la guía en la que la integración pasa y la e2e no: la rotura está en cómo el front usa la api o en lo que muestra, no en la api.

Con los reportes lo confirmé: el reporte de integración tiene los 3 tests en verde, y el de e2e muestra los tests que usan el formulario de alta fallando porque no encuentran el botón «Agregar película» (con la captura y la traza del reintento). El test de iniciar y cerrar sesión no depende de ese botón.

Si hubiera roto la api en lugar del front, la integración se habría puesto roja y la e2e ni siquiera habría corrido, porque su `needs` no se cumple: eso también es información, y habría apuntado directo a la api.

**Lo que el gate evitó:** que PROD recibiera una versión en la que la gente ya no puede agregar películas con el flujo que probamos, sin que nadie lo notara, porque el smoke y los unitarios daban verde.

**Lo que mi gate no ataja:** un bug en un flujo que ninguna de mis 3 e2e recorre (editar, filtrar, puntuar), un error que aparece solo con datos o volumen reales, o un problema visual que no cambia ningún texto ni botón que el test busque.

La rotura la provoqué **en la aplicación** (el texto del botón de alta en `FormularioPelicula.jsx`), sin tocar nada de `e2e/`, y elegí una que ningún unitario ataja (no testeo componentes). El arreglo restauró el texto y la cadena completa quedó verde hasta PROD.

### 7. Por qué mi integración es la amplia y no la estrecha

Mi integración le habla a la **api ya desplegada en QA**, con su base de verdad, en vez de armar la api dentro del pipeline con una base descartable. **Qué gano:** casi nada de configuración (no necesito un Postgres como servicio del job ni una forma de levantar la api en memoria), y además comprueba que el **despliegue** quedó bien: ejecuta la imagen exacta de esa corrida en el entorno real. **Qué pierdo:** aislamiento y velocidad de arranque. Depende de que QA esté despierto y desplegado, comparte la base con otras corridas y con la e2e, y necesita esperar al deploy. La estrecha llegaría antes, sin deploy, y con una base que nace y muere con la corrida, pero cuesta bastante más de armar.

### 8. El cold start del free tier en las e2e, y los tests flaky

En el plan gratis de Render, QA se duerme por inactividad y despierta lento (hasta un minuto). Mi configuración lo maneja con **tiempos largos y un solo reintento, nunca con `sleep`**: 60 segundos por test, 15 segundos por aserción y por acción (un `click` sobre un botón que no existe falla en 15 s y no en un minuto), y una espera automática a que el elemento aparezca. Además, el smoke y la integración corren antes y ya despertaron QA cuando arranca la e2e. Para que el entorno no falle en el primer pedido, hice que el front le mande un «hola» (`/health`) a su api cuando arranca el contenedor y al abrir la página, y que nginx espere hasta 120 segundos en vez de cortar.

Un test **flaky** es el que a veces pasa y a veces falla sin que el código haya cambiado. Es peor que no tener test porque entrena al equipo a ignorar el rojo. Mi configuración usa `retries: 1`: si un test pasa recién en el reintento, el reporte lo marca como *flaky* aunque la corrida salga en verde, así que lo veo y lo arreglo en vez de taparlo. El resumen de cada corrida muestra pasados, fallados, flaky y salteados.

### 9. Cómo logré que la misma imagen del front sirva en QA y en PROD

La dirección de la api **no está dentro de la imagen del front**. El nginx del front usa una plantilla que la imagen de nginx completa al arrancar con dos variables de entorno: `BACKEND_URL` (la api de ese entorno, sin barra al final) y `DNS_RESOLVER`. Cada servicio de Render tiene sus propios valores. Lo que queda **adentro** de la imagen es la aplicación compilada de React, la estructura de la configuración de nginx y los puertos. Ahora sí es *literalmente* la misma imagen en los dos entornos: lo único que los distingue son esas variables, que es lo que hace posible promover lo mismo que se verificó.

### 10. Límites conocidos de mi cadena

- **QA es una sola para todas las corridas.** Con dos merges seguidos, la corrida B puede redesplegar QA mientras la integración o la e2e de la A todavía la usan, y entonces un rojo (o un verde) puede no ser de mi código. Lo reconocería mirando en *Events* de QA qué imagen está corriendo y a qué hora corrió el `deploy-qa` de la otra corrida. Lo que haría: rechazar la corrida A con el motivo y quedarme con la B, sin re-correr sus pruebas (probarían la imagen de B). La regla que sigo es **un merge por vez** mientras la cadena corre. Un equipo real lo resuelve con un entorno por corrida.
- **El smoke contesta con la versión vieja** mientras se despliega la nueva (punto 4).
- **El gate humano depende de un click:** en el TP6 aprobé por error un cambio con una errata. La e2e es la red que atraparía algo así antes de llegar a mí.

### 11. Problemas encontrados y cómo los resolví

- **Error al abrir QA por el arranque en frío:** cada vez que abría el front de QA con la api dormida, aparecía «Ocurrió un error inesperado»: nginx cortaba antes de que la api despertara y devolvía un 502 o 504 en HTML, que el front no entendía. Lo resolví con tiempos de espera largos en nginx, un «hola» a la api desde el contenedor del front y al cargar la página, y un mensaje claro para los 502, 503 y 504, con sus tests.
- **Mi front no tenía nombres accesibles que Playwright pudiera encontrar:** los errores no tenían `role="alert"` y el botón «Eliminar» de cada fila no distinguía de qué película era. Los agregué, en vez de esquivarlo con selectores de CSS, que cambian con el diseño. Además, los campos «Año», «Género» y «Estado» están repetidos entre el formulario de alta y los filtros, y hay dos botones «Iniciar sesión» (la pestaña y el envío): acoté cada búsqueda al formulario correspondiente.
- **Las películas requieren sesión:** las pruebas necesitaban un usuario. Creé una cuenta de pruebas única que el propio test registra por la api, o inicia sesión si ya existía.
- **Un cambio de fuente en Render que no despliega:** cambiar la fuente a *Existing Image* no despliega nada; tuve que probar un servicio a mano con su hook antes de pasar a los otros tres, y recordar que *Manual Deploy* ya no es seguro (punto 2).
- **Errores míos con Git:** un `git add` con rutas de la raíz estando parado en `frontend/` (`pathspec did not match`), y correr las pruebas antes de copiar los archivos de las suites (`No tests found`). Se resolvieron volviendo a la raíz y copiando los archivos.

### 12. Declaración de uso de IA

Usé IA (Claude) para todo el TP7: los jobs de deploy con `imgURL`, los jobs `integracion` y `e2e` y el gate, las dos suites de Playwright (`api.spec.js` y `peliculas.spec.js`) y su configuración, los cambios de accesibilidad del front, el arreglo del arranque en frío, el cambio que rompió la aplicación a propósito y su arreglo, y el borrador del diagnóstico del punto 6, que corregí con lo que vi en la corrida. Lo verifiqué así: corrí las dos suites contra mi `docker compose` local antes de subirlas (3 de integración y 3 de e2e en verde), comprobé en *Actions* la cadena completa (la pausa de revisión y quién aprobó), vi la corrida roja con la integración verde y `deploy-prod` sin arrancar, miré en los *Events* de los cuatro servicios el deploy por hook con la etiqueta `sha-` de la corrida, y comprobé que mis pruebas detectan cambios (Claude rompió a propósito una regla del backend en una copia y la prueba correspondiente se puso roja). Puedo explicar qué verifica cada test (punto 5), contra qué entorno corre (la integración contra la api de QA y la e2e contra el front de QA) y qué pasa si falla: la integración roja frena todo y la e2e ni corre, y si la e2e queda roja, `deploy-prod` no pide aprobación.

## TP8 — Infraestructura como código: `preprod` declarado en `infra/`

### 1. Riel elegido y versiones fijadas

Elegí el riel **Terraform + Neon + Render**, el mismo de mi QA y producción desde el TP6, así que no necesité cuentas nuevas. Usé Terraform **v1.16.6** (no OpenTofu). Los providers están fijados por versión en `infra/main.tf`, y el `.terraform.lock.hcl` fija además los hashes exactos:

- `kislerdm/neon` **0.18.0** (su repositorio está archivado: el oficial saldrá con otro nombre de fuente, y ese día el ejercicio es cambiar el `source` y correr `terraform init -upgrade`),
- `render-oss/render` **1.9.1**,
- `hashicorp/random` **3.9.1**, que agregué yo: la guía no lo usa.

Declaré un entorno **nuevo**, `preprod`, y no toqué QA ni producción: ni los importé ni los modifiqué. Un proyecto de Neon propio (`peliculas-preprod`) con su rol y su base, y dos servicios de Render (`peliculas-api-preprod` y `peliculas-front-preprod`) que corren las **imágenes de mi TP7**, la del commit de la release `v7.0.0`. Elegí un entorno nuevo porque se puede **destruir entero** sin rozar nada de lo que ya tengo y porque Render no deja modificar un servicio del plan gratuito: acá se crea y se destruye, no se edita.

Ninguna credencial está escrita en el código: los providers leen `NEON_API_KEY`, `RENDER_API_KEY` y `RENDER_OWNER_ID` del entorno, y el id de mi organización de Neon entra por `TF_VAR_neon_org_id`. La conexión **se deriva, no se copia**: `DB_HOST`, `DB_USER`, `DB_PASSWORD` y `DB_NAME` de mi api salen de lo que devolvió Neon, y el `BACKEND_URL` del front sale de la URL de su propia api. El `JWT_SECRET` lo genera Terraform con `random_password` (queda solo en el estado, que no se versiona); lo haría distinto con un gestor de secretos. El estado vive en mi máquina y está cubierto por el `.gitignore` desde antes del primer `apply`; el `.terraform.lock.hcl` sí se commitea.

Lo comprobé como entorno y no como recursos: el `/health` de la api, un login de prueba que toca la tabla de usuarios, y el front cargando y llegando a su propia api (el mismo smoke del pipeline). Lo que sigue siendo artesanal en mi cuenta: QA y producción, que cliqueé a mano y ningún plan mío ve (un `plan` limpio no dice nada sobre ellos), y los deploy hooks de Render, que se copian del panel.

### 2. Problemas que tuve y cómo los resolví

- **Instalar Terraform con Chocolatey falló por permisos** (`Acceso denegado a C:\ProgramData\chocolatey`): la terminal no era de administrador. Lo resolví y quedó instalada la v1.16.6.
- **El primer `apply` falló con HTTP 400 `org_id is required`.** Neon ahora exige la organización dueña de cada proyecto, y mi código no la declaraba. Agregué `org_id` al proyecto, con una variable sin valor por defecto que se pasa por el entorno. Después del error, el `plan` siguiente no volvió a crear la contraseña que ya existía: Terraform retomó donde había quedado gracias al estado.
- **El segundo `apply` falló en Render: `invalid ownerID`.** Había cargado como `RENDER_OWNER_ID` el id de la URL del panel, que es el de mi usuario. Render pide el del workspace; el correcto lo saqué de su API (`/v1/owners`). Esa vez Neon ya había creado sus recursos y solo faltaban los dos servicios.
- **Hallazgos que la guía dice que la cátedra no midió:** al recrear los servicios de Render con el mismo nombre, volvieron **las mismas URLs**, sin sufijo; y en cambio el **host de la base de Neon cambió** (el proyecto nuevo trae su propio endpoint), y mi api de preprod apuntó sola a la base nueva, porque ese valor se deriva y no está escrito en ningún lado.
- **`prevent_destroy`:** probé la protección con preprod vivo y el `destroy` se negó (`Instance cannot be destroyed`) antes de pedir confirmación. La saqué después, porque preprod tiene que poder morir.

### 3. Declaración de uso de IA

Usé IA (Claude) para el TP8: el código de `infra/` (providers, recursos, variables y salidas, adaptados a las variables de mi backend en Go), el `.gitignore` con las reglas de Terraform y la guía de cada paso. Claude lo validó con OpenTofu contra el esquema real de los providers (formato y validación); el `plan` y el `apply` los corrí yo contra mis cuentas, no se probaron antes. Lo verifiqué así: leí el plan antes de aplicar, comprobé idempotencia (el segundo `plan` dijo «No changes»), abrí el estado y encontré la contraseña en texto plano, provoqué drift (borré la base a mano y el plan lo detectó y el `apply` la reparó; creé una base que mi código no declara y el plan no la vio), hice el ciclo `destroy` → `apply` mirando el orden inverso de la ejecución, y comprobé en las consolas de Neon y Render que mi QA y mi producción no se tocaron. Puedo leer un plan en voz alta y explicar qué haría cada acción, de dónde se derivan la conexión y la dirección del backend, y por qué el estado no está en el repositorio.

## TP9 — DevSecOps y observabilidad: dependencias vigiladas, un monitor externo y Sentry

### 1. El PR de Dependabot que mergeé, y por qué me animé

Configuré Dependabot (`.github/dependabot.yml`) para los tres ecosistemas de mi app: `gomod` en `/backend`, `npm` en `/frontend` y `github-actions` en la raíz, todos semanales. A los pocos minutos abrió **9 pull requests**. Cuatro pasaron mi pipeline en verde (tres actions y un paquete de Go) y cinco se pusieron en rojo, todos saltos de versión **mayor**: React y React DOM 18 a 19, vitest y su cobertura de 2 a 5, y vite de 5 a 8. Esos cinco no los mergeé: mi pipeline dijo que no, y dejé que lo dijera.

El que mergeé es el **#47: `gorm.io/driver/postgres` de 1.6.2 a 1.6.3** en `/backend`, el driver con el que mi backend habla con PostgreSQL. Me animé por tres cosas: es un cambio de *patch* (el de menor riesgo en SemVer), es una dependencia que corre en producción y toca la base, y mi CI estaba en verde sobre ese PR, con los tests de integración del TP5 contra un PostgreSQL de verdad, que son la prueba directa de que el driver nuevo anda. Esa confianza no la da Dependabot, que solo propone: la dan los tests, el coverage y el gate que armé durante el semestre. Después del merge, la cadena completa (`build → deploy-qa → integración → e2e → aprobación → deploy-prod`) quedó en verde.

- PR mergeado: https://github.com/benjaminatias/ingsoft3-tp01/pull/47
- Su corrida de `main` en verde: https://github.com/benjaminatias/ingsoft3-tp01/actions/runs/40124291026

### 2. El monitor y lo que mide Sentry

**El monitor.** Puse un monitor externo de UptimeRobot sobre el `/health` de mi API de producción, con un intervalo de 5 minutos. Para eso tuve que arreglar el `/health`, que solo contestaba GET: el plan gratuito pregunta con HEAD y habría marcado la app caída estando viva. Después provoqué la caída suspendiendo el servicio en Render. Lo suspendí a las **17:02** y el aviso de **DOWN** llegó a las **17:05**: el monitor **tardó unos 3 minutos** en avisarme. Reanudé a las 17:06 y el aviso de **UP** llegó a las 17:10. Sin el monitor me habría enterado cuando un usuario me lo contara. Lo pausé al terminar, porque despierto todo el mes se comería las horas gratuitas de Render que comparten todos mis servicios.

**Sentry.** Instalé Sentry en el backend (el DSN entra por la variable de entorno `SENTRY_DSN` en Render) y en el front (el DSN va escrito en el código, y no es un problema: solo sirve para *enviar* datos, nunca para leerlos). Provoqué un error real en producción con un endpoint temporal que después saqué, y quedó registrado («Falla de prueba del TP9», 4 eventos). De la tabla de transacciones del backend sale lo siguiente:

- **Mi operación más usada** es `GET /api/peliculas`, empatada con `GET /api/estadisticas`, con 15 pedidos cada una: el front pide las dos cada vez que muestra la lista.
- **Mi operación más lenta** es `POST /api/auth/login`, con un **P95 de unos 965 ms** y un P50 de 364 ms; el registro le sigue con unos 341 ms. Tarda unas diez veces más que el resto porque comparar el hash con bcrypt es lento a propósito, y mi servicio gratuito tiene poca CPU. Lo curioso es que **la más usada es de las más rápidas** (un P95 de unos 36 ms) y la más lenta es de las menos usadas (5 logins), así que optimizarla casi no cambiaría lo que siente la gente.
- Con ese dato abrí una tarjeta en mi tablero del TP3: *Revisar el costo del login: P95 de 965 ms en producción (Sentry)*, https://github.com/benjaminatias/ingsoft3-tp01/issues/57. La decisión que dejé escrita es **no bajar el costo de bcrypt**, porque ese costo es la protección.

### 3. Problemas que tuve y declaración de uso de IA

- **Mi `/health` solo contestaba GET** y el monitor gratuito pregunta con HEAD. Lo cambié para que acepte los dos y agregué un test.
- **El tablero de rendimiento del backend quedó vacío** aunque el error sí llegaba a Sentry. La causa era mía: en la versión de la librería que uso (`sentry-go` 0.50) hay que activar `EnableTracing` además de la tasa de muestreo, y yo había puesto solo la tasa. Lo detecté al comparar las dos pantallas (el error aparecía, pero el tablero decía «Get started with backend tracing») y lo arreglé en un PR aparte.
- **Mis primeros números de latencia eran del smoke test del pipeline**, no míos: el login más lento tenía 100 % de estado «no autenticado» y ningún usuario, que es el login falso que hace el smoke al terminar cada deploy. Lo reconocí por eso y volví a medir usando la aplicación de producción.
- **Cargar una variable en Render puede desplegar la imagen vieja.** Mis servicios corren una imagen fija en *Settings*, así que guardé `SENTRY_DSN` con «Save only», sin desplegar, y la tomó el deploy de mi pipeline.
- **Sentry creó el proyecto del backend con un nombre automático** (`go`) y le cambié el nombre a `peliculas-backend`.

**Uso de IA:** usé IA (Claude) para el `dependabot.yml`, la integración de Sentry en el backend (Go) y en el front (React), el cambio de `/health` para aceptar HEAD, el endpoint temporal de prueba, el borrador de la descripción de la tarjeta y la guía de los pasos en GitHub, Sentry y UptimeRobot. Lo verifiqué así: comprobé en producción que `/health` contesta 200 a HEAD, que el error de prueba aparecía en Sentry, que el monitor mandaba el DOWN y el UP, y que el pipeline completo pasó en verde en cada PR, con la aprobación de producción. Claude comprobó con un test que el panic de prueba llega a Sentry sin romper el servidor, y que sin `EnableTracing` el backend no manda mediciones y con él sí. Puedo mostrar cada pieza funcionando y explicar qué mira y qué no: Dependabot propone actualizaciones pero no las prueba, CodeQL lee el código sin ejecutarlo, el escáner de secretos solo detecta formatos conocidos, y el monitor sabe si la app responde pero no por qué falla (eso lo dice Sentry).