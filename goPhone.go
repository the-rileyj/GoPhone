package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"

	twilio "github.com/saintpete/twilio-go"
)

//Struct to hold the various key data
type datAuth struct {
	Number string `json:"number"`
	Pass   string `json:"pass"`
	Sid    string `json:"sid"`
	Token  string `json:"token"`
}

type phoneMessage struct {
	Message string
	Type    string
	Number  string
	Pass    string
}

var tpl *template.Template

func main() {
	fullPath := strings.Join(os.Args[1:], "") //Join together the path provided by cmd args index 1 to end
	go fmt.Printf("Reading %s\n", fullPath)   //Print out the path for user

	//Reading the data in from the location specified by command line argument
	var dat datAuth
	bdata, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Fatal("Could not read data properly", err)
	}

	//Unmarshaling the data into a datAuth struct to hold sensative information
	if json.Unmarshal(bdata, &dat) != nil {
		log.Fatal("Error Unmarshalling the data")
	}

	//Get the ID, and then get the lyrics with the returning ID from that function
	//lyrics := getSongLyrics(getSongID(artist, title, dat.Lkey), dat.Lkey)
	client, m, r, tpl := twilio.NewClient(dat.Sid, dat.Token, nil), melody.New(), gin.Default(), template.Must(template.New("").ParseGlob("*.gohtml"))
	//tpl := template.Must(template.New("").ParseGlob("*.gohtml"))
	//r, m := gin.Default(), melody.New()
	//m := melody.New()
	r.GET("/phone", func(c *gin.Context) {
		tpl.ExecuteTemplate(c.Writer, "index.gohtml", nil)
	})

	r.GET("/ws-phone", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})
	r.Run(":5000")
	/*if u, err := url.Parse("http://therileyjohnson.com/public/files/mp3/freshmanEdit.mp3"); err == nil {
		jack := "+17013186329"
		//matt := "+17014467380"
		client.Calls.MakeCall(dat.Number, jack, u)
	}
	//Split the lyrics up by the sections of lyrics seperated by two newlines
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
