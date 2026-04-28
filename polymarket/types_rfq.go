package polymarket

// RfqMatchType values for quote/request pairing (clob-client).
const (
	RfqMatchComplementary = "COMPLEMENTARY"
	RfqMatchMerge         = "MERGE"
	RfqMatchMint          = "MINT"
)

// RfqUserOrder is the minimal user input for RFQ request creation.
type RfqUserOrder struct {
	TokenID string
	Price   float64
	Size    float64
	Side    Side
}

// RfqUserQuote adds request id for quoting.
type RfqUserQuote struct {
	RfqUserOrder
	RequestID string
}

// CreateRfqRequestParams is the wire payload for POST /rfq/request.
type CreateRfqRequestParams struct {
	AssetIn   string `json:"assetIn"`
	AssetOut  string `json:"assetOut"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
	UserType  int    `json:"userType"`
}

// CancelRfqRequestParams cancels by request id.
type CancelRfqRequestParams struct {
	RequestID string `json:"requestId"`
}

// CreateRfqQuoteParams wire payload for POST /rfq/quote (userType added by client).
type CreateRfqQuoteParams struct {
	RequestID string `json:"requestId"`
	AssetIn   string `json:"assetIn"`
	AssetOut  string `json:"assetOut"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
	UserType  int    `json:"userType,omitempty"`
}

// CancelRfqQuoteParams cancels a quote.
type CancelRfqQuoteParams struct {
	QuoteID string `json:"quoteId"`
}

// AcceptQuoteParams accepts a quote as taker.
type AcceptQuoteParams struct {
	RequestID string `json:"requestId"`
	QuoteID   string `json:"quoteId"`
	Expiration int64  `json:"expiration"`
}

// ApproveOrderParams approves as maker.
type ApproveOrderParams struct {
	RequestID string `json:"requestId"`
	QuoteID   string `json:"quoteId"`
	Expiration int64 `json:"expiration"`
}

// GetRfqRequestsParams lists RFQ requests.
type GetRfqRequestsParams struct {
	Offset     string
	Limit      int
	State      string
	RequestIDs []string
	Markets    []string
	SizeMin    *float64
	SizeMax    *float64
	SizeUsdcMin *float64
	SizeUsdcMax *float64
	PriceMin   *float64
	PriceMax   *float64
	SortBy     string
	SortDir    string
}

// GetRfqQuotesParams lists RFQ quotes.
type GetRfqQuotesParams struct {
	Offset     string
	Limit      int
	State      string
	QuoteIDs   []string
	RequestIDs []string
	Markets    []string
	SizeMin    *float64
	SizeMax    *float64
	SizeUsdcMin *float64
	SizeUsdcMax *float64
	PriceMin   *float64
	PriceMax   *float64
	SortBy     string
	SortDir    string
}

// GetRfqBestQuoteParams selects best quote.
type GetRfqBestQuoteParams struct {
	RequestID string
}

// RfqQuote is one quote row (subset of fields used by accept/approve).
type RfqQuote struct {
	QuoteID    string  `json:"quoteId"`
	RequestID  string  `json:"requestId"`
	Token      string  `json:"token"`
	Complement string  `json:"complement"`
	Condition  string  `json:"condition"`
	Side       string  `json:"side"`
	SizeIn     string  `json:"sizeIn"`
	SizeOut    string  `json:"sizeOut"`
	Price      float64 `json:"price"`
	MatchType  string  `json:"matchType"`
	State      string  `json:"state"`
}

// RfqQuotesResponse is a paginated quote list.
type RfqQuotesResponse struct {
	Data       []RfqQuote `json:"data"`
	NextCursor string     `json:"next_cursor"`
	Limit      int        `json:"limit"`
	Count      int        `json:"count"`
}
