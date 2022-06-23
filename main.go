package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	flags "github.com/jessevdk/go-flags"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
)

type opts struct {
	DiscordToken     string `           long:"token"      env:"DISCORD_TOKEN"  description:"Discord Bot token" mandatory:"true"`
	DiscordChannelId string `           long:"id"      env:"DISCORD_CHANNEL"  description:"Discord Channel ID" mandatory:"true"`
	AlertDay         string `           long:"day"      env:"DAY"  description:"Alerting Day" mandatory:"true"`
	AlertHour        string `           long:"hour"      env:"HOUR"  description:"Alerting Hour" mandatory:"true"`
	AlertMinute      string `           long:"minute"      env:"MINUTE"  description:"Alerting Minute" mandatory:"true"`
}

var (
	argparser *flags.Parser
	arg       opts
)

func initArgparser() {
	argparser = flags.NewParser(&arg, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func main() {
	initArgparser()
	bot, err := discordgo.New("Bot " + arg.DiscordToken)
	if err != nil {
		log.Panicf(err.Error())
		return
	}

	go func() {
		c := cron.New()
		scheduledTime := fmt.Sprintf("0-59/20 %v %v %v * *", arg.AlertMinute, arg.AlertHour, arg.AlertDay)
		c.AddFunc(scheduledTime, func() {
			messageSend, err := bot.ChannelMessageSend(arg.DiscordChannelId, "Reminder: Book your Timesheets")
			if err != nil {
				log.Panicf(err.Error())
				return
			}
			log.Infof("%v: ==> %v", messageSend.EditedTimestamp, messageSend.Content)
		})
		c.Start()
	}()
	select {}
}
