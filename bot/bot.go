package bot

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	// "time"
	// "encoding/binary"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/hraban/opus.v2"
)

var (
	BotToken string
)

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		panic(err)
	}

	// dgv, err := discord.ChannelVoiceJoin(*GuildID, *ChannelID, false, false)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	discord.AddHandler(newMessage)

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "что"):
		channel, _ := discord.Channel(message.ChannelID)
		fmt.Print(channel.Type)
	case strings.Contains(message.Content, "сюда"):
		voiceJoin(discord, message)
	case strings.Contains(message.Content, "привет"):
		discord.ChannelMessageSend(message.ChannelID, "Шалом")
	}

}

func voiceJoin(discord *discordgo.Session, message *discordgo.MessageCreate) {
	channel, _ := discord.Channel(message.ChannelID)
	if channel.Type != 2 {
		discord.ChannelMessageSend(message.ChannelID, "Это команду нужно писать в текстовый чат голосового канала")
		return
	}
	voiceDs, err := discord.ChannelVoiceJoin(message.GuildID, message.ChannelID, false, false)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Open("z.opus")
	if err != nil {
		fmt.Println(err)
	}
	s, err := opus.NewStream(f)
	if err != nil {
		fmt.Println(err)
	}
	defer s.Close()
	pcmbuf := make([]int16, 16384)
	voice := make([]byte, 16384)
	for key, value := range pcmbuf {
		voice[key] = byte(value)
	}
	for {
		n, err := s.Read(pcmbuf)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		voiceDs.OpusSend <- voice[:n*2]

		// send pcm to audio device here, or write to a .wav file

	}
}
