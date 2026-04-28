package polymarket

// Types in this file match the Polymarket CLOB market WebSocket JSON payloads (event_type discriminator).
// See https://docs.polymarket.com/api-reference/wss/market

// MarketOrderSummary is one aggregated price level in a book snapshot.
type MarketOrderSummary struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// MarketBookEvent is event_type "book" (full orderbook snapshot).
type MarketBookEvent struct {
	EventType string               `json:"event_type"`
	AssetID   string               `json:"asset_id"`
	Market    string               `json:"market"`
	Bids      []MarketOrderSummary `json:"bids"`
	Asks      []MarketOrderSummary `json:"asks"`
	Timestamp string               `json:"timestamp"`
	Hash      string               `json:"hash"`
}

// MarketPriceChangeItem is one row inside MarketPriceChangeEvent.price_changes.
type MarketPriceChangeItem struct {
	AssetID string `json:"asset_id"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"`
	Hash    string `json:"hash"`
	BestBid string `json:"best_bid,omitempty"`
	BestAsk string `json:"best_ask,omitempty"`
}

// MarketPriceChangeEvent is event_type "price_change".
type MarketPriceChangeEvent struct {
	EventType    string                  `json:"event_type"`
	Market       string                  `json:"market"`
	PriceChanges []MarketPriceChangeItem `json:"price_changes"`
	Timestamp    string                  `json:"timestamp"`
}

// MarketLastTradePriceEvent is event_type "last_trade_price".
type MarketLastTradePriceEvent struct {
	EventType       string `json:"event_type"`
	AssetID         string `json:"asset_id"`
	Market          string `json:"market"`
	Price           string `json:"price"`
	Size            string `json:"size"`
	FeeRateBps      string `json:"fee_rate_bps,omitempty"`
	Side            string `json:"side"`
	Timestamp       string `json:"timestamp"`
	TransactionHash string `json:"transaction_hash,omitempty"`
}

// MarketTickSizeChangeEvent is event_type "tick_size_change".
type MarketTickSizeChangeEvent struct {
	EventType   string `json:"event_type"`
	AssetID     string `json:"asset_id"`
	Market      string `json:"market"`
	OldTickSize string `json:"old_tick_size"`
	NewTickSize string `json:"new_tick_size"`
	Timestamp   string `json:"timestamp"`
}

// MarketBestBidAskEvent is event_type "best_bid_ask" (requires custom_feature_enabled on subscribe).
type MarketBestBidAskEvent struct {
	EventType string `json:"event_type"`
	AssetID   string `json:"asset_id"`
	Market    string `json:"market"`
	BestBid   string `json:"best_bid"`
	BestAsk   string `json:"best_ask"`
	Spread    string `json:"spread"`
	Timestamp string `json:"timestamp"`
}

// MarketEventMessage is nested event_message on new_market / market_resolved payloads.
type MarketEventMessage struct {
	ID          string `json:"id,omitempty"`
	Ticker      string `json:"ticker,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// MarketNewMarketEvent is event_type "new_market" (requires custom_feature_enabled on subscribe).
type MarketNewMarketEvent struct {
	EventType             string              `json:"event_type"`
	ID                    string              `json:"id"`
	Question              string              `json:"question"`
	Market                string              `json:"market"`
	Slug                  string              `json:"slug"`
	Description           string              `json:"description,omitempty"`
	AssetsIDs             []string            `json:"assets_ids"`
	Outcomes              []string            `json:"outcomes"`
	EventMessage          *MarketEventMessage `json:"event_message,omitempty"`
	Timestamp             string              `json:"timestamp"`
	Tags                  []string            `json:"tags,omitempty"`
	ConditionID           string              `json:"condition_id,omitempty"`
	Active                bool                `json:"active,omitempty"`
	ClobTokenIDs          []string            `json:"clob_token_ids,omitempty"`
	SportsMarketType      string              `json:"sports_market_type,omitempty"`
	Line                  string              `json:"line,omitempty"`
	GameStartTime         string              `json:"game_start_time,omitempty"`
	OrderPriceMinTickSize string              `json:"order_price_min_tick_size,omitempty"`
	GroupItemTitle        string              `json:"group_item_title,omitempty"`
}

// MarketResolvedEvent is event_type "market_resolved" (requires custom_feature_enabled on subscribe).
type MarketResolvedEvent struct {
	EventType      string              `json:"event_type"`
	ID             string              `json:"id"`
	Market         string              `json:"market"`
	AssetsIDs      []string            `json:"assets_ids"`
	WinningAssetID string              `json:"winning_asset_id"`
	WinningOutcome string              `json:"winning_outcome"`
	EventMessage   *MarketEventMessage `json:"event_message,omitempty"`
	Timestamp      string              `json:"timestamp"`
	Tags           []string            `json:"tags,omitempty"`
}

// MarketChannelMessage holds at most one decoded market-channel server event.
type MarketChannelMessage struct {
	Book           *MarketBookEvent
	PriceChange    *MarketPriceChangeEvent
	LastTradePrice *MarketLastTradePriceEvent
	TickSizeChange *MarketTickSizeChangeEvent
	BestBidAsk     *MarketBestBidAskEvent
	NewMarket      *MarketNewMarketEvent
	Resolved       *MarketResolvedEvent
}

func (m MarketChannelMessage) empty() bool {
	return m.Book == nil && m.PriceChange == nil && m.LastTradePrice == nil &&
		m.TickSizeChange == nil && m.BestBidAsk == nil && m.NewMarket == nil && m.Resolved == nil
}
