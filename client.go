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



