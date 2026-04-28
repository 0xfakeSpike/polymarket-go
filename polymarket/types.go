package polymarket

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
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

type OptimizedImage struct {
	ID                        string `json:"id"`
	ImageUrlSource            string `json:"imageUrlSource"`
	ImageUrlOptimized         string `json:"imageUrlOptimized"`
	ImageSizeKbSource         int    `json:"imageSizeKbSource"`
	ImageSizeKbOptimized      int    `json:"imageSizeKbOptimized"`
	ImageOptimizedComplete    bool   `json:"imageOptimizedComplete"`
	ImageOptimizedLastUpdated string `json:"imageOptimizedLastUpdated"`
	RelID                     int    `json:"relID"`
	Field                     string `json:"field"`
	Relname                   string `json:"relname"`
}

// FlexibleStringSlice unmarshals JSON that is either a string array or a string containing a JSON-encoded array
// (e.g. `["Yes","No"]` or `"[\"0.5\",\"0.5\"]"`). Used for Market.outcomes and Market.outcomePrices across API variants.
type FlexibleStringSlice []string

func (s *FlexibleStringSlice) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*s = nil
		return nil
	}
	var direct []string
	if err := json.Unmarshal(data, &direct); err == nil {
		*s = direct
		return nil
	}
	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return fmt.Errorf("flexible string slice: expected []string or string: %w", err)
	}
	if asString == "" {
		*s = nil
		return nil
	}
	if err := json.Unmarshal([]byte(asString), &direct); err != nil {
		return fmt.Errorf("flexible string slice: parse string as JSON array: %w", err)
	}
	*s = direct
	return nil
}

// MarshalJSON encodes as a normal JSON array.
func (s FlexibleStringSlice) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]string(s))
}

