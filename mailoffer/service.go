package mailoffer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	models "mailinggo/mailer"
	"mailinggo/tools"
	"net/http"
	"net/smtp"
)

// SendEmails envía correos electrónicos a todos los usuarios.
func SendEmails(r *http.Request, w http.ResponseWriter) (*models.MailOffer, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header missing", http.StatusUnauthorized)
        return nil, fmt.Errorf("authorization header missing")
    }
    token := authHeader[len("bearer "):] 
    fmt.Println("Received token:", token)

    ctx := r.Context()

    // Buscar usuario suscritos de la base de datos
    subscribers, err := tools.GetAllSubscribers(ctx)
    if err != nil {
        log.Println("Error fetching users from DB:", err)
        return nil, err
    }
    fmt.Println("Fetched subscribers:", subscribers)

    var enrichedUsers []models.EnrichedUser
    for _, subscriber := range subscribers {
        // Buscar datos faltantes de el servicio de Auth
        userDetails, err := fetchUserDetailsFromAuthService(ctx, subscriber.ID, token)
        if err != nil {
            log.Printf("Error fetching user details for ID %s: %v", subscriber.ID, err)
            continue
        }

        // Solo agregar usuarios suscritos
        if userDetails.Subscribed {
            enrichedUsers = append(enrichedUsers, models.EnrichedUser{
                ID:        userDetails.ID.Hex(), 
                Name:      userDetails.Name,
                Email:     userDetails.Email,
                Subscribed: userDetails.Subscribed,
            })
        }
    }
    fmt.Println("Enriched users:", enrichedUsers)

    // traer todas las ofertas de la base de datos
    articles, err := tools.GetAllArticleOffers(ctx)
    if err != nil {
        log.Println("Error fetching articles from DB:", err)
        return nil, err
    }

    var enrichedArticles []models.EnrichedArticle
    for _, article := range articles {
        if article.Offer {  
            // traer datos faltantes de el servicio cataloggo
            articleDetails, err := fetchArticleDetailsFromCatalogService(ctx, article.ID, token)
            if err != nil {
                log.Printf("Error fetching article details for ID %s: %v", article.ID, err)
                continue
            }

            enrichedArticles = append(enrichedArticles, models.EnrichedArticle{
                Name:        articleDetails.Name,        
                Description: articleDetails.Description, 
                Price:       articleDetails.Price,       
            })
        }
    }
    fmt.Println("Enriched articles:", enrichedArticles)

    // Prepare the MailOffer struct
    mailOffer := &models.MailOffer{
        Articles: enrichedArticles,
        Users:    enrichedUsers,
    }

    // Send the offer emails
    if err := sendOfferEmailsToUsers(mailOffer); err != nil {
        log.Println("Error sending emails:", err)
        http.Error(w, "Failed to send emails", http.StatusInternalServerError)
        return nil, err
    }

    fmt.Println("Emails sent successfully!")
    return mailOffer, nil
}



// fetchUserDetailsFromAuthService llama al servicio Auth para obtener detalles del usuario
func fetchUserDetailsFromAuthService(ctx context.Context, userID string, token string) (*models.UserDetails, error) {
    url := fmt.Sprintf("http://localhost:3000/v1/user/%s", userID)
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "bearer "+token)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
    }

    var userDetails models.UserDetails
    if err := json.NewDecoder(resp.Body).Decode(&userDetails); err != nil {
        return nil, err
    }

    return &userDetails, nil
}

// fetchArticleDetailsFromCatalogService llama al servicio Cataloggo para obtener detalles de los articulos
func fetchArticleDetailsFromCatalogService(ctx context.Context, articleID string, token string) (*models.EnrichedArticle, error) {
    url := fmt.Sprintf("http://localhost:3002/v1/articles/%s", articleID)
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "bearer "+token)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
    }

    var enrichedArticle models.EnrichedArticle
    if err := json.NewDecoder(resp.Body).Decode(&enrichedArticle); err != nil {
        return nil, err
    }

    return &enrichedArticle, nil
}




// sendOfferEmailsToUsers envía los correos electrónicos con las ofertas a los usuarios.
func sendOfferEmailsToUsers(mailOffer *models.MailOffer) error {
	from := "joacomateocortez@gmail.com"
	password := "nxdw bukh weno nutr"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	for _, user := range mailOffer.Users {
		// Crear el contenido del correo electrónico
		subject := "Ofertas Especiales para Ti!"
		body := fmt.Sprintf("Hola %s, aquí tienes una lista de artículos en oferta:\n", user.Name)

		// Recorremos los artículos y añadimos sus detalles al cuerpo del correo
		for _, article := range mailOffer.Articles {
			
			body += fmt.Sprintf("Nombre Articulo: %s\nDescripción: %s\nPrecio: %f\n\n", article.Name, article.Description, article.Price)
		}

		body += "No te pierdas estas increíbles ofertas!"

		message := []byte("Subject: " + subject + "\r\n\r\n" + body)

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{user.Email}, message)
		if err != nil {
			log.Printf("Error sending email to %s: %v", user.Email, err)
			return err
		}

		fmt.Printf("Email sent to %s\n", user.Email)
	}

	return nil
}


