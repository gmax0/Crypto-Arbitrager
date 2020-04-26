package websocket

import (
	"errors"
	"net/url"
	"os"
	_ "time"

	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var f os.File

func init() {
	// open a file
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.SetLevel(logrus.DebugLevel)
}

type WSClient struct {
	connection *websocket.Conn
	ID         *uuid.UUID
	subMsg     []byte
	unsubMsg   []byte
}

//
func NewClient(socketUrl string) (*WSClient, error) {
	//Use default dialer for now
	u := url.URL{Scheme: "wss", Host: socketUrl}
	log.Info("Connecting to ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error("Dialer initialization error:", err)
		return nil, err
	}

	//Initiate an UID
	id, err := uuid.NewV4()
	if err != nil {
		log.Error("Error generating V4 UUID:", err)
	}

	client := &WSClient{c, id, nil, nil}
	return client, nil
}

//
func (c *WSClient) CloseUnderlyingConnection() {
	c.connection.Close()
}

//
func (c *WSClient) CloseConnection() error {
	err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Error("Client write close error: ", err)
		log.Info("Forcing connection close")
		return err
	}
	return nil
}

//
func (c *WSClient) SetSubscribeMessage(message []byte) {
	c.subMsg = message
	return
}

func (c *WSClient) SetUnsubscribeMessage(message []byte) {
	c.unsubMsg = message
	return
}

//
func (c *WSClient) StartStreaming(received chan<- []byte, interrupt <-chan os.Signal) error {
	if c.subMsg == nil {
		log.Error("No subscription message set")
		return errors.New("No subscription message")
	}
	err := c.connection.WriteMessage(websocket.TextMessage, c.subMsg)
	if err != nil {
		log.Error("Client write error: ", err)
		return err
	}

	msgReceived := 0

	for {
		select {
		default:
			_, message, err := c.connection.ReadMessage() // See https://github.com/gorilla/websocket/blob/master/conn.go#L980 on advancement of frames
			// We can only achieve this by calling Conn.NextReader() to skip messages
			if err != nil {
				log.Error("Client read error: ", err)
				defer f.Close()
				return err
			}
			msgReceived++
			received <- message

			log.Debug("recv ", msgReceived, ": ", string(message))
		case <-interrupt:
			log.Info("Received interrupt signal, stopped reading from WS client")
			defer f.Close()
			return nil
		}
	}
}

func (c *WSClient) StopStreaming() error {
	err := c.connection.WriteMessage(websocket.TextMessage, c.unsubMsg)
	if err != nil {
		log.Error("Client write error: ", err)
		return err
	}
	return nil
}