type Market struct {
	ID                           string              `json:"id"`
	Question                     string              `json:"question"`
	ConditionID                  string              `json:"conditionId"`
	Slug                         string              `json:"slug"`
	TwitterCardImage             string              `json:"twitterCardImage"`
	ResolutionSource             string              `json:"resolutionSource"`
	EndDate                      *time.Time          `json:"endDate"`
	Category                     string              `json:"category"`
	AmmType                      string              `json:"ammType"`
	Liquidity                    string              `json:"liquidity"`
	SponsorName                  string              `json:"sponsorName"`
	SponsorImage                 string              `json:"sponsorImage"`
	StartDate                    *time.Time          `json:"startDate"`
	XAxisValue                   string              `json:"xAxisValue"`
	YAxisValue                   string              `json:"yAxisValue"`
	DenominationToken            string              `json:"denominationToken"`
	Fee                          string              `json:"fee"`
	Image                        string              `json:"image"`
	Icon                         string              `json:"icon"`
	LowerBound                   string              `json:"lowerBound"`
	UpperBound                   string              `json:"upperBound"`
	Description                  string              `json:"description"`
	Outcomes                     FlexibleStringSlice `json:"outcomes"`
	OutcomesPrices               FlexibleStringSlice `json:"outcomePrices"`
	Volume                       string              `json:"volume"`
	Active                       bool                `json:"active"`
	MarketType                   string              `json:"marketType"`
	FormatType                   string              `json:"formatType"`
	LowerBoundDate               string              `json:"lowerBoundDate"`
	UpperBoundDate               string              `json:"upperBoundDate"`
	Closed                       bool                `json:"closed"`
	MarketMakerAddress           string              `json:"marketMakerAddress"`
	CreatedBy                    int                 `json:"createdBy"`
	UpdatedBy                    int                 `json:"updatedBy"`
	CreatedAt                    *time.Time          `json:"createdAt"`
	UpdatedAt                    *time.Time          `json:"updatedAt"`
	ClosedTime                   string              `json:"closedTime"`
	WideFormat                   bool                `json:"wideFormat"`
	New                          bool                `json:"new"`
	MailchimpTag                 string              `json:"mailchimpTag"`
	Featured                     bool                `json:"featured"`
	Archived                     bool                `json:"archived"`
	ResolvedBy                   string              `json:"resolvedBy"`
	Restricted                   bool                `json:"restricted"`
	MarketGroup                  int                 `json:"marketGroup"`
	GroupItemTitle               string              `json:"groupItemTitle"`
	GroupItemThreshold           string              `json:"groupItemThreshold"`
	QuestionID                   string              `json:"questionID"`
	UmaEndDate                   string              `json:"umaEndDate"`
	EnableOrderBook              bool                `json:"enableOrderBook"`
	OrderPriceMinTickSize        float64             `json:"orderPriceMinTickSize"`
	OrderMinSize                 float64             `json:"orderMinSize"`
	UmaResolutionStatus          string              `json:"umaResolutionStatus"`
	CurationOrder                int                 `json:"curationOrder"`
	VolumeNum                    float64             `json:"volumeNum"`
	LiquidityNum                 float64             `json:"liquidityNum"`
	EndDateIso                   string              `json:"endDateIso"`
	StartDateIso                 string              `json:"startDateIso"`
	UmaEndDateIso                string              `json:"umaEndDateIso"`
	HasReviewedDates             bool                `json:"hasReviewedDates"`
	ReadyForCron                 bool                `json:"readyForCron"`
	CommentsEnabled              bool                `json:"commentsEnabled"`
	Volume24hr                   float64             `json:"volume24hr"`
	Volume1wk                    float64             `json:"volume1wk"`
	Volume1mo                    float64             `json:"volume1mo"`
	Volume1yr                    float64             `json:"volume1yr"`
	GameStartTime                string              `json:"gameStartTime"`
	SecondsDelay                 int                 `json:"secondsDelay"`
	ClobTokenIds                 FlexibleStringSlice `json:"clobTokenIds"`
	DisqusThread                 string              `json:"disqusThread"`
	ShortOutcomes                string              `json:"shortOutcomes"`
	TeamAID                      string              `json:"teamAID"`
	TeamBID                      string              `json:"teamBID"`
	UmaBond                      string              `json:"umaBond"`
	UmaReward                    string              `json:"umaReward"`
	FpmmLive                     bool                `json:"fpmmLive"`
	Volume24hrAmm                float64             `json:"volume24hrAmm"`
	Volume1wkAmm                 float64             `json:"volume1wkAmm"`
	Volume1moAmm                 float64             `json:"volume1moAmm"`
	Volume1yrAmm                 float64             `json:"volume1yrAmm"`
	Volume24hrClob               float64             `json:"volume24hrClob"`
	Volume1wkClob                float64             `json:"volume1wkClob"`
	Volume1moClob                float64             `json:"volume1moClob"`
	Volume1yrClob                float64             `json:"volume1yrClob"`
	VolumeAmm                    float64             `json:"volumeAmm"`
	VolumeClob                   float64             `json:"volumeClob"`
	LiquidityAmm                 float64             `json:"liquidityAmm"`
	LiquidityClob                float64             `json:"liquidityClob"`
	MakerBaseFee                 float64             `json:"makerBaseFee"`
	TakerBaseFee                 float64             `json:"takerBaseFee"`
	CustomLiveness               float64             `json:"customLiveness"`
	AcceptingOrders              bool                `json:"acceptingOrders"`
	NotificationsEnabled         bool                `json:"notificationsEnabled"`
	Score                        float64             `json:"score"`
	ImageOptimized               *OptimizedImage     `json:"imageOptimized"`
	IconOptimized                *OptimizedImage     `json:"iconOptimized"`
	Events                       []Event             `json:"events"`
	Categories                   []Category          `json:"categories"`
	Tags                         []Tag               `json:"tags"`
	Creator                      string              `json:"creator"`
	Ready                        bool                `json:"ready"`
	Funded                       bool                `json:"funded"`
	PastSlugs                    string              `json:"pastSlugs"`
	ReadyTimestamp               *time.Time          `json:"readyTimestamp"`
	FundedTimestamp              *time.Time          `json:"fundedTimestamp"`
	AcceptingOrdersTimestamp     *time.Time          `json:"acceptingOrdersTimestamp"`
	Competitive                  float64             `json:"competitive"`
	RewardsMinSize               float64             `json:"rewardsMinSize"`
	RewardsMaxSpread             float64             `json:"rewardsMaxSpread"`
	Spread                       float64             `json:"spread"`
	AutomaticallyResolved        bool                `json:"automaticallyResolved"`
	OneDayPriceChange            float64             `json:"oneDayPriceChange"`
	OneHourPriceChange           float64             `json:"oneHourPriceChange"`
	OneWeekPriceChange           float64             `json:"oneWeekPriceChange"`
	OneMonthPriceChange          float64             `json:"oneMonthPriceChange"`
	OneYearPriceChange           float64             `json:"oneYearPriceChange"`
	LastTradePrice               float64             `json:"lastTradePrice"`
	BestBid                      float64             `json:"bestBid"`
	BestAsk                      float64             `json:"bestAsk"`
	AutomaticallyActive          bool                `json:"automaticallyActive"`
	ClearBookOnStart             bool                `json:"clearBookOnStart"`
	ChartColor                   string              `json:"chartColor"`
	SeriesColor                  string              `json:"seriesColor"`
	ShowGmpSeries                bool                `json:"showGmpSeries"`
	ShowGmpOutcome               bool                `json:"showGmpOutcome"`
	ManualActivation             bool                `json:"manualActivation"`
	NegRiskOther                 bool                `json:"negRiskOther"`
	GameId                       string              `json:"gameId"`
	GroupItemRange               string              `json:"groupItemRange"`
	SportsMarketType             string              `json:"sportsMarketType"`
	Line                         float64             `json:"line"`
	UmaResolutionStatuses        string              `json:"umaResolutionStatuses"`
	PendingDeployment            bool                `json:"pendingDeployment"`
	Deploying                    bool                `json:"deploying"`
	DeployingTimestamp           *time.Time          `json:"deployingTimestamp"`
	ScheduledDeploymentTimestamp *time.Time          `json:"scheduledDeploymentTimestamp"`
	RfqEnabled                   bool                `json:"rfqEnabled"`
	EventStartTime               *time.Time          `json:"eventStartTime"`
	SubmittedBy                  string              `json:"submitted_by"`
	NegRisk                      bool                `json:"negRisk"`
	NegRiskRequestID             string              `json:"negRiskRequestID"`
	Cyom                         bool                `json:"cyom"`
	PagerDutyNotificationEnabled bool                `json:"pagerDutyNotificationEnabled"`
	Approved                     bool                `json:"approved"`
	ClobRewards                  []ClobReward        `json:"clobRewards"`
	HoldingRewardsEnabled        bool                `json:"holdingRewardsEnabled"`
	FeesEnabled                  bool                `json:"feesEnabled"`
}

