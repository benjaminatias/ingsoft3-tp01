#!/bin/sh
# Corre los tests del backend midiendo cobertura y FRENA si no llega al umbral.
# Go no trae una bandera de umbral (a diferencia de vitest), así que el umbral se
# resuelve acá: se lee el total que calcula `go tool cover` y se lo compara.
#
# Lo usa la etapa `test` del Dockerfile como ENTRYPOINT: es la única receta para
# testear, en tu máquina y en el pipeline.
set -eu

UMBRAL="${COVERAGE_MIN:-70}"   # % mínimo de sentencias cubiertas
SALIDA="${COVERAGE_DIR:-/out}"  # carpeta donde quedan los reportes

# Binario estático: no hace falta un compilador de C en la imagen.
export CGO_ENABLED=0

mkdir -p "$SALIDA"

if [ -z "${DB_HOST:-}" ]; then
  echo "AVISO: DB_HOST no está definida. Los tests de integración (que ejercitan"
  echo "       películas y géneros contra PostgreSQL) se van a omitir y el número"
  echo "       de cobertura va a dar mucho más bajo."
fi

# -coverpkg=./internal/...  mide el código de la app aunque los tests vivan en otro
# paquete (./tests). Queda afuera cmd/api (el arranque: main.go) por no estar en internal.
# Los modelos son structs sin sentencias ejecutables: no suman ni restan.
ESTADO=0
go test ./tests/... -count=1 -coverpkg=./internal/... -coverprofile="$SALIDA/coverage.out" || ESTADO=$?

if [ ! -f "$SALIDA/coverage.out" ]; then
  echo "ERROR: los tests no generaron el archivo de cobertura."
  exit 1
fi

go tool cover -func="$SALIDA/coverage.out" > "$SALIDA/coverage-func.txt"
go tool cover -html="$SALIDA/coverage.out" -o "$SALIDA/coverage.html"

TOTAL=$(awk '/^total:/ { gsub("%", "", $3); print $3 }' "$SALIDA/coverage-func.txt")

# Resumen en Markdown por archivo (sentencias cubiertas / totales), para el Summary de la corrida.
# Formato del perfil: archivo:inicio,fin  nº-sentencias  veces-ejecutada
{
  echo "### Coverage del backend (Go)"
  echo ""
  echo "| métrica | valor |"
  echo "|---|---|"
  echo "| sentencias cubiertas | **${TOTAL}%** |"
  echo "| umbral exigido | ${UMBRAL}% |"
  echo ""
  echo "| archivo | sentencias | cubiertas | % |"
  echo "|---|---|---|---|"
  awk 'NR > 1 {
         split($1, p, ":"); f = p[1]
         sub(/^gestor-peliculas\//, "", f)
         tot[f] += $2
         if ($3 > 0) cub[f] += $2
       }
       END {
         for (f in tot) printf "| %s | %d | %d | %.1f |\n", f, tot[f], cub[f], 100 * cub[f] / tot[f]
       }' "$SALIDA/coverage.out" | sort
} > "$SALIDA/resumen.md"

echo ""
echo "Cobertura total (sentencias): ${TOTAL}%  |  umbral: ${UMBRAL}%"

if [ "$ESTADO" -ne 0 ]; then
  echo "ERROR: hay tests que fallan (código $ESTADO)."
  exit "$ESTADO"
fi

if ! awk -v t="$TOTAL" -v u="$UMBRAL" 'BEGIN { exit !(t + 0 >= u + 0) }'; then
  echo "ERROR: la cobertura ${TOTAL}% no llega al umbral de ${UMBRAL}%."
  exit 1
fi

echo "OK: la cobertura cumple el umbral."
