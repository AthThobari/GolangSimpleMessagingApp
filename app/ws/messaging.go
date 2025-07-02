package ws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
)


func ServeWSMessaging(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool) 
	var broadcast = make(chan models.MessagePayload)

	app.Get("message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true

		for {
			var msg models.MessagePayload
			if err := c.ReadJSON(&msg); err != nil {
				fmt.Println("error payload: ", err)
				break
			}
			msg.Date = time.Now()
			err := repository.InsertNewMessage(context.Background(), msg)
			if err != nil {
				fmt.Println(err)
			}
			broadcast <- msg
		}
	}))

	go func() {
		for {
		msg := <- broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("failed to write json: ", err)
				client.Close()
				delete(clients, client)
			}
		}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT_SOCKET", "8080"))))
}