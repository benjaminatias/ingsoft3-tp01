#!/bin/sh
# Se ejecuta cuando arranca el contenedor del front (la imagen de nginx corre todo lo que haya en
# /docker-entrypoint.d/ antes de levantar nginx). Le manda un «hola» a la API de ESTE entorno
# para que empiece a despertar YA, en paralelo, sin esperar a que alguien abra la página.
# En el plan gratis de Render la API se duerme por inactividad: sin esto, los dos arranques en
# frío (front y después API) se suman; con esto se solapan.
#
# Va en segundo plano y nunca falla: si la API no contesta, el front arranca igual.
( wget -q -T 120 -O /dev/null "${BACKEND_URL}/health" >/dev/null 2>&1 || true ) &
exit 0
