package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type model struct {
	topic        string
	messages     []string
	input        string
	conn         *websocket.Conn
	err          error
	userMessages []string
}

type msgReceived []byte

func initialModel(topic string) model {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8889/ws/"+topic, nil)
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

func sendMessage(m *model, input string) {
	if err := m.conn.WriteMessage(websocket.TextMessage, []byte(input)); err != nil {
		log.Println("خطا در ارسال پیام: ", err)
		return
	}
	m.userMessages = append(m.userMessages, input)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			sendMessage(&m, m.input)
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
		formattedMsg := msg
		if contains(m.userMessages, msg) {
			formattedMsg = fmt.Sprintf("%s%s%s", colorGreen, msg, colorReset)
		}
		view += fmt.Sprintf("%s\n", formattedMsg)
	}
	view += fmt.Sprintf("> %s", m.input)
	return view
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

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