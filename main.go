package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"github.com/martinlindhe/notify"
)

type Msg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type GetChHistoryMsg struct {
	Type string `json:"type"`
	Data struct {
		ChannelID string `json:"channel_id"`
		From      string `json:"from"`
	} `json:"data"`
}

func main() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\nGoodGame: goodbye!")
		os.Exit(0)
	}()

	speechRus := htgotts.Speech{Folder: ".", Language: voices.Russian, Handler: &handlers.Native{}}
	speechEng := htgotts.Speech{Folder: ".", Language: voices.English, Handler: &handlers.Native{}}
	ggcrID := os.Getenv("GGCR_ID")
	if ggcrID == "" {
		panic(errors.New("cannot get your GGCR_ID from env..."))
	}

	conn, _, err := websocket.DefaultDialer.Dial(
		"wss://chat-1.goodgame.ru/chat2/", nil)
	if err != nil {
		panic(err)
	}

	var welcomeMsg Msg
	conn.ReadJSON(&welcomeMsg)
	defer conn.Close()

	if welcomeMsg.Type != "welcome" {
		panic(errors.New("cannot connect Goodgame"))
	}
	fmt.Println("Goodgame: welcome!")

	for {
		lmid, err := os.ReadFile("lmid.txt")
		if err != nil {
			lmid = []byte("0")
		}
		conn.WriteJSON(GetChHistoryMsg{"get_channel_history", struct {
			ChannelID string `json:"channel_id"`
			From      string `json:"from"`
		}{ggcrID, string(lmid)}})

		var msg Msg
		conn.ReadJSON(&msg)
		switch msg.Type {
		case "error":
			var err struct {
				Data struct {
					ErrorMsg string `json:"errorMsg"`
				} `json:"data"`
			}
			json.Unmarshal(msg.Data, &err.Data)
			fmt.Fprintln(os.Stderr, "error: ", err.Data.ErrorMsg)
		case "channel_history":
			var sj struct {
				Data struct {
					Messages []struct {
						UserName  string `json:"user_name"`
						MessageID int    `json:"message_id"`
						Text      string `json:"text"`
					} `json:"messages"`
				} `json:"data"`
			}
			json.Unmarshal(msg.Data, &sj.Data)
			for i := len(sj.Data.Messages) - 1; i > 0; i-- {
				m := sj.Data.Messages[i]
				if strings.Contains(string(lmid), fmt.Sprintf("%d", m.MessageID)) ||
					i < len(sj.Data.Messages)-10 {
					sj.Data.Messages = sj.Data.Messages[i:]
					break
				}
			}
			for _, m := range sj.Data.Messages[1:] {
				t := fmt.Sprintf("%s: %s", m.UserName, m.Text)
				fmt.Println(t)
				notify.Notify("ggcr", "New message!", t, "")
				speechEng.Speak(m.UserName)
				speechRus.Speak(m.Text)
			}
			lm := sj.Data.Messages[len(sj.Data.Messages)-1]
			err = os.WriteFile("lmid.txt", []byte(fmt.Sprintf("%d", lm.MessageID)), 0777)
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second * 10)
	}
}
