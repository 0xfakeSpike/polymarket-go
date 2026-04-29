package polymarket

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

const (
	YESOPTION = "YES"
	NOOPTION  = "NO"
)

// Side is the CLOB order side.
type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

// BookQuote is one level of the order book (price × size).
type BookQuote struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// OrderRef is a minimal order handle from the CLOB.
type OrderRef struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// 订单类型常量
const (
	OrderTypeFOK = "FOK" // Fill-Or-Kill
	OrderTypeFAK = "FAK" // Fill-And-Kill
	OrderTypeGTC = "GTC" // Good-Til-Cancelled
	OrderTypeGTD = "GTD" // Good-Til-Date
)

type Book struct {
	Market         string       `json:"market"` // market CanditionID
	AssetID        string       `json:"asset_id"`
	Timestamp      string       `json:"timestamp"`
	Hash           string       `json:"hash"`
	Bids           []LimitOrder `json:"bids"`
	Asks           []LimitOrder `json:"asks"`
	MinOrderSize   string       `json:"min_order_size"`
	TickSize       string       `json:"tick_size"`
	NegRisk        bool         `json:"neg_risk"`
	LastTradePrice string       `json:"last_trade_price"`
}

func (b *Book) BidsData() []LimitOrder {
	bids := make([]LimitOrder, len(b.Bids))
	for i, bid := range b.Bids {
		bids[i] = LimitOrder{Price: bid.Price, Size: bid.Size}
	}
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price > bids[j].Price
	})
	return bids
}

func (b *Book) AsksData() []LimitOrder {
	asks := make([]LimitOrder, len(b.Asks))
	for i, ask := range b.Asks {
		asks[i] = LimitOrder{Price: ask.Price, Size: ask.Size}
	}
	sort.Slice(asks, func(i, j int) bool {
		return asks[i].Price < asks[j].Price
	})
	return asks
}

type LimitOrder struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

func parseWireFloat(raw json.RawMessage) (float64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return 0, nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return 0, err
		}
		if s == "" {
			return 0, nil
		}
		return strconv.ParseFloat(s, 64)
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return 0, err
	}
	return f, nil
}

func (lo *LimitOrder) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var w struct {
		Price json.RawMessage `json:"price"`
		Size  json.RawMessage `json:"size"`
	}
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	var err error
	if lo.Price, err = parseWireFloat(w.Price); err != nil {
		return fmt.Errorf("limit order price: %w", err)
	}
	if lo.Size, err = parseWireFloat(w.Size); err != nil {
		return fmt.Errorf("limit order size: %w", err)
	}
	return nil
}

// APIKeyCredentials represents CLOB API key credentials.
type APIKeyCredentials struct {
	ApiKey     string `json:"apiKey"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

// UnmarshalJSON accepts both create-key shape ("apiKey") and list-keys shape ("key") from the CLOB API.
func (a *APIKeyCredentials) UnmarshalJSON(data []byte) error {
	type wire struct {
		Key        string `json:"key"`
		APIKey     string `json:"apiKey"`
		Secret     string `json:"secret"`
		Passphrase string `json:"passphrase"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	a.ApiKey = w.Key
	if a.ApiKey == "" {
		a.ApiKey = w.APIKey
	}
	a.Secret = w.Secret
	a.Passphrase = w.Passphrase
	return nil
}

// APIKeysResponse matches clob-client GET /auth/api-keys.
type APIKeysResponse struct {
	APIKeys []APIKeyCredentials `json:"apiKeys"`
}

// ClosedOnlyModeStatus represents closed-only mode status
type ClosedOnlyModeStatus struct {
	ClosedOnly bool `json:"closed_only"`
}

// L1AuthHeaders represents L1 authentication headers
type L1AuthHeaders struct {
	PolyAddress   string `json:"POLY_ADDRESS"`
	PolySignature string `json:"POLY_SIGNATURE"`
	PolyTimestamp string `json:"POLY_TIMESTAMP"`
	PolyNonce     string `json:"POLY_NONCE"`
}

// L2AuthHeaders represents L2 authentication headers
type L2AuthHeaders struct {
	PolyAddress    string `json:"POLY_ADDRESS"`
	PolySignature  string `json:"POLY_SIGNATURE"`
	PolyTimestamp  string `json:"POLY_TIMESTAMP"`
	PolyAPIKey     string `json:"POLY_API_KEY"`
	PolyPassphrase string `json:"POLY_PASSPHRASE"`
}

// OrderRequest 订单请求参数
type OrderRequest struct {
	TokenID     string  `json:"tokenID"`
	Price       float64 `json:"price"`
	Side        Side    `json:"side"`
	Size        float64 `json:"size"`
	Metadata    string  `json:"metadata,omitempty"`    // bytes32 hex，默认 0x00...00
	BuilderCode string  `json:"builderCode,omitempty"` // bytes32 hex，默认 0x00...00
	Expiration  *int64  `json:"expiration,omitempty"`  // unix 秒，GTD 可用；nil 表示不过期
}

// MarketOrderRequest is a market (FOK/FAK) order. Price optional — computed from the book when 0.
type MarketOrderRequest struct {
	Side        Side
	TokenID     string
	Amount      float64
	Price       float64
	Metadata    string // bytes32 hex，默认 0x00...00
	BuilderCode string // bytes32 hex，默认 0x00...00
}

// OrderResponse 订单提交响应
type OrderResponse struct {
	Success            bool     `json:"success"`
	ErrorMsg           string   `json:"errorMsg"`
	OrderID            string   `json:"orderID"`
	TransactionsHashes []string `json:"transactionsHashes"`
	Status             string   `json:"status"`
	TakingAmount       string   `json:"takingAmount"`
	MakingAmount       string   `json:"makingAmount"`
}

type GetOrderResponse struct {
	AssociateTrades []string `json:"associate_trades"`
	ID              string   `json:"id"`
	Status          string   `json:"status"`
	Market          string   `json:"market"`
	OriginalSize    string   `json:"original_size"`
	Outcome         string   `json:"outcome"`
	MakerAddress    string   `json:"maker_address"`
	Owner           string   `json:"owner"`
	Price           string   `json:"price"`
	Side            string   `json:"side"`
	SizeMatched     string   `json:"size_matched"`
	AssetID         string   `json:"asset_id"`
	Expiration      string   `json:"expiration"`
	Type            string   `json:"type"`
	CreatedAt       string   `json:"created_at"`
}
