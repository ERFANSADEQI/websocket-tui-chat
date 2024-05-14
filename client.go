package main

import (
	"fmt"
	"bufio"
	"log"
	"strings"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

type model struct {
	topic    string
	messages []string
	input    string
	conn     *websocket.Conn
	err      error
}

func initialModel(topic string) model {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8889/ws/" + topic, nil)
	if err != nil {
		log.Fatal("خطا در اتصال به سرور وب سوکت: ", err)
	}
	return model{topic: topic, conn: conn}
}

func (m model) Init() tea.Cmd {
	return listenForMessages(m.conn)
}

func listenForMessages(conn *websocket.Conn) tea.Cmd {
	return func() tea.Msg {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return tea.Quit
			}
			return msgReceived(msg)
		}
	}
}

func sendMessage(conn *websocket.Conn, input string) {
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
			sendMessage(m.conn, m.input)
			m.input = ""
		default:
			m.input += string(msg.Runes)
		}
	case msgReceived:
		m.messages = append(m.messages, string(msg))
	}
	return m, listenForMessages(m.conn)
}

func (m model) View() string {
	view := fmt.Sprintf("Topic: %s\n\n", m.topic)
	for _, msg := range m.messages {
		view += fmt.Sprintf("%s\n", msg)
	}
	view += fmt.Sprintf("> %s", m.input)
	return view
}

type msgReceived []byte

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter topic: ")
	topic, _ := reader.ReadString('\n')
	topic = strings.TrimSpace(topic)

	p := tea.NewProgram(initialModel(topic))
	if err := p.Start(); err != nil {
		log.Fatalf("خطایی رخ داده است: %v", err)
	}
}