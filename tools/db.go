package tools

import (
	"context"
	"errors"
	"log"
	"time"

	models "mailinggo/mailer" // Importar los modelos desde el paquete correspondiente
	"mailinggo/tools/env"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func ConnectToMongo(mongoURL string) (*mongo.Client, context.Context, context.CancelFunc) {
    clientOptions := options.Client().ApplyURI(mongoURL)

    // Crear cliente
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatalf("Error al crear el cliente de MongoDB: %v", err)
    }

    // Contexto con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    // Conectar al servidor
    if err := client.Connect(ctx); err != nil {
        cancel()
        log.Fatalf("Error al conectar con MongoDB: %v", err)
    }

    return client, ctx, cancel
}

// Get obtiene la instancia de la base de datos
func Get(ctx context.Context) (*mongo.Database, error) {
	if database == nil {
		clientOptions := options.Client().ApplyURI(env.Get().MongoURL)

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("Error connecting to MongoDB: %v", err)
			return nil, err
		}

		if err = client.Ping(ctx, nil); err != nil {
			log.Printf("Error pinging MongoDB: %v", err)
			return nil, err
		}

		database = client.Database("mailinggo")
	}
	return database, nil
}

func InsertUserSubscriber(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	log.Println("\n dentro del insert")

	db, err := Get(ctx)
	if err != nil {
		log.Printf("Error al obtener la base de datos: %v", err)
		return nil, err
	}
	log.Println("Base de datos obtenida correctamente", db.Name())
	return db.Collection("users").InsertOne(ctx, user)

}

func GetAllSubscribers(ctx context.Context) ([]models.User, error) {
	db, err := Get(ctx)
	if err != nil {
		return nil, err
	}
	cursor, err := db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUserSubscriber(ctx context.Context, userID string) error {
	db, err := Get(ctx)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": userID}
	result, err := db.Collection("users").DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func UpdateUserState(ctx context.Context, userID string, subscribed bool) error {
	db, err := Get(ctx)
	if err != nil {
		return err
	}

	// Usamos el campo "subscribed" en lugar de "enabled"
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"subscribed": subscribed}}

	// Actualización de usuario en la colección "users"
	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("user not found") // Si no se encuentra el usuario, devolvemos un error
	}
	return nil
}

func FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	db, err := Get(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = db.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func InsertArticleOffer(ctx context.Context, article *models.Article) (*mongo.InsertOneResult, error) {
	db, err := Get(ctx)
	if err != nil {
		return nil, err
	}
	return db.Collection("articles_offers").InsertOne(ctx, article)
}

func GetAllArticleOffers(ctx context.Context) ([]models.Article, error) {
	db, err := Get(ctx)
	if err != nil {
		return nil, err
	}
	cursor, err := db.Collection("articles_offers").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	return articles, nil
}

func DeleteArticleOffer(ctx context.Context, articleID string) error {
	db, err := Get(ctx)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": articleID}
	result, err := db.Collection("articles_offers").DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		return errors.New("article not found")
	}
	return nil
}

func UpdateArticleState(ctx context.Context, articleID string, enabled bool) error {
	log.Printf("[UpdateArticleState] Intentando actualizar artículo con ID: %s a estado Offer: %t", articleID, enabled)
	db, err := Get(ctx)
	if err != nil {
		log.Printf("[UpdateArticleState] Error al obtener la base de datos: %v", err)
		return err
	}
	filter := bson.M{"_id": articleID}
	update := bson.M{"$set": bson.M{"offer": enabled}}

	log.Printf("[UpdateArticleState] Filtro: %+v", filter)
	log.Printf("[UpdateArticleState] Actualización: %+v", update)

	result, err := db.Collection("articles_offers").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("[UpdateArticleState] Error al actualizar el artículo: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Printf("[UpdateArticleState] No se encontró ningún artículo para actualizar con ID: %s", articleID)
		return errors.New("article not found")
	}

	log.Printf("[UpdateArticleState] Estado del artículo actualizado exitosamente.")
	return nil
}

func FindArticleByID(ctx context.Context, articleID string) (*models.Article, error) {
	log.Printf("[FindArticleByID] Buscando artículo con ID: %s", articleID)
	db, err := Get(ctx)
	if err != nil {
		log.Printf("[FindArticleByID] Error al obtener la base de datos: %v", err)
		return nil, err
	}

	var article models.Article
	err = db.Collection("articles_offers").FindOne(ctx, bson.M{"_id": articleID}).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[FindArticleByID] No se encontró ningún documento con ID: %s", articleID)
			return nil, nil
		}
		log.Printf("[FindArticleByID] Error al decodificar el documento: %v", err)
		return nil, err
	}

	log.Printf("[FindArticleByID] Artículo encontrado: %+v", article)
	return &article, nil
}
