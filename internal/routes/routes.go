package routes

import (
	"github.com/gin-gonic/gin"

	"ejemplo/crud-go-sqlite/internal/handlers"
)

// Nuevo crea el enrutador de Gin con todas las rutas de la API.
func Nuevo(productos *handlers.Productos) *gin.Engine {
	r := gin.Default() // Incluye los middleware Logger y Recovery.
	r.SetTrustedProxies(nil)

	grupo := r.Group("/productos")
	grupo.POST("", productos.Crear)
	grupo.GET("", productos.Listar)
	grupo.GET("/:id", productos.Obtener)
	grupo.PUT("/:id", productos.Actualizar)
	grupo.DELETE("/:id", productos.Eliminar)
	return r
}
