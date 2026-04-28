package polymarket

// CLOB path constants (aligned with @polymarket/clob-client/src/endpoints.ts).

// WebSocket channel URLs (production). See Polymarket API reference /api-reference/wss/*.
const (
	// UserWebSocketURL is the authenticated user channel for order and trade updates.
	// https://docs.polymarket.com/api-reference/wss/user
	UserWebSocketURL = "wss://ws-subscriptions-clob.polymarket.com/ws/user"
	// MarketWebSocketURL is the public market channel (order book, prices, lifecycle).
	// https://docs.polymarket.com/api-reference/wss/market
	MarketWebSocketURL = "wss://ws-subscriptions-clob.polymarket.com/ws/market"
	// SportsWebSocketURL is the public sports results stream (reply with text "pong" to server pings).
	// https://docs.polymarket.com/api-reference/wss/sports
	SportsWebSocketURL = "wss://sports-api.polymarket.com/ws"
)

// CLOBWebSocketPing is the client heartbeat text for user and market channels on ws-subscriptions-clob.
const CLOBWebSocketPing = "PING"

const (
	PathTime = "/time"

	PathCreateAPIKey       = "/auth/api-key"
	PathGetAPIKeys         = "/auth/api-keys"
	PathDeleteAPIKey       = "/auth/api-key"
	PathDeriveAPIKey       = "/auth/derive-api-key"
	PathClosedOnly         = "/auth/ban-status/closed-only"
	PathAccessStatus       = "/auth/access-status"
	PathCreateReadonlyKey  = "/auth/readonly-api-key"
	PathGetReadonlyKeys    = "/auth/readonly-api-keys"
	PathDeleteReadonlyKey  = "/auth/readonly-api-key"
	PathValidateReadonly   = "/auth/validate-readonly-api-key"
	PathCreateBuilderKey   = "/auth/builder-api-key"
	PathGetBuilderKeys     = "/auth/builder-api-key"
	PathRevokeBuilderKey   = "/auth/builder-api-key"

	PathSamplingSimplifiedMarkets = "/sampling-simplified-markets"
	PathSamplingMarkets           = "/sampling-markets"
	PathSimplifiedMarkets         = "/simplified-markets"
	PathMarkets                   = "/markets"
	PathMarketPrefix              = "/markets/"
	PathOrderBook                 = "/book"
	PathOrderBooks                = "/books"
	PathMidpoint                  = "/midpoint"
	PathMidpoints                 = "/midpoints"
	PathPrice                     = "/price"
	PathPrices                    = "/prices"
	PathSpread                    = "/spread"
	PathSpreads                   = "/spreads"
	PathLastTradePrice            = "/last-trade-price"
	PathLastTradesPrices          = "/last-trades-prices"
	PathTickSize                  = "/tick-size"
	PathNegRisk                   = "/neg-risk"
	PathFeeRate                   = "/fee-rate"

	PathPostOrder        = "/order"
	PathPostOrders       = "/orders"
	PathCancelOrder      = "/order"
	PathCancelOrders     = "/orders"
	PathDataOrderPrefix  = "/data/order/"
	PathCancelAll        = "/cancel-all"
	PathCancelMarket     = "/cancel-market-orders"
	PathDataOrders       = "/data/orders"
	PathDataTrades       = "/data/trades"
	PathOrderScoring     = "/order-scoring"
	PathOrdersScoring    = "/orders-scoring"
	PathPricesHistory    = "/prices-history"
	PathNotifications    = "/notifications"
	PathBalanceAllowance = "/balance-allowance"
	PathBalanceUpdate    = "/balance-allowance/update"
	PathLiveActivity     = "/markets/live-activity/"
	PathRewardsUser              = "/rewards/user"
	PathRewardsUserTotal         = "/rewards/user/total"
	PathRewardsPercentages       = "/rewards/user/percentages"
	PathRewardsMarketsCurrent    = "/rewards/markets/current"
	PathRewardsMarketsPrefix     = "/rewards/markets/"
	PathRewardsUserMarkets       = "/rewards/user/markets"
	PathBuilderTrades            = "/builder/trades"
	PathHeartbeats               = "/v1/heartbeats"

	PathRFQCreateRequest   = "/rfq/request"
	PathRFQCancelRequest   = "/rfq/request"
	PathRFQRequests        = "/rfq/data/requests"
	PathRFQCreateQuote     = "/rfq/quote"
	PathRFQCancelQuote     = "/rfq/quote"
	PathRFQRequestAccept   = "/rfq/request/accept"
	PathRFQQuoteApprove    = "/rfq/quote/approve"
	PathRFQRequesterQuotes = "/rfq/data/requester/quotes"
	PathRFQQuoterQuotes    = "/rfq/data/quoter/quotes"
	PathRFQBestQuote       = "/rfq/data/best-quote"
	PathRFQConfig          = "/rfq/config"
)

// Pagination cursors (clob-client/src/constants.ts).
const (
	InitialCursor = "MA=="
	EndCursor     = "LTE="
)
