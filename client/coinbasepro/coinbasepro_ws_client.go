package coinbasepro

import (
	"log"
	"net/url"
    "strings"
	_ "time"

	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"

    "../../bookkeeper/pricebook"
)

/*
 *
 * See https://docs.pro.coinbase.com/#websocket-feed for latest specifications
 *
 */

//CoinbasePro available channels to subscribe to
const (
	HeartbeatChannel = "hearbeat"
	StatusChannel    = "status"
	TickerChannel    = "ticker"
	Level2Channel    = "level2"
	UserChannel      = "user"
	MatchesChannel   = "matches"
	FullChannel      = "full"
)

//Possible message types sent by the CoinbasePro websocket
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

//Subscription Message
type channel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type SubscriptionMessage struct {
	MessageType string    `json:"type"`
	Channels    []channel `json:"channels"`
}

//Update Messages (e.g. price tickers, heartbeats, etc.)
type UpdateMessage struct {
    messageType string `json:"type"`
    productId   string `json:"product_id"`

}

type SnapshotMessage struct {
    messageType string    `json:"type"`
    productId   string    `json:"product_id"`
    asks        []string  `json:"asks"`
    bids        []string  `json:"bids"`
}

type CoinbaseProWSClient struct {
	Connection *websocket.Conn
	ID         *uuid.UUID
}

/*******************************************************************************/

//
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
		log.Println("CoinbasePro Client write close error: ", err)
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
			log.Println("CoinbasePro Client read error: ", err)
			return
		}
        msgReceived++
		received <- message

		log.Printf("recv %d: %s", msgReceived, message)
	}
}

//
func (c *CoinbaseProWSClient) Subscribe(message []byte) error {
	err := c.Connection.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("CoinbasePro Client write error: ", err)
		return err
	}

	return nil
}

// ParsePricePairs parses the initial subscription message that is sent to Coinbase Pro's websocket.
// All price pairs that are to be subscribed to will receive an entry in the adjacency matrix (map[map[X]]string) graph. 
// e.g. ETH-BTC will result in two entries with keys
func ParsePricePairs(subscriptionMessage SubscriptionMessage) *map[string]map[string]pricebook.Edge {


    //TODO: Add validation for SubscriptionMessage (must have L2). Assuming one channel, L2 channel is contained in the subscriptionMessage
    var m map[string]map[string]pricebook.Edge
    m = make(map[string]map[string]pricebook.Edge)


    for _, pair := range subscriptionMessage.Channels[0].ProductIds {
        split := strings.Split(pair, "-")
        if m[split[0]] == nil {
            m[split[0]] = make(map[string]pricebook.Edge)
        }
        m[split[0]][split[1]] = pricebook.Edge{ High: -1, Low: -1}

        if m[split[1]] == nil {
            m[split[1]] = make(map[string]pricebook.Edge)
        }
        m[split[1]][split[0]] = pricebook.Edge{ High: -2, Low: -2}
    }

    return &m  
}
/*
func (c *CoinbaseProWSClient) Unsubscribe(message []byte) error {

}
*/
