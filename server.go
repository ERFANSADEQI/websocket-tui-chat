package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var wsupgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Subscription struct {
	Conn  *websocket.Conn
	Topic string
}

var subscriptions = make([]*Subscription, 0)

func wsHandler(w http.ResponseWriter, r *http.Request, topic string) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer conn.Close()

	subscription := &Subscription{Conn: conn, Topic: topic}
	subscriptions = append(subscriptions, subscription)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		publish(topic, msg)
	}
}

func publish(topic string, msg []byte) {
	for _, sub := range subscriptions {
		if sub.Topic == topic {
			if err := sub.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ws/:topic", func(c *gin.Context) {
		topic := c.Param("topic")
		wsHandler(c.Writer, c.Request, topic)
	})

	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Run(":8889")
}
