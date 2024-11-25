package rest

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"

	models "mailinggo/mailer"
	"mailinggo/tools"

	"github.com/gin-gonic/gin"
)

// Maneja el estado de oferta de un artículo
func HandleArticleOffer(c *gin.Context) {
	
		// Leer el cuerpo de la solicitud para depurar
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error leyendo el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		log.Printf("Cuerpo de la solicitud: %s\n", body)
	
		// Volver a establecer el cuerpo de la solicitud para que Gin lo pueda leer correctamente después
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	// Estructura para parsear el request
	var articleOffer struct {
		ID    string `json:"_id" binding:"required"`
		Offer bool   `json:"offer" `
	}

	// Validar datos de entrada
	if err := c.ShouldBindJSON(&articleOffer); err != nil {
		log.Printf("[HandleArticleOffer] Error en la validación de datos: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	log.Printf("[HandleArticleOffer] Datos recibidos: %+v", articleOffer)

	ctx := context.TODO()

	// Buscar artículo por ID
	log.Printf("[HandleArticleOffer] Buscando artículo con ID: %s", articleOffer.ID)
	article, err := tools.FindArticleByID(ctx, articleOffer.ID)
	if err != nil {
		log.Printf("[HandleArticleOffer] Error buscando artículo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find article"})
		return
	}
	log.Printf("[HandleArticleOffer] Artículo encontrado: %+v", article)

	// Si el artículo no existe, lo creamos
	if article == nil {
		log.Printf("[HandleArticleOffer] Artículo no encontrado. Creando uno nuevo.")
		newArticle := models.Article{
			ID:    articleOffer.ID,
			Offer: articleOffer.Offer,
		}

		log.Printf("[HandleArticleOffer] Insertando artículo: %+v", newArticle)
		result, err := tools.InsertArticleOffer(ctx, &newArticle)
		if err != nil {
			log.Printf("[HandleArticleOffer] Error al insertar artículo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert article"})
			return
		}
		log.Printf("[HandleArticleOffer] Artículo creado con ID: %v", result.InsertedID)

		c.JSON(http.StatusOK, gin.H{
			"message": "Article successfully created",
			"id":      result.InsertedID,
		})
		return
	}

	// Si el artículo existe, actualizamos su estado si es necesario
	if article.Offer != articleOffer.Offer {
		log.Printf("[HandleArticleOffer] Actualizando estado del artículo con ID: %s", articleOffer.ID)
		err := tools.UpdateArticleState(ctx, articleOffer.ID, articleOffer.Offer)
		if err != nil {
			log.Printf("[HandleArticleOffer] Error al actualizar estado del artículo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article state"})
			return
		}

		state := "offer enabled"
		if !articleOffer.Offer {
			state = "offer disabled"
		}
		log.Printf("[HandleArticleOffer] Estado actualizado: %s", state)

		c.JSON(http.StatusOK, gin.H{
			"message": "Article state successfully updated",
			"state":   state,
		})
		return
	}

	// No hubo cambios en el estado
	log.Printf("[HandleArticleOffer] No se realizaron cambios en el estado del artículo.")
	c.JSON(http.StatusOK, gin.H{"message": "No changes were made"})
}
