package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/kevinburke/twilio-go"
	"log"
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
	client = twilio.NewClient(sid, token, nil)
}

func main() {
	msg := "Sale:\n"
	send := false

	for _, w := range Watches {
		if resp, err := soup.Get(w.Url); err != nil {
			log.Fatal(err)
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
		log.Fatal("Nothing's on sale.")
	}

	// Send a message via Twilio
	if _, err := client.Messages.SendMessage(twilioNum, myNum, msg, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sent Successfully to", myNum)
}
