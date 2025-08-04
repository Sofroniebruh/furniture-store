package main

import (
	"auth-service/db"
	"auth-service/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type body struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

func main() {
	err := db.Init()
	//conn, ch, err := config.InitRabbitMq()
	//
	//_, _ = config.DeclareQueue(
	//	ch,
	//	"codeQueue",
	//	false,
	//	false,
	//	false,
	//	false)
	//
	//if err != nil {
	//	log.Fatal("Failed to initialize RabbitMQ: ", err)
	//}
	//
	//defer conn.Close()
	//defer ch.Close()
	//
	//go func() {
	//	callback := func(response config.ConsumedMessage[body]) error {
	//		log.Println("Received a message: ", response)
	//
	//		var errorMessage struct {
	//			StatusCode int    `json:"status_code"`
	//			Message    string `json:"message"`
	//		}
	//
	//		emailResponse, err := handlers.SendEmail(response.Body.Code, response.Body.Email)
	//
	//		if err != nil {
	//			log.Println("Failed to send email: ", err)
	//			errorMessage.StatusCode = 500
	//			errorMessage.Message = err.Error()
	//			errorData, _ := json.Marshal(errorMessage)
	//			err = ch.Publish(
	//				"",
	//				"replyCodeQueue",
	//				false,
	//				false,
	//				amqp091.Publishing{
	//					ContentType:   "application/json",
	//					Body:          errorData,
	//					CorrelationId: response.CorrelationId.String(),
	//				})
	//
	//			if err != nil {
	//				log.Println("Failed to send error response: ", err)
	//				return err
	//			}
	//			return err
	//		}
	//
	//		data, err := json.Marshal(emailResponse)
	//
	//		err = ch.Publish(
	//			"",
	//			"replyCodeQueue",
	//			false,
	//			false,
	//			amqp091.Publishing{
	//				ContentType:   "application/json",
	//				Body:          data,
	//				CorrelationId: response.CorrelationId.String(),
	//			})
	//
	//		if err != nil {
	//			log.Println("Failed to send email: ", err)
	//			return err
	//		}
	//
	//		return nil
	//	}
	//
	//	err = config.ConsumeMessage(ch, "codeQueue", callback)
	//
	//	if err != nil {
	//		log.Fatal("Failed to consume message: ", err)
	//	}
	//}()

	if err != nil {
		log.Fatal("Error initializing database: ", err)
		return
	}

	r := chi.NewRouter()

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	r.Post("/registration", handlers.Signup)
	r.Post("/login", handlers.Login)
	r.Post("/logout", handlers.Logout)
	r.Post("/refresh", handlers.Refresh)
	r.Post("/send-code", handlers.GenerateCode)

	log.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", handler)

	if err != nil {
		log.Fatal("Failed listening on port 8080: ", err)
		return
	}
}
