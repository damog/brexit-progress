package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context) (string, error) {
	loc, _ := time.LoadLocation("Europe/London")

	now := time.Now().In(loc)
	start := time.Date(2017, time.March, 29, 12, 0, 0, 0, loc)
	end := time.Date(2019, time.October, 31, 23, 0, 0, 0, loc)

	duration_end := end.Sub(start)
	duration_now := now.Sub(start)
	perc := (duration_now.Seconds() * 100) / duration_end.Seconds()

	if !(perc < 100) {
		// nothing to do
		// Brexit is here!
		os.Exit(0)
	}

	duration_until := end.Sub(now)

	hours, minutes, seconds := int(duration_until.Hours()), int(duration_until.Minutes()), int(duration_until.Seconds())

	days := int(hours / 24)

	hours = hours - (days * 24)
	minutes = minutes - (days * 1440) - (hours * 60)
	seconds = seconds - (days * 86400) - (hours * 3600) - (minutes * 60)

	eu_emoji, wave_emoji, uk_emoji := "🇪🇺", "👋", "🇬🇧"

	tweet := fmt.Sprintf("Time until 3rd #Brexit deadline on 31 Oct 2019: %dd %dh %dm %ds %s%s%s\n", days, hours, minutes, seconds, eu_emoji, wave_emoji, uk_emoji)

	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}
	client, err := getClient(&creds)

	tw, resp, err := client.Statuses.Update(tweet, nil)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v\n", resp)
	log.Printf("%+v\n", tw)

	return fmt.Sprintf("Am I done?"), nil

}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client, nil
}
