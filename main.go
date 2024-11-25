package main

import (
	"log"
	"mailinggo/rabbit/consume"
	"mailinggo/rest"
	db "mailinggo/tools" // Importar tools/db para la conexión a MongoDB
	"mailinggo/tools/env"
	"strconv"
)

func main() {
	// Cargar la configuración
	config := env.Get()

	// Conectar a MongoDB
	client, ctx, cancel := db.ConnectToMongo(config.MongoURL)
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error al desconectar de MongoDB: %v", err)
		} else {
			log.Println("Conexión a MongoDB cerrada correctamente")
		}
	}()

	// Registrar estado de la conexión
	log.Println("Verificando conexión a MongoDB...")
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("No se pudo conectar a MongoDB: %v", err)
	} else {
		log.Println("Conexión exitosa a MongoDB")
	}

	// Inicializar el enrutador de Gin desde el paquete rest
	router := rest.NewRouter()

	port := ":" + strconv.Itoa(config.Port)

	// Iniciar el consumidor RabbitMQ
	go consume.ConsumePriceChangeNotifications() // Consumidor en una goroutine

	// Iniciar el servidor HTTP de Gin
	log.Printf("El servidor se está iniciando en el puerto %s...", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
