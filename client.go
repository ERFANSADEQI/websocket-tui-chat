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
type sendMessage string

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
    return nil
}

