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

En el **backend** puse **70 % de sentencias**, y hoy da **77,3 %**. Go no tiene métrica de ramas,
solo sentencias, así que no puedo reportar el número de rama de ese lado. En el **frontend** puse
**80 % de líneas y 85 % de ramas**: medí 82,4 % y 87,0 %, y dejé unos 2 puntos de margen, chico a
propósito para frenar una caída real sin fallar por ruido. Puse umbral en las dos métricas porque
con vitest 2.1.9 una función que nadie llama suma líneas sin cubrir pero casi no mueve las ramas
(en mi demostración las ramas quedaron en 86,4 %, arriba del umbral, y frenó por líneas). Hoy el
frontend está en 84,7 % de líneas y 89,3 % de ramas, con los tests de `recomendaciones.js`.

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
