### Overview

Golang based application that aggregates real-time cryptocurrency price data across various exchanges and processes this data to determine exploitable arbitrage opportunities. 

If performance bottlenecks cannot be resolved by improving application architecture, they may necessitate a change to a more performant language, e.g. Rust.

### Exchange Integrations

Most cryptocurrency exchanges provide websocket implementations, allowing clients to subscribe to receive price data for price pairs of interest.
Typically, a [level-2 order book](https://www.investopedia.com/articles/trading/06/level2quotes.asp) offers a sufficent amount of information to calculate arbitrage opportunities. For a given price pair's order book, bids and asks are represented as price levels, the mid-market price lying between the delta of the highest bid and the lowest ask. Websockets transmit updates (usually JSON-formatted messages) detailing changes to price levels and these updates are to be used to synchronize an application's in-memory order book with the exchange server's.

Furthermore, in addition to websockets, the exchanges offer REST APIs that allow clients to execute orders. Clients can leverage the APIs to place trades based on automated trading strategies.

Worth noting are the two classifications of exchanges. The first, "centralized" exchanges, is encompassed within the scope of this application. These exchanges are privately owned and operated, are highly reliable, and are of relatively low ease to use. The second, "decentralized" exchanges, provide P2P trading by utilizing smart contract capabilities. [Further research is needed...]

#### Websocket Client Implementation

The code for various exchange integrations can be found under client/{ExchangeName}. Each implementation uses the gorilla/websocket library to interface with exchange websockets. JSON-formatted byte messages are transmitted to and from the exchange websockets by the clients. Generally, to allow for accurate synchronization of price data between client and server, following subscription, a snapshot of the current level-2 order book is provided to the client by the server, followed by subsequent updates containing individual price levels.

Relationships between threads, clients, exchanges, price-pairs:
- 1 client per 1 thread
- 1 client per 1 exchange
- N price-pair subscriptions per 1 client
Note that these relationships may change depending on scaling needs.

#### Client-Side Order Book Implementation

See bookkeeper/ and bookkeeper/orderbook.

Bookkeeper retains all the orderbooks for each price-pair for every exchange (this is subject to change depending on performance bottlenecks). It uses parser/{Exchange} to parse raw byte messages into canoniclized structures used by the underlying orderbook as well as the graph implementations.

The data structure selected for the orderbook implementation requires quick lookups, insertions, deletions, and ideally constant retrieval of the minimum ask price level and maximum bid price level. Balanced, ordered binary search trees are ideal for such requirements, offering O(logN) upserts, deletions, and O(1) min/max retrievals given a trivial implementation maintaining such values on upsertion/deletion.

#### Graph Implementation + Bellman-Ford Algorithm for Arbitrage Detection
To-do.
