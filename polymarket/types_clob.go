package polymarket

import (
	"encoding/json"
)

// PaginationPayload is the CLOB paginated list envelope (cursor-based).
type PaginationPayload struct {
	Limit      int             `json:"limit"`
	Count      int             `json:"count"`
	NextCursor string          `json:"next_cursor"`
	Data       json.RawMessage `json:"data"`
}

// BookParams is one token/side pair for batch price endpoints.
type BookParams struct {
	TokenID string `json:"token_id"`
	Side    Side   `json:"side"`
}

// TradeParams filters CLOB trades listing.
type TradeParams struct {
	ID            string `json:"id,omitempty"`
	MakerAddress  string `json:"maker_address,omitempty"`
	Market        string `json:"market,omitempty"`
	AssetID       string `json:"asset_id,omitempty"`
	Before        string `json:"before,omitempty"`
	After         string `json:"after,omitempty"`
}

// OpenOrderParams filters open orders.
type OpenOrderParams struct {
	ID      string `json:"id,omitempty"`
	Market  string `json:"market,omitempty"`
	AssetID string `json:"asset_id,omitempty"`
}

// OrderPayload cancels a single order by server id.
type OrderPayload struct {
	OrderID string `json:"orderID"`
}

// OrderMarketCancelParams cancels orders in a market.
type OrderMarketCancelParams struct {
	Market  string `json:"market,omitempty"`
	AssetID string `json:"asset_id,omitempty"`
}

// DropNotificationParams drops notifications by id.
type DropNotificationParams struct {
	IDs []string `json:"ids"`
}

// BalanceAllowanceParams queries collateral / conditional balance.
type BalanceAllowanceParams struct {
	AssetType string `json:"asset_type"` // COLLATERAL | CONDITIONAL
	TokenID   string `json:"token_id,omitempty"`
}

// BalanceAllowanceResponse is CLOB balance snapshot.
type BalanceAllowanceResponse struct {
	Balance   string `json:"balance"`
	Allowance string `json:"allowance"`
}

// OrderScoringParams asks if one order is scoring liquidity.
type OrderScoringParams struct {
	OrderID string `json:"order_id"`
}

// OrderScoring is the scoring flag response.
type OrderScoring struct {
	Scoring bool `json:"scoring"`
}

// OrdersScoringParams batch order scoring.
type OrdersScoringParams struct {
	OrderIDs []string `json:"orderIds"`
}

// OrdersScoring maps order id -> scoring.
type OrdersScoring map[string]bool

// PriceHistoryFilterParams filters /prices-history.
type PriceHistoryFilterParams struct {
	Market   string `json:"market,omitempty"`
	StartTs  int64  `json:"startTs,omitempty"`
	EndTs    int64  `json:"endTs,omitempty"`
	Fidelity int    `json:"fidelity,omitempty"`
	Interval string `json:"interval,omitempty"`
}

// MarketPrice is one point in price history.
type MarketPrice struct {
	T int64   `json:"t"`
	P float64 `json:"p"`
}

// HeartbeatResponse from POST /v1/heartbeats.
type HeartbeatResponse struct {
	HeartbeatID string `json:"heartbeat_id"`
	Error       string `json:"error,omitempty"`
}

// OpenOrder is one row from GET /data/orders.
type OpenOrder struct {
	ID             string   `json:"id"`
	Status         string   `json:"status"`
	Owner          string   `json:"owner"`
	MakerAddress   string   `json:"maker_address"`
	Market         string   `json:"market"`
	AssetID        string   `json:"asset_id"`
	Side           string   `json:"side"`
	OriginalSize   string   `json:"original_size"`
	SizeMatched    string   `json:"size_matched"`
	Price          string   `json:"price"`
	AssociateTrades []string `json:"associate_trades"`
	Outcome        string   `json:"outcome"`
	CreatedAt      float64  `json:"created_at"`
	Expiration     string   `json:"expiration"`
	OrderType      string   `json:"order_type"`
}

// Trade is a single trade from GET /data/trades.
type Trade struct {
	ID              string `json:"id"`
	TakerOrderID    string `json:"taker_order_id"`
	Market          string `json:"market"`
	AssetID         string `json:"asset_id"`
	Side            string `json:"side"`
	Size            string `json:"size"`
	FeeRateBps      string `json:"fee_rate_bps"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	MatchTime       string `json:"match_time"`
	LastUpdate      string `json:"last_update"`
	Outcome         string `json:"outcome"`
	BucketIndex     int    `json:"bucket_index"`
	Owner           string `json:"owner"`
	MakerAddress    string `json:"maker_address"`
	TransactionHash string `json:"transaction_hash"`
	TraderSide      string `json:"trader_side"`
}

// TradesPage is one page from GET /data/trades.
type TradesPage struct {
	Trades     []Trade `json:"data"`
	NextCursor string  `json:"next_cursor"`
	Limit      int     `json:"limit"`
	Count      int     `json:"count"`
}
