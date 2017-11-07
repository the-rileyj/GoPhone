package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/static"
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

	client, m, r, tpl := twilio.NewClient(dat.Sid, dat.Token, nil), melody.New(), gin.Default(), template.Must(template.New("").ParseGlob("*.gohtml"))
	r.GET("/phone", func(c *gin.Context) {
		tpl.ExecuteTemplate(c.Writer, "index.gohtml", nil)
	})

	r.Use(func(c *gin.Context) {
		fmt.Println(c.Request.URL.Path)
	})

	r.GET("/public/:fi", static.Serve("/public", static.LocalFile("public/", true)))

	r.GET("/ws-phone", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var pmsg phoneMessage
		if json.Unmarshal(msg, &pmsg) != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		if pmsg.Pass == dat.Pass {
			switch strings.ToLower(pmsg.Type) {
			case "init":
				if rmsg, err := json.Marshal(phoneMessage{"Connected", "init", "", dat.Pass}); err == nil {
					if err = s.Write(rmsg); err != nil {
						fmt.Printf("Error sending websocket message: %s\n", err)
					}
				} else {
					fmt.Printf("Error Marshalling outgoing message: %s\n", err)
				}
			case "text":
				_, err := client.Messages.SendMessage(dat.Number, pmsg.Number, pmsg.Message, nil)
				if err != nil {
					emsg := fmt.Sprintf("Error sending sms message: %s\n", err)
					fmt.Printf(emsg)
					if rmsg, err := json.Marshal(phoneMessage{emsg, "error", "", dat.Pass}); err == nil {
						if err = s.Write(rmsg); err != nil {
							fmt.Printf("Error sending websocket error message: %s\n", err)
						}
					} else {
						fmt.Printf("Error Marshalling outgoing error message: %s\n", err)
					}
				} else {
					if rmsg, err := json.Marshal(phoneMessage{pmsg.Message, "confirm", pmsg.Number, dat.Pass}); err == nil {
						if err = m.Broadcast(rmsg); err != nil {
							fmt.Printf("Error sending websocket confirm message: %s\n", err)
						}
					} else {
						fmt.Printf("Error Marshalling outgoing confirm message: %s\n", err)
					}
				}
			case "call":
				if pmsg.Message != "" {
					if u, err := url.Parse(pmsg.Message); err == nil {
						if _, err := client.Calls.MakeCall(dat.Number, pmsg.Number, u); err != nil {
							if rmsg, err := json.Marshal(phoneMessage{fmt.Sprintf("Error making outgoing call: %s\n", err), "error", pmsg.Number, dat.Pass}); err == nil {
								if err = m.Broadcast(rmsg); err != nil {
									fmt.Printf("Error sending websocket confirm error making outgoing call: %s\n", err)
								}
							} else {
								fmt.Printf("Error Marshalling outgoing confirm error making outgoing call: %s\n", err)
							}
						} else {
							if rmsg, err := json.Marshal(phoneMessage{pmsg.Message, "confirm", pmsg.Number, dat.Pass}); err == nil {
								if err = m.Broadcast(rmsg); err != nil {
									fmt.Printf("Error sending websocket confirm making outgoing call: %s\n", err)
								}
							} else {
								fmt.Printf("Error Marshalling outgoing confirm making outgoing call: %s\n", err)
							}
						}
					} else {
						if rmsg, err := json.Marshal(phoneMessage{fmt.Sprintf("Error parsing the url for the call: %s\n", err), "error", pmsg.Number, dat.Pass}); err == nil {
							if err = m.Broadcast(rmsg); err != nil {
								fmt.Printf("Error sending websocket confirm message: %s\n", err)
							}
						} else {
							fmt.Printf("Error Marshalling outgoing confirm message: %s\n", err)
						}
					}
				} else {
					if _, err := client.Calls.MakeCall(dat.Number, pmsg.Number, nil); err != nil {
						if rmsg, err := json.Marshal(phoneMessage{fmt.Sprintf("Error making outgoing call: %s\n", err), "error", pmsg.Number, dat.Pass}); err == nil {
							if err = m.Broadcast(rmsg); err != nil {
								fmt.Printf("Error sending websocket confirm error making outgoing call: %s\n", err)
							}
						} else {
							fmt.Printf("Error Marshalling outgoing confirm error making outgoing call: %s\n", err)
						}
					} else {
						if rmsg, err := json.Marshal(phoneMessage{pmsg.Message, "confirm", pmsg.Number, dat.Pass}); err == nil {
							if err = m.Broadcast(rmsg); err != nil {
								fmt.Printf("Error sending websocket confirm making outgoing call: %s\n", err)
							}
						} else {
							fmt.Printf("Error Marshalling outgoing confirm making outgoing call: %s\n", err)
						}
					}
				}
			}
		} else {
			if rmsg, err := json.Marshal(phoneMessage{"Incorrect Password", "init", "", pmsg.Pass}); err == nil {
				if err = s.Write(rmsg); err != nil {
					fmt.Printf("Error sending auth error message: %s\n", err)
				}
			} else {
				fmt.Printf("Error Marshalling auth error message: %s\n", err)
			}
		}
	})
	r.Run(":5000")
}
