#!/usr/bin/env bash
# Smoke test post-deploy: espera a que un entorno (API + base + front) conteste.
# Uso:  bash scripts/smoke.sh <URL_API> <URL_FRONT>
# Reintenta porque en el plan gratis de Render el entorno tarda en estar listo
# (build + arranque en frío): un curl seco daría falsos rojos.
set -u

URL_API="${1:?falta la URL de la API}"
URL_FRONT="${2:?falta la URL del front}"
INTENTOS="${SMOKE_INTENTOS:-30}"
ESPERA="${SMOKE_ESPERA:-20}"

verificar() {
  # 1. El proceso está vivo y alcanza la base (el /health hace ping a PostgreSQL).
  curl -fsS --max-time 10 "$URL_API/health" > /dev/null || return 1

  # 2. La app consulta la tabla de usuarios: un login inexistente tiene que dar 401.
  #    Si la base estuviera vacía o inalcanzable, daría 500, y un 200 sería rarísimo.
  local codigo
  codigo=$(curl -s --max-time 10 -o /dev/null -w '%{http_code}' \
    -X POST -H 'Content-Type: application/json' \
    -d '{"email":"smoke@invalido.test","password":"no-existe-123"}' \
    "$URL_API/api/auth/login")
  [ "$codigo" = "401" ] || { echo "  login de prueba devolvió $codigo (se esperaba 401)"; return 1; }

  # 3. El front sirve la página.
  curl -fsS --max-time 10 "$URL_FRONT/" > /dev/null || return 1

  # 4. El front llega a la API de SU entorno a través de nginx.
  curl -fsS --max-time 10 "$URL_FRONT/health" > /dev/null || return 1
}

echo "Esperando a que responda: API=$URL_API  FRONT=$URL_FRONT"
for i in $(seq 1 "$INTENTOS"); do
  if verificar; then
    echo "✅ Responde (API + base + front + proxy)"
    exit 0
  fi
  echo "  intento $i/$INTENTOS sin respuesta, reintento en ${ESPERA}s"
  sleep "$ESPERA"
done
echo "❌ No respondió a tiempo"
exit 1
