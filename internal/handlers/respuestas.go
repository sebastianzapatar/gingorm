package handlers

import "github.com/gin-gonic/gin"

// fallo responde con un JSON de la forma {"error": "..."}.
func fallo(c *gin.Context, estado int, mensaje string) {
	c.JSON(estado, gin.H{"error": mensaje})
}