// FavoredSidePNL is the PnL snapshot for the side whose live best ask is above 0.5.
// PnL assumes buying one share now and holding until market settlement.
type FavoredSidePNL struct {
	Outcome       string        `json:"outcome"`
	TokenID       string        `json:"token_id,omitempty"`
	Price         float64       `json:"price"`
	BestAsk       float64       `json:"best_ask,omitempty"`
	BestBid       float64       `json:"best_bid,omitempty"`
	HoldingPeriod time.Duration `json:"holding_period"`
	PnLPerShare   float64       `json:"pnl_per_share"`
	ROI           float64       `json:"roi"`
	// AnnualizedReturn is the compound annualized rate implied by ROI over HoldingPeriod:
	// (1+ROI)^(1/years) - 1, with years = holding / 365.25d. Nil if not computable (e.g. non-positive holding).
	AnnualizedReturn *float64  `json:"annualized_return,omitempty"`
	SettlementTime   time.Time `json:"settlement_time"`
	ComputedAt       time.Time `json:"computed_at"`
}

// FavoredSidePNLFromOrderBooks computes favored-side PnL from live order books.
// It selects the side with best ask > 0.5, then uses that best ask as entry price.
//
// Formula:
//   - Cost per share = best ask
//   - PnL per share at settlement (if that side resolves true) = 1 - best ask
//   - ROI = (1 - best ask) / best ask
//   - Annualized return (compound): (1+ROI)^(1/years)-1, years = (EndDate-now) / 365.25 days
func (m *Market) FavoredSidePNLFromOrderBooks(now time.Time, booksByToken map[string]*Book) (*FavoredSidePNL, error) {
	if m == nil {
		return nil, fmt.Errorf("market is nil")
	}
	if m.EndDate == nil {
		return nil, fmt.Errorf("market endDate is nil")
	}
	if len(m.ClobTokenIds) == 0 {
		return nil, fmt.Errorf("market clobTokenIds is empty")
	}
	if now.IsZero() {
		now = time.Now()
	}
	if !m.EndDate.After(now) {
		return nil, fmt.Errorf("market already settled or passed endDate")
	}

	bestIdx := -1
	bestAsk := 0.0
	bestBid := 0.0
	bestTokenID := ""
	for i, tokenID := range m.ClobTokenIds {
		book := booksByToken[tokenID]
		if book == nil {
			continue
		}
		asks := book.AsksData()
		bids := book.BidsData()
		if len(asks) == 0 {
			continue
		}
		ask := asks[0].Price
		bid := 0.0
		if len(bids) > 0 {
			bid = bids[0].Price
		}
		if ask > 0.5 && ask > bestAsk {
			bestIdx = i
			bestAsk = ask
			bestBid = bid
			bestTokenID = tokenID
		}
	}
	if bestIdx < 0 {
		return nil, fmt.Errorf("no order book best ask > 0.5")
	}

	return m.buildFavoredSidePNL(now, bestIdx, bestAsk, bestAsk, bestBid, bestTokenID)
}

func (m *Market) buildFavoredSidePNL(now time.Time, outcomeIdx int, price float64, bestAsk float64, bestBid float64, tokenID string) (*FavoredSidePNL, error) {
	if now.IsZero() {
		now = time.Now()
	}
	if m.EndDate == nil || !m.EndDate.After(now) {
		return nil, fmt.Errorf("market already settled or passed endDate")
	}
	if price <= 0 {
		return nil, fmt.Errorf("invalid price %.8f", price)
	}
	outcome := ""
	if outcomeIdx >= 0 && outcomeIdx < len(m.Outcomes) {
		outcome = m.Outcomes[outcomeIdx]
	}
	pnlPerShare := 1 - price
	roi := pnlPerShare / price
	holding := m.EndDate.Sub(now)
	return &FavoredSidePNL{
		Outcome:          outcome,
		TokenID:          tokenID,
		Price:            price,
		BestAsk:          bestAsk,
		BestBid:          bestBid,
		HoldingPeriod:    holding,
		PnLPerShare:      pnlPerShare,
		ROI:              roi,
		AnnualizedReturn: annualizedReturnFromROI(roi, holding),
		SettlementTime:   *m.EndDate,
		ComputedAt:       now,
	}, nil
}

