package CoinbasePro

import (
    "log"
    "net/url"
    "time"

    "github.com/gorilla/websocket"
    "github.com/nu7hatch/gouuid"
)

/*
 *
 * See https://docs.pro.coinbase.com/#websocket-feed for latest specifications
 *
 */

//CoinbasePro available channels to subscribe to
const (
    HeartbeatChannel = "hearbeat"
    StatusChannel = "status"
    TickerChannel = "ticker"
    Level2Channel = "level2"
    UserChannel = "user"
    MatchesChannel = "matches"
    FullChannel = "full"
)

//Possible message types sent by the CoinbasePro websocket
const (
    HeartbeatType = "heartbeat"
    StatusType = "status"
    TickerType = "ticker"
    SnapshotType = "snapshot"
    Level2Type = "l2update"

    //Full Channel Exclusive Message Types
    ReceivedType = "received"
    OpenType = "open"
    DoneType = "done"
    MatchType = "match"
    ChangeType = "change"
    ActivateType = "activate"
)

/*******************************************************************************/

//Subscription Message
type channel struct {
    name        string   `json:"name"`
    productIds  []string `json:"product_ids"`
}

type SubscriptionMessage struct {
    messageType string    `json:"type"`
    channels    []channel `json:"channels"`
}

//Update Messages (e.g. price tickers, heartbeats, etc.)
/*
type UpdateMessage struct {
    type      string `json:"type"`
    productId string `json:"product_id"`

}
*/

type CoinbaseProWSClient struct {
    Connection *websocket.Conn
    Id         *uuid.UUID
}


/*******************************************************************************/

func NewClient(socketUrl string) (*CoinbaseProWSClient, error) {
    //Use default dialer for now
    u := url.URL{Scheme: "wss", Host: socketUrl}
    log.Printf("Connecting to %s", u.String())

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Println("Dialer initialization error:", err)
        return nil, err
    }

    id, err := uuid.NewV4()
    if err != nil {
        log.Println("Error generating V4 UUID:", err)
    }
    client := &CoinbaseProWSClient{ c, id }
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

func (c *CoinbaseProWSClient) StreamMessages(received chan<- []byte) {
    msgReceived := 0

    for {
        time.Sleep(time.Second * 10)

        // log.Println("Connection buffer size:", c.Connection.br.Size())
        _, message, err := c.Connection.ReadMessage()
        if err != nil {
            log.Println("CoinbasePro Client read error: ", err)
            return
        }
        received <- message

        //msgReceived++
        //if (msgReceived % 20 == 0) {
            log.Printf("recv %d: %s", msgReceived, message)
        //}
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

