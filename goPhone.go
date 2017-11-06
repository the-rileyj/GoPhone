package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	twilio "github.com/saintpete/twilio-go"
)

//Struct to hold the various key data
type datAuth struct {
	Sid    string `json:"sid"`
	Token  string `json:"token"`
	Number string `json:"number"`
}

func main() {
	/*var artist, title, keys, span, from, to string //Variables to hold the command line args corresponding to
	flag.StringVar(&artist, "artist", "", "The name of the artist of the song you want to lookup")
	flag.StringVar(&title, "title", "", "The name of the song you want to lookup")
	flag.StringVar(&to, "to", "", "The phone number you're sending to in the format '+(Country Code)(Area Code)(Phone Number)', ex. '+17015559999")
	flag.StringVar(&from, "from", "", "The twilio number you're sending from, if not included, it's assumed that you have it in your keys .json file")
	flag.StringVar(&keys, "keys", "", "The location of the keys for the Twilio and MusixMatch API's, json should look like: ")
	flag.StringVar(&span, "span", "", "The time span over which the lyrics are to be sent every 'span' / 'number of verses' amount of time")
	flag.Parse()*/
	fullPath := strings.Join(os.Args[1:], "")
	go fmt.Printf("Reading %s\n", fullPath)
	//Reading the data in from the location specified by the "keys" argument
	var dat datAuth
	bdata, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Fatal("Could not read data properly", err)
	}

	//Unmarshaling the data into a datAuth struct to hold sensative information
	if json.Unmarshal(bdata, &dat) != nil {
		log.Fatal("Error Unmarshalling the data")
	}

	/*/Making sure that there is a sending Twilio number
	if from == "" && dat.Number != "" {
		from = dat.Number
	} else if from == "" && dat.Number == "" {
		log.Fatal("Need a sending Twilio number")
	}*/

	//Get the ID, and then get the lyrics with the returning ID from that function
	//lyrics := getSongLyrics(getSongID(artist, title, dat.Lkey), dat.Lkey)
	client := twilio.NewClient(dat.Sid, dat.Token, nil)
	if u, err := url.Parse("http://therileyjohnson.com/public/files/mp3/freshmanEdit.mp3"); err == nil {
		jack := "+17013186329"
		//matt := "+17014467380"
		client.Calls.MakeCall(dat.Number, jack, u)
	}
	/*/Split the lyrics up by the sections of lyrics seperated by two newlines
	slyrics := strings.Split(lyrics, "\n\n")

	//Assuring the starting date argument isn't empty and if it is defaulting to sending the lyrics right now
	if span != "" {
		tTime, err := time.ParseDuration(span) //Parse the sending interval and check for success
		if err != nil {
			log.Fatal("Error parsing sending interval")
		}

		//Range over all values except the last which isn't lyrical material,
		//Then between iterations sleep for the set interval
		for _, l := range slyrics[:len(slyrics)-1] {
			time.Sleep(time.Duration(tTime.Nanoseconds() / int64(len(slyrics)-1)))
			_, err = client.Messages.SendMessage(dat.Number, to, l, nil)
		}
	} else {
		//Joining together the split up lyrics, but making sure not to include the final value because it ins't lyrical material
		_, err = client.Messages.SendMessage(dat.Number, to, strings.Join(slyrics[:len(slyrics)-1], "\n\n"), nil)
	}*/
}
