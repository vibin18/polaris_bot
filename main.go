package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	flags "github.com/jessevdk/go-flags"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type opts struct {
	DiscordToken     string `           long:"token"      env:"DISCORD_TOKEN"  description:"Discord Bot token" required:"true"`
	DiscordChannelId string `           long:"id"      env:"DISCORD_CHANNEL"  description:"Discord Channel ID" required:"true"`
	AlertDay         string `           long:"day"      env:"DAY"  description:"Alerting Day" required:"true"`
	AlertHour        string `           long:"hour"      env:"HOUR"  description:"Alerting Hour" required:"true"`
	AlertMinute      string `           long:"minute"      env:"MINUTE"  description:"Alerting Minute" required:"true"`
	AlertTimeZone    string `           long:"timezone"      env:"TIMEZONE"  description:"Alerting Timezone" default:"Europe/Berlin" required:"true"`
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
		timeZone, err := time.LoadLocation(arg.AlertTimeZone)
		if err != nil {
			log.Panicf(err.Error())
			return
		}

		c := cron.NewWithLocation(timeZone)
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
