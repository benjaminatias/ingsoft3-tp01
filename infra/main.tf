# preprod: un entorno NUEVO, entre QA y producción, declarado como código.
# No toca QA ni producción: convive con ellos y no los importa. Nace con `terraform apply`
# y muere con `terraform destroy`.

terraform {
  required_version = ">= 1.5"

  required_providers {
    # Fijados por versión: los nombres de atributos cambian entre versiones de un provider.
    neon = {
      source  = "kislerdm/neon"
      version = "0.18.0" # su repositorio está archivado: el oficial saldrá con otro nombre de fuente
    }
    render = {
      source  = "render-oss/render"
      version = "1.9.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.9.1"
    }
  }
}

# Ninguna credencial está escrita acá: los providers las leen del entorno.
provider "neon" {}   # lee NEON_API_KEY
provider "render" {} # lee RENDER_API_KEY y RENDER_OWNER_ID

# ─────────────────────────────── La base de datos de preprod ───────────────────────────────

# Un proyecto propio para preprod (no una base más en el de tu app): se destruye entero sin
# rozar QA ni producción.
resource "neon_project" "preprod" {
  name      = "peliculas-preprod"
  org_id    = var.neon_org_id # Neon exige la organización dueña del proyecto: sin esto, el apply da HTTP 400
  region_id = "aws-us-east-2" # Ohio, la misma región de los servicios de Render

  # Sin esta línea el primer apply falla: el provider pide por defecto un día de historial
  # y el tope del plan gratuito es de seis horas (21600 segundos).
  history_retention_seconds = 21600
}

resource "neon_role" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id # default_branch_id, no branch_id
  name       = "app_preprod_owner"
}

resource "neon_database" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "app_preprod"
  owner_name = neon_role.preprod.name
}

# El secreto con el que la api firma los tokens de sesión: lo genera Terraform, vive solo en el
# estado (que no se versiona) y nadie lo escribe ni lo copia. Cada preprod nace con uno distinto.
resource "random_password" "jwt" {
  length  = 48
  special = false
}

# ─────────────────────────────── Los dos servicios de preprod ───────────────────────────────
# Corren las IMÁGENES del TP7, ya construidas: Render no compila nada. La etiqueta es fija a
# propósito: preprod corre la imagen de un commit concreto, escrito en el código.
# Render NO deja modificar un servicio del plan gratuito: acá se crea y se destruye, no se edita.

resource "render_web_service" "api" {
  name   = "peliculas-api-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_api # SIN la etiqueta pegada
      tag       = var.image_tag      # la etiqueta va en su propio campo
    }
  }

  health_check_path = "/health"

  # La conexión a la base NO se copia: cada valor se DERIVA de lo que devolvió Neon.
  # No se declara APP_COMMIT ni nada de la imagen: una variable del servicio pisa a la de la imagen.
  env_vars = {
    DB_HOST     = { value = neon_project.preprod.database_host }
    DB_PORT     = { value = "5432" }
    DB_USER     = { value = neon_role.preprod.name }
    DB_PASSWORD = { value = neon_role.preprod.password }
    DB_NAME     = { value = neon_database.preprod.name }
    DB_SSLMODE  = { value = "require" } # Neon exige conexión segura
    JWT_SECRET  = { value = random_password.jwt.result }
    SERVER_PORT = { value = "8080" }
    PORT        = { value = "8080" }
  }
}

resource "render_web_service" "front" {
  name   = "peliculas-front-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_front
      tag       = var.image_tag
    }
  }

  env_vars = {
    # La dirección del backend SALE de la api que declara este mismo archivo, sin la barra final
    # (con ella nginx mandaría todo a "/" y cada llamada daría 404).
    BACKEND_URL  = { value = trimsuffix(render_web_service.api.url, "/") }
    DNS_RESOLVER = { value = "8.8.8.8" }
    PORT         = { value = "80" }
  }
}
