package models

import (
	"errors"
	"strings"
)

// Producto representa una fila de la tabla productos.
type Producto struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre"`
	Precio int64  `json:"precio"` // Pesos enteros, por ejemplo: 25000.
}

// Entrada define los campos que el cliente puede enviar en POST y PUT.
type Entrada struct {
	Nombre string `json:"nombre"`
	Precio *int64 `json:"precio"` // Puntero para detectar un campo omitido.
}

var ErrEntradaInvalida = errors.New("Envia nombre y precio entero mayor o igual a 0")

// Validar limpia el nombre y comprueba que los campos obligatorios sean validos.
func (e *Entrada) Validar() error {
	e.Nombre = strings.TrimSpace(e.Nombre)
	if e.Nombre == "" || e.Precio == nil || *e.Precio < 0 {
		return ErrEntradaInvalida
	}
	return nil
}
