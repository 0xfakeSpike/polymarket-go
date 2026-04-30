package polymarket

import "time"

// PublicProfileResponse is the Gamma API body for GET /public-profile.
// See https://docs.polymarket.com/api-reference/profiles/get-public-profile-by-wallet-address
type PublicProfileResponse struct {
	CreatedAt             *time.Time          `json:"createdAt"`
	ProxyWallet           *string             `json:"proxyWallet"`
	ProfileImage          *string             `json:"profileImage"`
	DisplayUsernamePublic *bool               `json:"displayUsernamePublic"`
	Bio                   *string             `json:"bio"`
	Pseudonym             *string             `json:"pseudonym"`
	Name                  *string             `json:"name"`
	Users                 []PublicProfileUser `json:"users"`
	XUsername             *string             `json:"xUsername"`
	VerifiedBadge         *bool               `json:"verifiedBadge"`
}

// PublicProfileUser is an entry in PublicProfileResponse.Users.
type PublicProfileUser struct {
	ID      string `json:"id"`
	Creator bool   `json:"creator"`
	Mod     bool   `json:"mod"`
}

// CurrentPositionsParams filters GET https://data-api.polymarket.com/positions
// (https://docs.polymarket.com/api-reference/core/get-current-positions-for-a-user).
type CurrentPositionsParams struct {
	User          string   `json:"user"`
	Market        []string `json:"market,omitempty"`
	EventID       []int    `json:"event_id,omitempty"`
	SizeThreshold *float64 `json:"size_threshold,omitempty"`
	Redeemable    *bool    `json:"redeemable,omitempty"`
	Mergeable     *bool    `json:"mergeable,omitempty"`
	Limit         *int     `json:"limit,omitempty"`
	Offset        *int     `json:"offset,omitempty"`
	SortBy        string   `json:"sort_by,omitempty"`
	SortDirection string   `json:"sort_direction,omitempty"`
	Title         string   `json:"title,omitempty"`
}

// Position is one row from GET /positions.
type Position struct {
	ProxyWallet        string  `json:"proxyWallet"`
	Asset              string  `json:"asset"`
	ConditionID        string  `json:"conditionId"`
	Size               float64 `json:"size"`
	AvgPrice           float64 `json:"avgPrice"`
	InitialValue       float64 `json:"initialValue"`
	CurrentValue       float64 `json:"currentValue"`
	CashPnl            float64 `json:"cashPnl"`
	PercentPnl         float64 `json:"percentPnl"`
	TotalBought        float64 `json:"totalBought"`
	RealizedPnl        float64 `json:"realizedPnl"`
	PercentRealizedPnl float64 `json:"percentRealizedPnl"`
	CurPrice           float64 `json:"curPrice"`
	Redeemable         bool    `json:"redeemable"`
	Mergeable          bool    `json:"mergeable"`
	Title              string  `json:"title"`
	Slug               string  `json:"slug"`
	Icon               string  `json:"icon"`
	EventSlug          string  `json:"eventSlug"`
	Outcome            string  `json:"outcome"`
	OutcomeIndex       int     `json:"outcomeIndex"`
	OppositeOutcome    string  `json:"oppositeOutcome"`
	OppositeAsset      string  `json:"oppositeAsset"`
	EndDate            string  `json:"endDate"`
	NegativeRisk       bool    `json:"negativeRisk"`
}

// ClosedPositionsParams filters GET /closed-positions.
type ClosedPositionsParams struct {
	User          string   `json:"user"`
	Market        []string `json:"market,omitempty"`
	Title         string   `json:"title,omitempty"`
	EventID       []int    `json:"event_id,omitempty"`
	Limit         *int     `json:"limit,omitempty"`
	Offset        *int     `json:"offset,omitempty"`
	SortBy        string   `json:"sort_by,omitempty"`
	SortDirection string   `json:"sort_direction,omitempty"`
}

// ClosedPosition is one row from GET /closed-positions.
type ClosedPosition struct {
	ProxyWallet     string  `json:"proxyWallet"`
	Asset           string  `json:"asset"`
	ConditionID     string  `json:"conditionId"`
	AvgPrice        float64 `json:"avgPrice"`
	TotalBought     float64 `json:"totalBought"`
	RealizedPnl     float64 `json:"realizedPnl"`
	CurPrice        float64 `json:"curPrice"`
	Timestamp       int64   `json:"timestamp"`
	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	Icon            string  `json:"icon"`
	EventSlug       string  `json:"eventSlug"`
	Outcome         string  `json:"outcome"`
	OutcomeIndex    int     `json:"outcomeIndex"`
	OppositeOutcome string  `json:"oppositeOutcome"`
	OppositeAsset   string  `json:"oppositeAsset"`
	EndDate         string  `json:"endDate"`
}

// UserActivityParams filters GET /activity.
type UserActivityParams struct {
	Limit         *int     `json:"limit,omitempty"`
	Offset        *int     `json:"offset,omitempty"`
	User          string   `json:"user"`
	Market        []string `json:"market,omitempty"`
	EventID       []int    `json:"event_id,omitempty"`
	ActivityTypes []string `json:"activity_types,omitempty"` // TRADE, SPLIT, MERGE, … → query param "type"
	Start         *int64   `json:"start,omitempty"`
	End           *int64   `json:"end,omitempty"`
	SortBy        string   `json:"sort_by,omitempty"`
	SortDirection string   `json:"sort_direction,omitempty"`
	Side          string   `json:"side,omitempty"` // BUY | SELL
}

// Activity is one row from GET /activity.
type Activity struct {
	ProxyWallet           string  `json:"proxyWallet"`
	Timestamp             int64   `json:"timestamp"`
	ConditionID           string  `json:"conditionId"`
	Type                  string  `json:"type"`
	Size                  float64 `json:"size"`
	USDCSize              float64 `json:"usdcSize"`
	TransactionHash       string  `json:"transactionHash"`
	Price                 float64 `json:"price"`
	Asset                 string  `json:"asset"`
	Side                  string  `json:"side"`
	OutcomeIndex          int     `json:"outcomeIndex"`
	Title                 string  `json:"title"`
	Slug                  string  `json:"slug"`
	Icon                  string  `json:"icon"`
	EventSlug             string  `json:"eventSlug"`
	Outcome               string  `json:"outcome"`
	Name                  string  `json:"name"`
	Pseudonym             string  `json:"pseudonym"`
	Bio                   string  `json:"bio"`
	ProfileImage          string  `json:"profileImage"`
	ProfileImageOptimized string  `json:"profileImageOptimized"`
}
