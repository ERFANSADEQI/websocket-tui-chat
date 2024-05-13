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

func sendMessage(input string, topic string) {
    conn := connectToWebsocketServer(topic)
    defer conn.Close()

    if err := conn.WriteMessage(websocket.TextMessage, []byte(input)); err != nil {
        log.Println("خطا در ارسال پیام: ", err)
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			sendMessage(m.input, m.topic)
			m.input = ""
		default:
			m.input += string(msg.Runes)
		}
	case msgReceived:
		m.messages = append(m.messages, string(msg))
	}
	return m, nil
}

