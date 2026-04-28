package polymarket

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// CreateRfqRequest builds and posts an RFQ request (matches rfq-client createRfqRequest).
func (c *Client) CreateRfqRequest(user RfqUserOrder) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	tick, err := c.ResolveTickSize(user.TokenID, nil)
	if err != nil {
		return nil, err
	}
	cfg, ok := roundingConfig[tick]
	if !ok {
		return nil, fmt.Errorf("unsupported tick %q", tick)
	}
	roundedPrice := roundNormal(user.Price, cfg.price)
	roundedSize := roundDown(user.Size, cfg.size)
	priceStr := strconv.FormatFloat(roundedPrice, 'f', cfg.price, 64)
	sizeStr := strconv.FormatFloat(roundedSize, 'f', cfg.size, 64)
	priceNum, _ := strconv.ParseFloat(priceStr, 64)
	sizeNum, _ := strconv.ParseFloat(sizeStr, 64)

	var amountIn, amountOut, assetIn, assetOut string
	if user.Side == SideBuy {
		amountIn, _ = parseUnitsHuman(sizeNum, collateralTokenDecimals)
		out := sizeNum * priceNum
		amountOut, _ = parseUnitsHuman(out, collateralTokenDecimals)
		assetIn = user.TokenID
		assetOut = "0"
	} else {
		out := sizeNum * priceNum
		amountIn, _ = parseUnitsHuman(out, collateralTokenDecimals)
		amountOut, _ = parseUnitsHuman(sizeNum, collateralTokenDecimals)
		assetIn = "0"
		assetOut = user.TokenID
	}

	payload := CreateRfqRequestParams{
		AssetIn: assetIn, AssetOut: assetOut,
		AmountIn: amountIn, AmountOut: amountOut,
		UserType: c.signatureType,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	path := PathRFQCreateRequest
	h, err := c.buildL2AuthHeaders("POST", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", path, nil, h, b)
}

// CancelRfqRequest cancels an RFQ request.
func (c *Client) CancelRfqRequest(p CancelRfqRequestParams) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	path := PathRFQCancelRequest
	h, err := c.buildL2AuthHeaders("DELETE", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("DELETE", path, nil, h, b)
}

func appendRfqRequestQuery(q url.Values, p *GetRfqRequestsParams) {
	if p == nil {
		return
	}
	if p.Offset != "" {
		q.Set("offset", p.Offset)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.State != "" {
		q.Set("state", p.State)
	}
	for _, id := range p.RequestIDs {
		q.Add("requestIds", id)
	}
	for _, m := range p.Markets {
		q.Add("markets", m)
	}
	setFloat := func(k string, v *float64) {
		if v != nil {
			q.Set(k, strconv.FormatFloat(*v, 'f', -1, 64))
		}
	}
	setFloat("sizeMin", p.SizeMin)
	setFloat("sizeMax", p.SizeMax)
	setFloat("sizeUsdcMin", p.SizeUsdcMin)
	setFloat("sizeUsdcMax", p.SizeUsdcMax)
	setFloat("priceMin", p.PriceMin)
	setFloat("priceMax", p.PriceMax)
	if p.SortBy != "" {
		q.Set("sortBy", p.SortBy)
	}
	if p.SortDir != "" {
		q.Set("sortDir", p.SortDir)
	}
}

// GetRfqRequests lists RFQ requests.
func (c *Client) GetRfqRequests(params *GetRfqRequestsParams) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRFQRequests
	h, err := c.buildL2AuthHeaders("GET", path, "")
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	appendRfqRequestQuery(q, params)
	return c.clobRequest("GET", path, q, h, nil)
}

// CreateRfqQuote posts a quote on a request.
func (c *Client) CreateRfqQuote(user RfqUserQuote) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	tick, err := c.ResolveTickSize(user.TokenID, nil)
	if err != nil {
		return nil, err
	}
	cfg, ok := roundingConfig[tick]
	if !ok {
		return nil, fmt.Errorf("unsupported tick %q", tick)
	}
	roundedPrice := roundNormal(user.Price, cfg.price)
	roundedSize := roundDown(user.Size, cfg.size)
	priceStr := strconv.FormatFloat(roundedPrice, 'f', cfg.price, 64)
	sizeStr := strconv.FormatFloat(roundedSize, 'f', cfg.size, 64)
	priceNum, _ := strconv.ParseFloat(priceStr, 64)
	sizeNum, _ := strconv.ParseFloat(sizeStr, 64)

	var amountIn, amountOut, assetIn, assetOut string
	if user.Side == SideSell {
		out := sizeNum * priceNum
		amountIn, _ = parseUnitsHuman(out, collateralTokenDecimals)
		amountOut, _ = parseUnitsHuman(sizeNum, collateralTokenDecimals)
		assetIn = "0"
		assetOut = user.TokenID
	} else {
		amountIn, _ = parseUnitsHuman(sizeNum, collateralTokenDecimals)
		out := sizeNum * priceNum
		amountOut, _ = parseUnitsHuman(out, collateralTokenDecimals)
		assetIn = user.TokenID
		assetOut = "0"
	}

	payload := map[string]any{
		"requestId": user.RequestID,
		"assetIn":   assetIn,
		"assetOut":  assetOut,
		"amountIn":  amountIn,
		"amountOut": amountOut,
		"userType":  c.signatureType,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	path := PathRFQCreateQuote
	h, err := c.buildL2AuthHeaders("POST", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", path, nil, h, b)
}

func appendRfqQuotesQuery(q url.Values, p *GetRfqQuotesParams) {
	if p == nil {
		return
	}
	if p.Offset != "" {
		q.Set("offset", p.Offset)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.State != "" {
		q.Set("state", p.State)
	}
	for _, id := range p.QuoteIDs {
		q.Add("quoteIds", id)
	}
	for _, id := range p.RequestIDs {
		q.Add("requestIds", id)
	}
	for _, m := range p.Markets {
		q.Add("markets", m)
	}
	setFloat := func(k string, v *float64) {
		if v != nil {
			q.Set(k, strconv.FormatFloat(*v, 'f', -1, 64))
		}
	}
	setFloat("sizeMin", p.SizeMin)
	setFloat("sizeMax", p.SizeMax)
	setFloat("sizeUsdcMin", p.SizeUsdcMin)
	setFloat("sizeUsdcMax", p.SizeUsdcMax)
	setFloat("priceMin", p.PriceMin)
	setFloat("priceMax", p.PriceMax)
	if p.SortBy != "" {
		q.Set("sortBy", p.SortBy)
	}
	if p.SortDir != "" {
		q.Set("sortDir", p.SortDir)
	}
}

// GetRfqRequesterQuotes returns quotes on the user's requests.
func (c *Client) GetRfqRequesterQuotes(params *GetRfqQuotesParams) (*RfqQuotesResponse, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRFQRequesterQuotes
	h, err := c.buildL2AuthHeaders("GET", path, "")
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	appendRfqQuotesQuery(q, params)
	data, err := c.clobRequest("GET", path, q, h, nil)
	if err != nil {
		return nil, err
	}
	var out RfqQuotesResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRfqQuoterQuotes returns quotes created by the user.
func (c *Client) GetRfqQuoterQuotes(params *GetRfqQuotesParams) (*RfqQuotesResponse, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRFQQuoterQuotes
	h, err := c.buildL2AuthHeaders("GET", path, "")
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	appendRfqQuotesQuery(q, params)
	data, err := c.clobRequest("GET", path, q, h, nil)
	if err != nil {
		return nil, err
	}
	var out RfqQuotesResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRfqBestQuote returns best quote for optional request id.
func (c *Client) GetRfqBestQuote(params *GetRfqBestQuoteParams) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRFQBestQuote
	h, err := c.buildL2AuthHeaders("GET", path, "")
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if params != nil && params.RequestID != "" {
		q.Set("requestId", params.RequestID)
	}
	return c.clobRequest("GET", path, q, h, nil)
}

// CancelRfqQuote cancels a quote.
func (c *Client) CancelRfqQuote(p CancelRfqQuoteParams) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	path := PathRFQCancelQuote
	h, err := c.buildL2AuthHeaders("DELETE", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("DELETE", path, nil, h, b)
}

// RFQConfig returns RFQ server config.
func (c *Client) RFQConfig() (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRFQConfig
	h, err := c.buildL2AuthHeaders("GET", path, "")
	if err != nil {
		return nil, err
	}
	return c.clobRequest("GET", path, nil, h, nil)
}

// AcceptRfqQuote accepts a quote (taker): fetches quote, builds matching order, posts accept.
func (c *Client) AcceptRfqQuote(payload AcceptQuoteParams) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	quotes, err := c.GetRfqRequesterQuotes(&GetRfqQuotesParams{QuoteIDs: []string{payload.QuoteID}})
	if err != nil {
		return nil, err
	}
	if len(quotes.Data) == 0 {
		return nil, fmt.Errorf("RFQ quote not found")
	}
	rfq := quotes.Data[0]
	oc, err := rfqRequestOrderCreationPayload(&rfq)
	if err != nil {
		return nil, err
	}
	req := OrderRequest{
		TokenID:    oc.TokenID,
		Price:      oc.Price,
		Size:       oc.Size,
		Side:       oc.Side,
	}
	exp := payload.Expiration
	req.Expiration = &exp
	signed, err := c.BuildSignedLimitOrder(req)
	if err != nil {
		return nil, err
	}
	acceptPayload := map[string]any{
		"requestId":     payload.RequestID,
		"quoteId":       payload.QuoteID,
		"owner":         c.apiKeyCredentials.ApiKey,
		"salt":          signed.Salt.Int64(),
		"maker":         signed.Maker.Hex(),
		"signer":        signed.Signer.Hex(),
		"tokenId":       signed.TokenID.String(),
		"makerAmount":   signed.MakerAmount.String(),
		"takerAmount":   signed.TakerAmount.String(),
		"side":          string(signed.Side),
		"signatureType": int(signed.SignatureType),
		"timestamp":     signed.Timestamp.String(),
		"metadata":      signed.Metadata.Hex(),
		"builder":       signed.Builder.Hex(),
		"expiration":    signed.Expiration.String(),
		"signature":     encodeOrderSignatureHex(signed.Signature),
	}
	b, err := json.Marshal(acceptPayload)
	if err != nil {
		return nil, err
	}
	path := PathRFQRequestAccept
	h, err := c.buildL2AuthHeaders("POST", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", path, nil, h, b)
}

type orderCreationPayload struct {
	TokenID string
	Side    Side
	Size    float64
	Price   float64
}

func rfqRequestOrderCreationPayload(q *RfqQuote) (orderCreationPayload, error) {
	switch q.MatchType {
	case RfqMatchComplementary:
		var side Side
		if q.Side == "BUY" {
			side = SideSell
		} else {
			side = SideBuy
		}
		sizeStr := q.SizeOut
		if side == SideBuy {
			sizeStr = q.SizeIn
		}
		sz, err := strconv.ParseFloat(sizeStr, 64)
		if err != nil {
			return orderCreationPayload{}, err
		}
		return orderCreationPayload{TokenID: q.Token, Side: side, Size: sz, Price: q.Price}, nil
	case RfqMatchMint, RfqMatchMerge:
		side := SideBuy
		if q.Side != "BUY" {
			side = SideSell
		}
		sizeStr := q.SizeIn
		if side == SideSell {
			sizeStr = q.SizeOut
		}
		sz, err := strconv.ParseFloat(sizeStr, 64)
		if err != nil {
			return orderCreationPayload{}, err
		}
		return orderCreationPayload{
			TokenID: q.Complement,
			Side:    side,
			Size:    sz,
			Price:   1 - q.Price,
		}, nil
	default:
		return orderCreationPayload{}, fmt.Errorf("invalid RFQ match type %q", q.MatchType)
	}
}

// ApproveRfqOrder approves a quote (maker side).
func (c *Client) ApproveRfqOrder(payload ApproveOrderParams) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	quotes, err := c.GetRfqQuoterQuotes(&GetRfqQuotesParams{QuoteIDs: []string{payload.QuoteID}})
	if err != nil {
		return nil, err
	}
	if len(quotes.Data) == 0 {
		return nil, fmt.Errorf("RFQ quote not found")
	}
	rfq := quotes.Data[0]
	side := SideBuy
	if rfq.Side != "BUY" {
		side = SideSell
	}
	var size float64
	if rfq.Side == "BUY" {
		size, _ = strconv.ParseFloat(rfq.SizeIn, 64)
	} else {
		size, _ = strconv.ParseFloat(rfq.SizeOut, 64)
	}
	exp := payload.Expiration
	req := OrderRequest{
		TokenID: rfq.Token, Price: rfq.Price, Size: size, Side: side,
		Expiration: &exp,
	}
	signed, err := c.BuildSignedLimitOrder(req)
	if err != nil {
		return nil, err
	}
	approvePayload := map[string]any{
		"requestId":     payload.RequestID,
		"quoteId":       payload.QuoteID,
		"owner":         c.apiKeyCredentials.ApiKey,
		"salt":          signed.Salt.Int64(),
		"maker":         signed.Maker.Hex(),
		"signer":        signed.Signer.Hex(),
		"tokenId":       signed.TokenID.String(),
		"makerAmount":   signed.MakerAmount.String(),
		"takerAmount":   signed.TakerAmount.String(),
		"side":          string(side),
		"signatureType": int(signed.SignatureType),
		"timestamp":     signed.Timestamp.String(),
		"metadata":      signed.Metadata.Hex(),
		"builder":       signed.Builder.Hex(),
		"expiration":    signed.Expiration.String(),
		"signature":     encodeOrderSignatureHex(signed.Signature),
	}
	b, err := json.Marshal(approvePayload)
	if err != nil {
		return nil, err
	}
	path := PathRFQQuoteApprove
	h, err := c.buildL2AuthHeaders("POST", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", path, nil, h, b)
}