// annualizedReturnFromROI returns the compound annualized rate (1+ROI)^(1/y)-1
// using an average year of 365.25 solar days. Returns nil if the horizon is not positive
// or if 1+ROI is not positive.
func annualizedReturnFromROI(roi float64, holding time.Duration) *float64 {
	if holding <= 0 {
		return nil
	}
	if roi <= -1 {
		return nil
	}
	years := holding.Seconds() / (365.25 * 24 * 3600)
	if years <= 0 {
		return nil
	}
	v := math.Pow(1+roi, 1/years) - 1
	return &v
}

// ClobReward represents CLOB reward information
type ClobReward struct {
	ID               string  `json:"id"`
	ConditionID      string  `json:"conditionId"`
	AssetAddress     string  `json:"assetAddress"`
	RewardsAmount    float64 `json:"rewardsAmount"`
	RewardsDailyRate float64 `json:"rewardsDailyRate"`
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
}

// Token represents a stake in a specific Yes/No outcome in a Market
// Price fluctuates between 0-1 and is redeemable for $1 USDC upon resolution
type Token struct {
	ID      string `json:"id"`
	TokenID string `json:"token_id"`
	Outcome string `json:"outcome"`
	Price   string `json:"price"`
	Winner  *bool  `json:"winner"`
}

// Event represents a collection of related markets
type Event struct {
	ID                           string          `json:"id"`
	Ticker                       string          `json:"ticker"`
	Slug                         string          `json:"slug"`
	Title                        string          `json:"title"`
	Subtitle                     string          `json:"subtitle"`
	Description                  string          `json:"description"`
	ResolutionSource             string          `json:"resolutionSource"`
	StartDate                    *time.Time      `json:"startDate"`
	CreationDate                 *time.Time      `json:"creationDate"`
	EndDate                      *time.Time      `json:"endDate"`
	Image                        string          `json:"image"`
	Icon                         string          `json:"icon"`
	Active                       bool            `json:"active"`
	Closed                       bool            `json:"closed"`
	Archived                     bool            `json:"archived"`
	New                          bool            `json:"new"`
	Featured                     bool            `json:"featured"`
	Restricted                   bool            `json:"restricted"`
	Liquidity                    float64         `json:"liquidity"`
	Volume                       float64         `json:"volume"`
	OpenInterest                 float64         `json:"openInterest"`
	SortBy                       string          `json:"sortBy"`
	Category                     string          `json:"category"`
	Subcategory                  string          `json:"subcategory"`
	IsTemplate                   bool            `json:"isTemplate"`
	TemplateVariables            string          `json:"templateVariables"`
	PublishedAt                  string          `json:"published_at"`
	CreatedBy                    string          `json:"createdBy"`
	UpdatedBy                    string          `json:"updatedBy"`
	CreatedAt                    *time.Time      `json:"createdAt"`
	UpdatedAt                    *time.Time      `json:"updatedAt"`
	CommentsEnabled              bool            `json:"commentsEnabled"`
	Competitive                  float64         `json:"competitive"`
	Volume24hr                   float64         `json:"volume24hr"`
	Volume1wk                    float64         `json:"volume1wk"`
	Volume1mo                    float64         `json:"volume1mo"`
	Volume1yr                    float64         `json:"volume1yr"`
	FeaturedImage                string          `json:"featuredImage"`
	DisqusThread                 string          `json:"disqusThread"`
	ParentEvent                  string          `json:"parentEvent"`
	EnableOrderBook              bool            `json:"enableOrderBook"`
	LiquidityAmm                 float64         `json:"liquidityAmm"`
	LiquidityClob                float64         `json:"liquidityClob"`
	NegRisk                      bool            `json:"negRisk"`
	NegRiskMarketID              string          `json:"negRiskMarketID"`
	NegRiskFeeBips               float64         `json:"negRiskFeeBips"`
	CommentCount                 int             `json:"commentCount"`
	ImageOptimized               *OptimizedImage `json:"imageOptimized"`
	IconOptimized                *OptimizedImage `json:"iconOptimized"`
	FeaturedImageOptimized       *OptimizedImage `json:"featuredImageOptimized"`
	SubEvents                    []string        `json:"subEvents"`
	Markets                      []Market        `json:"markets"`
	Series                       []Series        `json:"series"`
	Categories                   []Category      `json:"categories"`
	Collections                  []Collection    `json:"collections"`
	Tags                         []Tag           `json:"tags"`
	CYOM                         bool            `json:"cyom"`
	ClosedTime                   *time.Time      `json:"closedTime"`
	ShowAllOutcomes              bool            `json:"showAllOutcomes"`
	ShowMarketImages             bool            `json:"showMarketImages"`
	AutomaticallyResolved        bool            `json:"automaticallyResolved"`
	EnableNegRisk                bool            `json:"enableNegRisk"`
	AutomaticallyActive          bool            `json:"automaticallyActive"`
	EventDate                    string          `json:"eventDate"`
	StartTime                    *time.Time      `json:"startTime"`
	EventWeek                    int             `json:"eventWeek"`
	SeriesSlug                   string          `json:"seriesSlug"`
	Score                        string          `json:"score"`
	Elapsed                      string          `json:"elapsed"`
	Period                       string          `json:"period"`
	Live                         bool            `json:"live"`
	Ended                        bool            `json:"ended"`
	FinishedTimestamp            *time.Time      `json:"finishedTimestamp"`
	GmpChartMode                 string          `json:"gmpChartMode"`
	EventCreators                []EventCreator  `json:"eventCreators"`
	TweetCount                   int             `json:"tweetCount"`
	Chats                        []Chat          `json:"chats"`
	FeaturedOrder                int             `json:"featuredOrder"`
	EstimateValue                bool            `json:"estimateValue"`
	CantEstimate                 bool            `json:"cantEstimate"`
	EstimatedValue               string          `json:"estimatedValue"`
	Templates                    []Template      `json:"templates"`
	SpreadsMainLine              float64         `json:"spreadsMainLine"`
	TotalsMainLine               float64         `json:"totalsMainLine"`
	CarouselMap                  string          `json:"carouselMap"`
	PendingDeployment            bool            `json:"pendingDeployment"`
	Deploying                    bool            `json:"deploying"`
	DeployingTimestamp           *time.Time      `json:"deployingTimestamp"`
	ScheduledDeploymentTimestamp *time.Time      `json:"scheduledDeploymentTimestamp"`
	GameStatus                   string          `json:"gameStatus"`
}

