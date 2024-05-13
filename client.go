package main

import (
	"fmt"
	"log"
	"strings"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

type model struct {
	topic string 
	messages []string 
	input string
}

type msgReceived string

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
    return nil
}

func connectToWebsocketServer(topic string) *websocket.Conn {
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8889/ws/" + topic, nil)
    if err != nil {
        log.Fatal("خطا در اتصال به سرور وب سوکت: ", err)
    }
    return conn
}

func listenForMessages(m *model, topic string) {
    conn := connectToWebsocketServer(topic)
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("خطا در خواندن پیام: ", err)
            break 
        }
        m.messages = append(m.messages, string(msg))
    }
}

