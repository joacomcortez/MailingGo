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

func HandleSubscribedUser(c *gin.Context) {
	log.Println("Entrando al manejador de suscripción...")

	// Cuerpo solicitud
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error leyendo el cuerpo de la solicitud: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	log.Printf("Cuerpo de la solicitud: %s\n", body)

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	var userInput struct {
		ID         string `json:"_id" binding:"required"`         
		Subscribed bool   `json:"subscribed"` 
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		log.Printf("Error en la validación de datos: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Log de los datos recibidos
	log.Printf("Datos recibidos - ID: %s, Subscribed: %t\n", userInput.ID, userInput.Subscribed)

	ctx := context.TODO()

	log.Println("Buscando usuario en la base de datos...")
	user, err := tools.FindUserByID(ctx, userInput.ID)
	if err != nil {
		log.Printf("Error buscando usuario: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	if user == nil {
		log.Println("Usuario no encontrado. Creando uno nuevo...")
		// Si el usuario no existe, crear uno nuevo con el estado recibido
		newUser := models.User{
			ID:         userInput.ID,
			Subscribed: userInput.Subscribed,
		}

		log.Printf("Datos del nuevo usuario: %+v\n", newUser)

		_, err := tools.InsertUserSubscriber(ctx, &newUser)
		if err != nil {
			log.Printf("Error al insertar el usuario: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
			return
		}

		log.Println("Usuario creado exitosamente.")
		c.JSON(http.StatusOK, gin.H{
			"message": "User successfully created",
			"id":      newUser.ID,
			"state":   newUser.Subscribed,
		})
		return
	}

	// Usuario encontrado. Verificar si el estado debe actualizarse
	log.Println("Usuario encontrado. Verificando si el estado ha cambiado...")
	if user.Subscribed != userInput.Subscribed {
		log.Printf("Actualizando estado del usuario (ID: %s) de %t a %t\n", userInput.ID, user.Subscribed, userInput.Subscribed)
		err := tools.UpdateUserState(ctx, userInput.ID, userInput.Subscribed)
		if err != nil {
			log.Printf("Error al actualizar el estado del usuario: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user state"})
			return
		}

		state := "subscribed"
		if !userInput.Subscribed {
			state = "unsubscribed"
		}

		log.Println("Estado del usuario actualizado exitosamente.")
		c.JSON(http.StatusOK, gin.H{
			"message": "User state successfully updated",
			"state":   state,
		})
		return
	}

	// No se realizaron cambios
	log.Println("No se realizaron cambios en el estado del usuario.")
	c.JSON(http.StatusOK, gin.H{"message": "No changes were made"})
}