// Collection represents a collection of events
type Collection struct {
	ID                   string          `json:"id"`
	Ticker               string          `json:"ticker"`
	Slug                 string          `json:"slug"`
	Title                string          `json:"title"`
	Subtitle             string          `json:"subtitle"`
	CollectionType       string          `json:"collectionType"`
	Description          string          `json:"description"`
	Tags                 string          `json:"tags"`
	Image                string          `json:"image"`
	Icon                 string          `json:"icon"`
	HeaderImage          string          `json:"headerImage"`
	Layout               string          `json:"layout"`
	Active               bool            `json:"active"`
	Closed               bool            `json:"closed"`
	Archived             bool            `json:"archived"`
	New                  bool            `json:"new"`
	Featured             bool            `json:"featured"`
	Restricted           bool            `json:"restricted"`
	IsTemplate           bool            `json:"isTemplate"`
	TemplateVariables    string          `json:"templateVariables"`
	PublishedAt          string          `json:"publishedAt"`
	CreatedBy            string          `json:"createdBy"`
	UpdatedBy            string          `json:"updatedBy"`
	CreatedAt            *time.Time      `json:"createdAt"`
	UpdatedAt            *time.Time      `json:"updatedAt"`
	CommentsEnabled      bool            `json:"commentsEnabled"`
	ImageOptimized       *OptimizedImage `json:"imageOptimized"`
	IconOptimized        *OptimizedImage `json:"iconOptimized"`
	HeaderImageOptimized *OptimizedImage `json:"headerImageOptimized"`
}

// Chat represents a chat channel
type Chat struct {
	ID           string     `json:"id"`
	ChannelId    string     `json:"channelId"`
	ChannelName  string     `json:"channelName"`
	ChannelImage string     `json:"channelImage"`
	Live         bool       `json:"live"`
	StartTime    *time.Time `json:"startTime"`
	EndTime      *time.Time `json:"endTime"`
}

// EventCreator represents an event creator
type EventCreator struct {
	ID            string     `json:"id"`
	CreatorName   string     `json:"creatorName"`
	CreatorHandle string     `json:"creatorHandle"`
	CreatorUrl    string     `json:"creatorUrl"`
	CreatorImage  string     `json:"creatorImage"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

// Template represents a market template
type Template struct {
	ID               string `json:"id"`
	EventTitle       string `json:"eventTitle"`
	EventSlug        string `json:"eventSlug"`
	EventImage       string `json:"eventImage"`
	MarketTitle      string `json:"marketTitle"`
	Description      string `json:"description"`
	ResolutionSource string `json:"resolutionSource"`
	NegRisk          bool   `json:"negRisk"`
	SortBy           string `json:"sortBy"`
	ShowMarketImages bool   `json:"showMarketImages"`
	SeriesSlug       string `json:"seriesSlug"`
	Outcomes         string `json:"outcomes"`
}

// Series represents a series of related events
type Series struct {
	ID                string       `json:"id"`
	Ticker            string       `json:"ticker"`
	Slug              string       `json:"slug"`
	Title             string       `json:"title"`
	Subtitle          string       `json:"subtitle"`
	SeriesType        string       `json:"seriesType"`
	Recurrence        string       `json:"recurrence"`
	Description       string       `json:"description"`
	Image             string       `json:"image"`
	Icon              string       `json:"icon"`
	Layout            string       `json:"layout"`
	Active            bool         `json:"active"`
	Closed            bool         `json:"closed"`
	Archived          bool         `json:"archived"`
	New               bool         `json:"new"`
	Featured          bool         `json:"featured"`
	Restricted        bool         `json:"restricted"`
	IsTemplate        bool         `json:"isTemplate"`
	TemplateVariables bool         `json:"templateVariables"`
	PublishedAt       string       `json:"publishedAt"`
	CreatedBy         string       `json:"createdBy"`
	UpdatedBy         string       `json:"updatedBy"`
	CreatedAt         *time.Time   `json:"createdAt"`
	UpdatedAt         *time.Time   `json:"updatedAt"`
	CommentsEnabled   bool         `json:"commentsEnabled"`
	Competitive       string       `json:"competitive"`
	Volume24hr        float64      `json:"volume24hr"`
	Volume            float64      `json:"volume"`
	Liquidity         float64      `json:"liquidity"`
	StartDate         *time.Time   `json:"startDate"`
	PythTokenID       string       `json:"pythTokenID"`
	CgAssetName       string       `json:"cgAssetName"`
	Score             float64      `json:"score"`
	Events            []Event      `json:"events"`
	Collections       []Collection `json:"collections"`
	Categories        []Category   `json:"categories"`
	Tags              []Tag        `json:"tags"`
	CommentCount      int          `json:"commentCount"`
	Chats             []Chat       `json:"chats"`
}

// Tag represents an event tag
type Tag struct {
	ID          string     `json:"id"`
	Label       string     `json:"label"`
	Slug        string     `json:"slug"`
	ForceShow   bool       `json:"forceShow"`
	PublishedAt string     `json:"publishedAt"`
	CreatedBy   int        `json:"createdBy"`
	UpdatedBy   int        `json:"updatedBy"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	ForceHide   bool       `json:"forceHide"`
	IsCarousel  bool       `json:"isCarousel"`
}

