package CoinbasePro

import (
    "log"
    "net/url"

    "github.com/gorilla/websocket"
)

/*
struct channel {
    name        string   `json:"name"`
    productIds  []string `json:"product_ids"`
}

struct message {
    messageType string   `json:"type"`
    channels    channels `json:"channels"`
}
*/

type CoinbaseProWSClient struct {
    Connection *websocket.Conn
}

func NewClient(socketUrl string) (*CoinbaseProWSClient, error) {
    //Use default dialer for now
    u := url.URL{Scheme: "wss", Host: socketUrl}
    log.Printf("Connecting to %s", u.String())

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Println("Dialer initialization error:", err)
        return nil, err
    }

    client := &CoinbaseProWSClient{ c }
    return client, nil
}

func (c *CoinbaseProWSClient) CloseUnderlyingConnection() {
    c.Connection.Close()
}

func (c *CoinbaseProWSClient) CloseConnection() error {
    err := c.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
    if err != nil {
        log.Println("CoinbasePro Client write close error: ", err)
        return err
    }
    return nil
}

func (c *CoinbaseProWSClient) StreamMessages() {
    for {
        _, message, err := c.Connection.ReadMessage()
        if err != nil {
            log.Println("CoinbasePro Client read error: ", err)
            return
        }
        log.Printf("recv: %s", message)
    }
}

func (c *CoinbaseProWSClient) Subscribe(message []byte) error {
    err := c.Connection.WriteMessage(websocket.TextMessage, message)
    if err != nil {
        log.Println("CoinbasePro Client write error: ", err)
        return err
    }

    return nil
}


/*
func (c *CoinbaseProWSClient) Unsubscribe(message []byte) error {

}
*/

