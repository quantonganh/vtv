package main

import (
	"bufio"
	"embed"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed VNese_poems.txt
var poems embed.FS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func shuffle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	for {
		line := getRandomLine()
		words := strings.Split(line, " ")
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := ws.WriteMessage(websocket.TextMessage, []byte(strings.Join(words, " / "))); err != nil {
			log.Println(err)
			return
		}

		_, p, err := ws.ReadMessage()
		if err == nil && string(p) == "\n" {
			if err := ws.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
				log.Println(err)
				return
			}
		} else {
			time.Sleep(2 * time.Second)

			if err := ws.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
				log.Println(err)
				return
			}
		}

		ws.WriteMessage(websocket.TextMessage, []byte("\n"))
	}
}

func getRandomLine() string {
	f, err := poems.Open("VNese_poems.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	lineNum := 1
	var pick string
	for scanner.Scan() {
		line := scanner.Text()
		roll := generator.Intn(lineNum)
		if roll == 0 {
			pick = line
		}
		lineNum += 1
	}
	return pick
}