// Category represents a market category/tag
type Category struct {
	ID             string     `json:"id"`
	Label          string     `json:"label"`
	ParentCategory string     `json:"parentCategory"`
	Slug           string     `json:"slug"`
	PublishedAt    string     `json:"publishedAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedBy      string     `json:"updatedBy"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

// MarketsParams represents query parameters for listing markets
type MarketsParams struct {
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Order     string `json:"order,omitempty"`
	Ascending bool   `json:"ascending,omitempty"`

	// Filters
	Closed   *bool  `json:"closed,omitempty"`
	Archived *bool  `json:"archived,omitempty"`
	Slug     string `json:"slug,omitempty"`
	EventID  string `json:"event_id,omitempty"`
	TagID    string `json:"tag_id,omitempty"`
}

// MarketsKeysetParams represents query parameters for /markets/keyset.
// This endpoint uses cursor-based pagination and rejects offset.
type MarketsKeysetParams struct {
	Limit       int    `json:"limit,omitempty"`
	Order       string `json:"order,omitempty"` // e.g. volume_num,liquidity_num
	Ascending   bool   `json:"ascending,omitempty"`
	AfterCursor string `json:"after_cursor,omitempty"`
	Locale      string `json:"locale,omitempty"`

	// Common filters
	ID                  []string   `json:"id,omitempty"`
	Slug                []string   `json:"slug,omitempty"`
	Closed              *bool      `json:"closed,omitempty"`
	Decimalized         *bool      `json:"decimalized,omitempty"`
	ClobTokenIDs        []string   `json:"clob_token_ids,omitempty"`
	ConditionIDs        []string   `json:"condition_ids,omitempty"`
	QuestionIDs         []string   `json:"question_ids,omitempty"`
	MarketMakerAddress  []string   `json:"market_maker_address,omitempty"`
	LiquidityNumMin     *float64   `json:"liquidity_num_min,omitempty"`
	LiquidityNumMax     *float64   `json:"liquidity_num_max,omitempty"`
	VolumeNumMin        *float64   `json:"volume_num_min,omitempty"`
	VolumeNumMax        *float64   `json:"volume_num_max,omitempty"`
	StartDateMin        *time.Time `json:"start_date_min,omitempty"`
	StartDateMax        *time.Time `json:"start_date_max,omitempty"`
	EndDateMin          *time.Time `json:"end_date_min,omitempty"`
	EndDateMax          *time.Time `json:"end_date_max,omitempty"`
	TagID               []int      `json:"tag_id,omitempty"`
	RelatedTags         *bool      `json:"related_tags,omitempty"`
	TagMatch            string     `json:"tag_match,omitempty"`
	CYOM                *bool      `json:"cyom,omitempty"`
	RFQEnabled          *bool      `json:"rfq_enabled,omitempty"`
	UMAResolutionStatus string     `json:"uma_resolution_status,omitempty"`
	GameID              string     `json:"game_id,omitempty"`
	SportsMarketTypes   []string   `json:"sports_market_types,omitempty"`
	IncludeTag          *bool      `json:"include_tag,omitempty"`
}

// EventsParams represents query parameters for listing events
type EventsParams struct {
	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Sorting
	Order     string `json:"order,omitempty"` // volume24hr
	Ascending bool   `json:"ascending,omitempty"`

	// Basic filters
	ID       []string `json:"id,omitempty"`
	Slug     []string `json:"slug,omitempty"`
	Active   *bool    `json:"active,omitempty"`
	Closed   *bool    `json:"closed,omitempty"`
	Archived *bool    `json:"archived,omitempty"`

	// Advanced filters
	TagID        *int   `json:"tag_id,omitempty"`
	ExcludeTagID []int  `json:"exclude_tag_id,omitempty"`
	RelatedTags  *bool  `json:"related_tags,omitempty"`
	Featured     *bool  `json:"featured,omitempty"`
	CYOM         *bool  `json:"cyom,omitempty"`
	Recurrence   string `json:"recurrence,omitempty"`
	TagSlug      string `json:"tag_slug,omitempty"`

	// Date filters
	StartDateMin *time.Time `json:"start_date_min,omitempty"`
	StartDateMax *time.Time `json:"start_date_max,omitempty"`
	EndDateMin   *time.Time `json:"end_date_min,omitempty"`
	EndDateMax   *time.Time `json:"end_date_max,omitempty"`
}

// GetEventParams represents query parameters for getting a single event by ID
type GetEventParams struct {
	IncludeChat     *bool `json:"include_chat,omitempty"`
	IncludeTemplate *bool `json:"include_template,omitempty"`
}

// EventsPaginationParams represents query parameters for the events pagination endpoint
type EventsPaginationParams struct {
	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Sorting
	Order     string `json:"order,omitempty"` // volume, volume24hr, etc.
	Ascending bool   `json:"ascending,omitempty"`

	// Basic filters
	Active   *bool `json:"active,omitempty"`
	Archived *bool `json:"archived,omitempty"`
	Closed   *bool `json:"closed,omitempty"`

	// Tag filters
	TagSlug []string `json:"tag_slug,omitempty"`
}

// EventsKeysetParams represents query parameters for GET /events/keyset (cursor pagination; offset is rejected).
// See https://docs.polymarket.com/api-reference/events/list-events-keyset-pagination
type EventsKeysetParams struct {
	Limit       int    `json:"limit,omitempty"`
	Order       string `json:"order,omitempty"`
	Ascending   bool   `json:"ascending,omitempty"`
	AfterCursor string `json:"after_cursor,omitempty"`
	Locale      string `json:"locale,omitempty"`

	ID               []int      `json:"id,omitempty"`
	Slug             []string   `json:"slug,omitempty"`
	Closed           *bool      `json:"closed,omitempty"`
	Live             *bool      `json:"live,omitempty"`
	Featured         *bool      `json:"featured,omitempty"`
	CYOM             *bool      `json:"cyom,omitempty"`
	TitleSearch      string     `json:"title_search,omitempty"`
	LiquidityMin     *float64   `json:"liquidity_min,omitempty"`
	LiquidityMax     *float64   `json:"liquidity_max,omitempty"`
	VolumeMin        *float64   `json:"volume_min,omitempty"`
	VolumeMax        *float64   `json:"volume_max,omitempty"`
	StartDateMin     *time.Time `json:"start_date_min,omitempty"`
	StartDateMax     *time.Time `json:"start_date_max,omitempty"`
	EndDateMin       *time.Time `json:"end_date_min,omitempty"`
	EndDateMax       *time.Time `json:"end_date_max,omitempty"`
	StartTimeMin     *time.Time `json:"start_time_min,omitempty"`
	StartTimeMax     *time.Time `json:"start_time_max,omitempty"`
	TagID            []int      `json:"tag_id,omitempty"`
	TagSlug          string     `json:"tag_slug,omitempty"`
	ExcludeTagID     []int      `json:"exclude_tag_id,omitempty"`
	RelatedTags      *bool      `json:"related_tags,omitempty"`
	TagMatch         string     `json:"tag_match,omitempty"`
	SeriesID         []int      `json:"series_id,omitempty"`
	GameID           []int      `json:"game_id,omitempty"`
	EventDate        *time.Time `json:"event_date,omitempty"`
	EventWeek        *int       `json:"event_week,omitempty"`
	FeaturedOrder    *bool      `json:"featured_order,omitempty"`
	Recurrence       string     `json:"recurrence,omitempty"`
	CreatedBy        []string   `json:"created_by,omitempty"`
	ParentEventID    *int       `json:"parent_event_id,omitempty"`
	IncludeChildren  *bool      `json:"include_children,omitempty"`
	PartnerSlug      string     `json:"partner_slug,omitempty"`
	IncludeChat      *bool      `json:"include_chat,omitempty"`
	IncludeTemplate  *bool      `json:"include_template,omitempty"`
	IncludeBestLines *bool      `json:"include_best_lines,omitempty"`
}

// Comment represents a comment on a market, event, or series
type Comment struct {
	ID               string       `json:"id"`
	Body             string       `json:"body"`
	ParentEntityType string       `json:"parentEntityType"`
	ParentEntityID   string       `json:"parentEntityID"`
	UserAddress      string       `json:"userAddress"`
	CreatedAt        *time.Time   `json:"createdAt"`
	Profile          *UserProfile `json:"profile"`
	Reactions        []Reaction   `json:"reactions"`
	ReportCount      int          `json:"reportCount"`
	ReactionCount    int          `json:"reactionCount"`
}

// UserProfile represents user profile information
type UserProfile struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
	Bio         string `json:"bio"`
	Website     string `json:"website"`
	Verified    bool   `json:"verified"`
}

