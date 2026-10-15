package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ejemplo/crud-go-sqlite/internal/models"
)

// Productos agrupa los handlers del CRUD y la conexion que comparten.
type Productos struct {
	db *gorm.DB
}

func NuevoProductos(db *gorm.DB) *Productos {
	return &Productos{db: db}
}

func (h *Productos) Crear(c *gin.Context) {
	e, ok := leer(c)
	if !ok {
		return
	}
	p := models.Producto{Nombre: e.Nombre, Precio: *e.Precio}
	if err := h.db.Create(&p).Error; err != nil {
		fallo(c, 500, "No se pudo crear")
		return
	}
	c.JSON(201, p)
}

func (h *Productos) Listar(c *gin.Context) {
	productos := []models.Producto{}
	if err := h.db.Order("id").Find(&productos).Error; err != nil {
		fallo(c, 500, "No se pudo listar")
		return
	}
	c.JSON(200, productos)
}

func (h *Productos) Obtener(c *gin.Context) {
	p, ok := h.buscar(c)
	if ok {
		c.JSON(200, p)
	}
}

func (h *Productos) Actualizar(c *gin.Context) {
	p, ok := h.buscar(c)
	if !ok {
		return
	}
	e, ok := leer(c)
	if !ok {
		return
	}
	p.Nombre, p.Precio = e.Nombre, *e.Precio
	if err := h.db.Model(&p).Select("Nombre", "Precio").Updates(p).Error; err != nil {
		fallo(c, 500, "No se pudo actualizar")
		return
	}
	c.JSON(200, p)
}

func (h *Productos) Eliminar(c *gin.Context) {
	p, ok := h.buscar(c)
	if !ok {
		return
	}
	if err := h.db.Delete(&p).Error; err != nil {
		fallo(c, 500, "No se pudo eliminar")
		return
	}
	c.Status(http.StatusNoContent)
}

// leer convierte el cuerpo JSON en una Entrada y la valida.
func leer(c *gin.Context) (models.Entrada, bool) {
	var e models.Entrada
	if err := c.ShouldBindJSON(&e); err != nil {
		fallo(c, 400, "JSON invalido")
		return e, false
	}
	if err := e.Validar(); err != nil {
		fallo(c, 400, err.Error())
		return e, false
	}
	return e, true
}

// buscar lee el ID de la ruta y consulta el producto por clave primaria.
func (h *Productos) buscar(c *gin.Context) (models.Producto, bool) {
	var p models.Producto
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		fallo(c, 400, "ID invalido")
		return p, false
	}
	err = h.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fallo(c, 404, "Producto no encontrado")
		return p, false
	}
	if err != nil {
		fallo(c, 500, "Error de base de datos")
		return p, false
	}
	return p, true
}