func SendCartNotification(r *http.Request, w http.ResponseWriter) error {

    token := r.Header.Get("Authorization")
    if token == "" {
        return fmt.Errorf("authorization header missing")
    }

    ctx := r.Context()

    users, err := tools.GetAllSubscribers(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch users: %v", err)
    }

    // Iterar sobre cada usuario
    for _, user := range users {
        userID := user.ID 

        // Validar si el usuario tiene un carrito abierto (llamada al servicio Cartgo)
        isCartNotEmpty, err := fetchUserCart(userID, token)
        if err != nil {
            log.Printf("Failed to fetch cart for user %s: %v", userID, err)
            continue
        }

        // Continuar solo si el carrito está abierto y no vacío
        if !isCartNotEmpty {
            continue
        }

        // completar datos del usuario si los necesitamos
        enrichedUser, err := fetchUserDetailsFromAuthService(ctx, userID, token)
        if err != nil {
            log.Printf("Failed to fetch enriched user details for %s: %v", userID, err)
            continue
        }

        cartOpenUser := CartOpen{
            Email: enrichedUser.Email, 
            Name:  enrichedUser.Name,  
        }

        if err := SendCartOpenEmail(cartOpenUser); err != nil {
            log.Printf("Failed to send email to %s: %v", cartOpenUser.Email, err)
            continue
        }

        log.Printf("Notification sent successfully to %s", cartOpenUser.Email)
    }

    return nil
}



func fetchUserCart(userID, token string) (bool, error) {
    url := fmt.Sprintf("http://localhost:3003/v1/cart/%s/not-empty", userID)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false, err
    }
    req.Header.Set("Authorization", token)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("failed to fetch cart, status code: %d", resp.StatusCode)
    }

    // Decode the response body into a map
    var response map[string]bool
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return false, fmt.Errorf("failed to decode response: %v", err)
    }

    // Extract the 'cartNotEmpty' field from the decoded response
    isNotEmpty, ok := response["cartNotEmpty"]
    if !ok {
        return false, fmt.Errorf("'cartNotEmpty' field not found in the response")
    }

    fmt.Printf("Cart for user %s is not empty: %v\n", userID, isNotEmpty)

    return isNotEmpty, nil
}




func SendCartOpenEmail(user CartOpen) error {
    tmpl, err := template.ParseFiles("templates/cart_open.tmpl")
    if err != nil {
        return err
    }

    var body bytes.Buffer
    if err := tmpl.Execute(&body, user); err != nil {
        return err
    }

    from := "joacomateocortez@gmail.com"
    password := "nxdw bukh weno nutr"
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    auth := smtp.PlainAuth("", from, password, smtpHost)

    subject := "Unfinished business"
    message := []byte("Subject: " + subject + "\r\n\r\n" + body.String())
    to := []string{user.Email}
    // Send the email
    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
    if err != nil {
        log.Printf("Failed to send email to %s: %v", user.Name, err)
        return err
    }

    log.Printf("Email sent to %s\n", user.Name)
    return nil
}

func SendPriceNotification(updatedArticle string, token string) error {
	client := &http.Client{}

    //Aca se busca todos los carritos que contengan el articulo recibido
	url := fmt.Sprintf("http://localhost:3003/v1/updated-carts?updatedArticle=%s", updatedArticle)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch updated carts: status code %d", resp.StatusCode)
	}

	var carts []Cart
	if err := json.NewDecoder(resp.Body).Decode(&carts); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	users, err := tools.GetAllSubscribers(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching subscribers: %w", err)
	}

	// Crear un mapa de usuarios para acceso rápido por ID
	userMap := make(map[string]models.User)
	for _, user := range users {
		userMap[user.ID] = user 
    }

	// Para cada carrito, enviar un correo si el usuario está suscrito
	for _, cart := range carts {
		// Obtener detalles del usuario desde el servicio Auth
		user, exists := userMap[cart.UserID]
		if !exists {
			log.Printf("Skipping cart %s: user not found", cart.UserID)
			continue
		}

		// Verificar si el usuario está habilitado para recibir notificaciones
		if !user.Subscribed {
			log.Printf("Skipping cart %s: user is not enabled for notifications", cart.UserID)
			continue
		}

		// Obtener detalles completos del usuario desde el servicio Auth
		userDetails, err := fetchUserDetailsFromAuthService(context.Background(), cart.UserID, token)
		if err != nil {
			log.Printf("Failed to fetch user details for %s: %v", cart.UserID, err)
			continue
		}

		cartOpenUser := CartOpen{
			Email: userDetails.Email, 
			Name:  userDetails.Name,   
		}

		// Enviar el correo de cambio de precio
		if err := SendPriceChangeEmail(cartOpenUser); err != nil {
			log.Printf("Failed to send price change email to %s: %v", cartOpenUser.Email, err)
		}
	}

	return nil
}





func SendPriceChangeEmail(user CartOpen) error {

    tmpl, err := template.ParseFiles("templates/price_update.tmpl")
    if err != nil {
        return err
    }

    var body bytes.Buffer
    if err := tmpl.Execute(&body, user); err != nil {
        return err
    }

	from := "joacomateocortez@gmail.com" 
	password := "nxdw bukh weno nutr" 
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

    auth := smtp.PlainAuth("", from, password, smtpHost)

    subject := "Price Changes in Your Cart"
    message := []byte("Subject: " + subject + "\r\n\r\n" + body.String())
    to := []string{user.Email}

    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
    if err != nil {
        log.Printf("Failed to send email to %s: %v", user.Name, err)
        return err
    }

    log.Printf("Price change email sent to %s\n", user.Name)
    return nil
}