// Reaction represents a reaction to a comment
type Reaction struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	UserAddress string     `json:"userAddress"`
	CreatedAt   *time.Time `json:"createdAt"`
}

// CommentsParams represents query parameters for listing comments
type CommentsParams struct {
	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Sorting
	Order     string `json:"order,omitempty"`
	Ascending bool   `json:"ascending,omitempty"`

	// Filters
	ParentEntityType string `json:"parent_entity_type,omitempty"` // "Event", "Series", "market"
	ParentEntityID   *int   `json:"parent_entity_id,omitempty"`
	GetPositions     *bool  `json:"get_positions,omitempty"`
	HoldersOnly      *bool  `json:"holders_only,omitempty"`
}

// SearchParams represents query parameters for search
type SearchParams struct {
	// Required
	Q string `json:"q"` // Search query

	Type string `json:"type"`

	Presets []string `json:"presets,omitempty"`

	// Pagination
	Page         int `json:"page,omitempty"`
	LimitPerType int `json:"limit_per_type,omitempty"`

	// Sorting
	Sort      string `json:"sort,omitempty"`
	Ascending bool   `json:"ascending,omitempty"`

	// Filters and options
	Cache             *bool    `json:"cache,omitempty"`
	EventsStatus      string   `json:"events_status,omitempty"`
	EventsTag         []string `json:"events_tag,omitempty"`
	KeepClosedMarkets *int     `json:"keep_closed_markets,omitempty"`
	SearchTags        *bool    `json:"search_tags,omitempty"`
	SearchProfiles    *bool    `json:"search_profiles,omitempty"`
	Recurrence        string   `json:"recurrence,omitempty"`
	ExcludeTagID      []int    `json:"exclude_tag_id,omitempty"`
	Optimized         *bool    `json:"optimized,omitempty"`
}

