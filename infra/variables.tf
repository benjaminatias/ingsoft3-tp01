# Estos valores NO son secretos (el nombre de una imagen pública, una etiqueta, una región),
# así que van con default en un archivo que SÍ se commitea.

variable "image_repo_api" {
  description = "Imagen del backend, SIN la etiqueta."
  default     = "ghcr.io/benjaminatias/ingsoft3-tp01-backend"
}

variable "image_repo_front" {
  description = "Imagen del frontend, SIN la etiqueta."
  default     = "ghcr.io/benjaminatias/ingsoft3-tp01-frontend"
}

variable "image_tag" {
  description = "Etiqueta de las dos imágenes: la del commit de la release v7.0.0."
  default     = "sha-b7aa87bf6d3e468ba96ca6a8aa5b7c3f12536dd2"
}

variable "region" {
  description = "Región de Render (frankfurt | ohio | oregon | singapore | virginia)."
  default     = "ohio"
}

# El id de tu organización de Neon (empieza con "org-"). No es una credencial, pero tampoco se
# escribe en el código: no tiene default y se pasa por el entorno, igual que las claves:
#   export TF_VAR_neon_org_id="org-..."
variable "neon_org_id" {
  description = "Id de la organización de Neon dueña del proyecto de preprod."
  type        = string
}
