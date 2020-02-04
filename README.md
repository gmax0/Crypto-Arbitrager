
### Ideal Feature List
1. Arbitrage opportunity notifications/simulations based on cross exchange and cross cryptocurrency data
2. Eventual integration of trade executions in response to discovered arbitrage opportunities
3. Arbitrage discovery should encompass multiple exchanges + N maximum currency hops ("N-angular" arbitrage), all configurable by user

### To-Do's:
---

#### Research 
- [ ] Implement buffered channels to store price pair updates (can gorilla/websockets discard websocket messages in the buffer that are X (mille)seconds long? i.e. determine stale message discard mechanism) 
- [ ] Unified/common representation of price data streamed from various exchanges
- [ ] Ability to stream entire playbook and act off those insights (i.e. all sell orders, all buy orders, etc) instead of using only ticker updates
- [ ] Scalable arbitrage algorithm

#### WebSocket Implementations
Avoid integrations with fake-volume exchanges.

Using this [list here](https://nomics.com/exchanges). Preference given for regulated exchanges with level 2 order book websocket functionality.

[ ] Liquid [documentation here](https://developers.liquid.com/#iii.-liquid-tap-websocket)
[ ] Poloniex websocket client code
[ ] Kraken websocket client code [documentation here](https://docs.kraken.com/websockets/#overview)

#### Potential Algorithms for Arbitrage
- Bellman-Ford and negative cycle detection

#### Performance Considerations
- Is a level 2 order book strictly required? Consider the memory overhead of maintaining vs maintaining only the mid-market price...
- JSON Parsing benchmark tests, look at alternative libraries
- Lock-free ring buffers + atomic operations as an alternative to channels