// SearchResults represents the unified search response
type SearchResults struct {
	Events     []Event       `json:"events"`
	Tags       []Tag         `json:"tags"`
	Profiles   []UserProfile `json:"profiles"`
	Pagination Pagination    `json:"pagination"`
}

// Pagination represents pagination information in search results
type Pagination struct {
	HasMore      bool `json:"hasMore"`
	TotalResults int  `json:"totalResults"`
}

// LiveVolume represents live volume data for an event
type LiveVolume struct {
	Total   float64        `json:"total"`
	Markets []MarketVolume `json:"markets"`
}

// MarketVolume represents volume data for a specific market
type MarketVolume struct {
	Market string  `json:"market"` // Market address/ID
	Value  float64 `json:"value"`  // Volume value
}

// Book represents the order book for a market
type Book struct {
	Market       string       `json:"market"` // market CanditionID
	AssetID      string       `json:"asset_id"`
	Timestamp    string       `json:"timestamp"`
	Hash         string       `json:"hash"`
	Bids         []LimitOrder `json:"bids"`
	Asks         []LimitOrder `json:"asks"`
	MinOrderSize string       `json:"min_order_size"`
	TickSize     string       `json:"tick_size"`
	NegRisk      bool         `json:"neg_risk"`
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

// BiggestMover represents a biggest mover item
type BiggestMover struct {
	ID          string  `json:"id"`
	Question    string  `json:"question"`
	Slug        string  `json:"slug"`
	PriceChange float64 `json:"priceChange"`
	Volume      float64 `json:"volume"`
	Liquidity   float64 `json:"liquidity"`
	LastPrice   float64 `json:"lastPrice"`
	Image       string  `json:"image"`
	Icon        string  `json:"icon"`
	Category    string  `json:"category"`
	Active      bool    `json:"active"`
	Closed      bool    `json:"closed"`
	EndDate     string  `json:"endDate"`
}

type EventPaginationResponse struct {
	Events     []Event    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// EventsKeysetResponse is the response envelope from GET /events/keyset (KeysetEventsResponse).
type EventsKeysetResponse struct {
	Events     []Event `json:"events"`
	NextCursor string  `json:"next_cursor,omitempty"`
}

// MarketsKeysetResponse is the response envelope from /markets/keyset.
type MarketsKeysetResponse struct {
	Markets    []Market `json:"markets"`
	NextCursor string   `json:"next_cursor,omitempty"`
}

// APIKeyCredentials represents API key credentials
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

// APIKeyInfo represents API key information
type APIKeyInfo struct {
	ApiKey    string    `json:"apiKey"`
	CreatedAt time.Time `json:"createdAt"`
	IsActive  bool      `json:"isActive"`
}

// APIKeysResponse matches clob-client GET /auth/api-keys.
type APIKeysResponse struct {
	APIKeys []APIKeyCredentials `json:"apiKeys"`
}

// AccessStatus represents access status information
type AccessStatus struct {
	CertRequired bool `json:"cert_required"`
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
