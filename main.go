package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	tok string
)

func init() {
	flag.StringVar(&tok, "t", "", "<user token>")
	flag.Parse()

	if tok == "" {
		tok = os.Getenv("DISCORD_TOKEN")
	}
}

func main() {
	log.Println("Starting Texcord...")

	discord, err := discordgo.New(tok)
	if err != nil {
		log.Panic(err)
		return
	}

	discord.AddHandler(messageEvent)

	err = discord.Open()
	if err != nil {
		log.Println("Error opening connection,", err)
		return
	}

	log.Println("Bot started")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != s.State.User.ID || !strings.HasPrefix(m.Message.Content, "TEX:") {
		return
	}

	tex := strings.Replace(m.Message.Content, "TEX:", "", 1)
	tex = strings.Replace(tex, "```", "", -1)
	tex = strings.Replace(tex, "\n", "", -1)

	u := "https://chart.googleapis.com/chart?cht=tx&chs=100&chl=" + url.QueryEscape(tex)

	var c string

	msg := &discordgo.MessageEdit{
		ID:      m.ID,
		Content: &c,
		Channel: m.ChannelID,
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				Height: 100,
				URL:    u,
			},
		},
	}
	s.ChannelMessageEditComplex(msg)
}
