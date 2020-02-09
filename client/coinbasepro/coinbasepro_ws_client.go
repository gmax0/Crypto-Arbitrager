package coinbasepro

import (
	"net/url"
	_ "time"

	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
    "github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {

}

/*
 *
 * See https://docs.pro.coinbase.com/#websocket-feed for latest specifications
 *
 */

// CoinbasePro available channels to subscribe to
const (
	HeartbeatChannel = "hearbeat"
	StatusChannel    = "status"
	TickerChannel    = "ticker"
	Level2Channel    = "level2"
	UserChannel      = "user"
	MatchesChannel   = "matches"
	FullChannel      = "full"
)

// Possible message types sent by the CoinbasePro websocket
const (
	HeartbeatType = "heartbeat"
	StatusType    = "status"
	TickerType    = "ticker"
	SnapshotType  = "snapshot"
	Level2Type    = "l2update"

	//Full Channel Exclusive Message Types
	ReceivedType = "received"
	OpenType     = "open"
	DoneType     = "done"
	MatchType    = "match"
	ChangeType   = "change"
	ActivateType = "activate"
)

/*******************************************************************************/

// Subscription Message
type Channel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type SubscriptionMessage struct {
	MessageType string    `json:"type"`
	Channels    []Channel `json:"channels"`
}

/*******************************************************************************/

// L2 Channel Response Messages
/*
type L2UpdateMessage struct {
    messageType string `json:"type"`
    productId   string `json:"product_id"`

}
*/

/*
type Ask struct {
    Entry []string
}

type Bid struct {
    Entry []string
}
*/

type L2SnapshotMessage struct {
    MessageType string    `json:"type"`
    ProductId   string    `json:"product_id"`
    Asks        [][]string `json:"asks"`
    Bids        [][]string `json:"bids"`
}


/*******************************************************************************/
// Status Channel Response Message
/*
type Currency struct {

}
type StatusResponseMessage struct {
    MessageType string `json:"type"`
    Currencies []Currency `json:"currencies"`
    Products   []Product `json:"products"`
}
*/

/*******************************************************************************/

type CoinbaseProWSClient struct {
    Connection *websocket.Conn
    ID         *uuid.UUID
}

//
func NewClient(socketUrl string) (*CoinbaseProWSClient, error) {
	//Use default dialer for now
	u := url.URL{Scheme: "wss", Host: socketUrl}
	log.Info("Connecting to ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error("Dialer initialization error:", err)
		return nil, err
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Error("Error generating V4 UUID:", err)
	}
	client := &CoinbaseProWSClient{c, id}
	return client, nil
}

//
func (c *CoinbaseProWSClient) CloseUnderlyingConnection() {
	c.Connection.Close()
}

//
func (c *CoinbaseProWSClient) CloseConnection() error {
	err := c.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Error("CoinbasePro Client write close error: ", err)
		return err
	}
	return nil
}

//
func (c *CoinbaseProWSClient) Subscribe(message []byte) error {
    err := c.Connection.WriteMessage(websocket.TextMessage, message)
    if err != nil {
        log.Error("CoinbasePro Client write error: ", err)
        return err
    }

    return nil
}

//
func (c *CoinbaseProWSClient) StreamMessages(received chan<- []byte) {
	msgReceived := 0

	for {
		_, message, err := c.Connection.ReadMessage() // See https://github.com/gorilla/websocket/blob/master/conn.go#L980 on advancement of frames
                                                      // We can only achieve this by calling Conn.NextReader() to skip messages 
		if err != nil {
			log.Error("CoinbasePro Client read error: ", err)
			return
		}
        msgReceived++
		received <- message

		log.Trace("recv ", msgReceived, ": ", string(message))
	}
}


/*
func (c *CoinbaseProWSClient) Unsubscribe(message []byte) error {

}
*/
