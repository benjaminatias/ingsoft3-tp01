# Etapa 1: build de React con Vite
FROM node:20-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY index.html vite.config.js ./
COPY src ./src

RUN npm run build

# Etapa 2: tests (parte del build: ya tiene Node y las dependencias de desarrollo).
# Va en el medio: la última etapa del archivo tiene que seguir siendo nginx.
FROM build AS test
COPY tests ./tests
ENTRYPOINT ["npm", "run", "test:ci"]

# Etapa 3: Nginx sirve los archivos compilados
FROM nginx:1.27-alpine

# La plantilla va en templates/ (no en conf.d/): al arrancar, la imagen de nginx reemplaza
# ${BACKEND_URL} y ${DNS_RESOLVER} por los valores del entorno. Los valores por defecto son
# los de Docker Compose, así que `docker compose up` sigue andando sin tocar nada.
COPY default.conf.template /etc/nginx/templates/default.conf.template
# Un «hola» a la API al arrancar el contenedor (la imagen de nginx corre /docker-entrypoint.d/*.sh).
COPY --chmod=755 40-despertar-api.sh /docker-entrypoint.d/40-despertar-api.sh
COPY --from=build /app/dist /usr/share/nginx/html

ENV BACKEND_URL=http://backend:8080 \
    DNS_RESOLVER=127.0.0.11

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
