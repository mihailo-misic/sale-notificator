package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/kevinburke/twilio-go"
	"github.com/mihailo-misic/sale-notificator/env"
	"os"
	"strings"
)

var client *twilio.Client

type Watch struct {
	Name    string
	Url     string
	Element []string
	Price   string
}

var Watches = []Watch{
	{
		"GANT DENVILLE",
		"https://goo.gl/VEFdCv",
		[]string{"span", "class", "price_big"},
		"14.990,00RSD",
	},
	{
		"PIERRE LANNIER ELEGANCE",
		"https://goo.gl/pQa87s",
		[]string{"span", "class", "price_big"},
		"18.990,00RSD",
	},
}

func init() {
	client = twilio.NewClient(env.Sid, env.Token, nil)
}

func main() {
	msg := "Sale:\n"
	send := false

	for _, w := range Watches {
		if resp, err := soup.Get(w.Url); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			doc := soup.HTMLParse(resp)
			price := doc.Find(w.Element...).Text()
			if price != w.Price {
				msg += fmt.Sprintf("%s: %s (%s)\n", w.Name, price, w.Url)
				send = true
			}
		}
	}

	msg = strings.TrimRight(msg, "\n")

	// Don't send if non are on sale.
	if !send {
		fmt.Println("Nothing's on sale.")
		os.Exit(1)
	}

	// Send a message via Twilio
	if _, err := client.Messages.SendMessage(env.TwilioNum, env.SendNum, msg, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Sent Successfully to", env.SendNum)
}
