# Solo salen cosas que NO son secretas. La contraseña y la cadena de conexión no se imprimen:
# viajan del proveedor al servicio sin pasar por ningún lado visible.

output "db_host" {
  value = neon_project.preprod.database_host
}

output "database" {
  value = neon_database.preprod.name
}

output "api_url" {
  value = render_web_service.api.url
}

output "front_url" {
  value = render_web_service.front.url
}
