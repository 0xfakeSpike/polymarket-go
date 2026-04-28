package polymarket

import (
	"net/url"
	"strconv"
	"time"
)

// Values encodes *MarketsParams as query string parameters (net/url.Values).
func (p *MarketsParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit > 0 {
		v.Add("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		v.Add("offset", strconv.Itoa(p.Offset))
	}
	if p.Order != "" {
		v.Add("order", p.Order)
	}
	v.Add("ascending", strconv.FormatBool(p.Ascending))
	if p.Closed != nil {
		v.Add("closed", strconv.FormatBool(*p.Closed))
	}
	if p.Archived != nil {
		v.Add("archived", strconv.FormatBool(*p.Archived))
	}
	if p.Slug != "" {
		v.Add("slug", p.Slug)
	}
	if p.EventID != "" {
		v.Add("event_id", p.EventID)
	}
	if p.TagID != "" {
		v.Add("tag_id", p.TagID)
	}
	return v
}

// Values encodes *MarketsKeysetParams as query string parameters.
func (p *MarketsKeysetParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit > 0 {
		v.Add("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		v.Add("order", p.Order)
	}
	v.Add("ascending", strconv.FormatBool(p.Ascending))
	if p.AfterCursor != "" {
		v.Add("after_cursor", p.AfterCursor)
	}
	if p.Locale != "" {
		v.Add("locale", p.Locale)
	}
	for _, id := range p.ID {
		v.Add("id", id)
	}
	for _, slug := range p.Slug {
		v.Add("slug", slug)
	}
	if p.Closed != nil {
		v.Add("closed", strconv.FormatBool(*p.Closed))
	}
	if p.Decimalized != nil {
		v.Add("decimalized", strconv.FormatBool(*p.Decimalized))
	}
	for _, clobID := range p.ClobTokenIDs {
		v.Add("clob_token_ids", clobID)
	}
	for _, condID := range p.ConditionIDs {
		v.Add("condition_ids", condID)
	}
	for _, qid := range p.QuestionIDs {
		v.Add("question_ids", qid)
	}
	for _, mm := range p.MarketMakerAddress {
		v.Add("market_maker_address", mm)
	}
	if p.LiquidityNumMin != nil {
		v.Add("liquidity_num_min", strconv.FormatFloat(*p.LiquidityNumMin, 'f', -1, 64))
	}
	if p.LiquidityNumMax != nil {
		v.Add("liquidity_num_max", strconv.FormatFloat(*p.LiquidityNumMax, 'f', -1, 64))
	}
	if p.VolumeNumMin != nil {
		v.Add("volume_num_min", strconv.FormatFloat(*p.VolumeNumMin, 'f', -1, 64))
	}
	if p.VolumeNumMax != nil {
		v.Add("volume_num_max", strconv.FormatFloat(*p.VolumeNumMax, 'f', -1, 64))
	}
	if p.StartDateMin != nil {
		v.Add("start_date_min", p.StartDateMin.Format(time.RFC3339))
	}
	if p.StartDateMax != nil {
		v.Add("start_date_max", p.StartDateMax.Format(time.RFC3339))
	}
	if p.EndDateMin != nil {
		v.Add("end_date_min", p.EndDateMin.Format(time.RFC3339))
	}
	if p.EndDateMax != nil {
		v.Add("end_date_max", p.EndDateMax.Format(time.RFC3339))
	}
	for _, tagID := range p.TagID {
		v.Add("tag_id", strconv.Itoa(tagID))
	}
	if p.RelatedTags != nil {
		v.Add("related_tags", strconv.FormatBool(*p.RelatedTags))
	}
	if p.TagMatch != "" {
		v.Add("tag_match", p.TagMatch)
	}
	if p.CYOM != nil {
		v.Add("cyom", strconv.FormatBool(*p.CYOM))
	}
	if p.RFQEnabled != nil {
		v.Add("rfq_enabled", strconv.FormatBool(*p.RFQEnabled))
	}
	if p.UMAResolutionStatus != "" {
		v.Add("uma_resolution_status", p.UMAResolutionStatus)
	}
	if p.GameID != "" {
		v.Add("game_id", p.GameID)
	}
	for _, t := range p.SportsMarketTypes {
		v.Add("sports_market_types", t)
	}
	if p.IncludeTag != nil {
		v.Add("include_tag", strconv.FormatBool(*p.IncludeTag))
	}
	return v
}

// Values encodes *EventsParams as query string parameters.
func (p *EventsParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit > 0 {
		v.Add("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		v.Add("offset", strconv.Itoa(p.Offset))
	}
	if p.Order != "" {
		v.Add("order", p.Order)
	}
	v.Add("ascending", strconv.FormatBool(p.Ascending))
	for _, id := range p.ID {
		v.Add("id", id)
	}
	for _, slug := range p.Slug {
		v.Add("slug", slug)
	}
	if p.Active != nil {
		v.Add("active", strconv.FormatBool(*p.Active))
	}
	if p.Closed != nil {
		v.Add("closed", strconv.FormatBool(*p.Closed))
	}
	if p.Archived != nil {
		v.Add("archived", strconv.FormatBool(*p.Archived))
	}
	if p.TagID != nil {
		v.Add("tag_id", strconv.Itoa(*p.TagID))
	}
	for _, tagID := range p.ExcludeTagID {
		v.Add("exclude_tag_id", strconv.Itoa(tagID))
	}
	if p.RelatedTags != nil {
		v.Add("related_tags", strconv.FormatBool(*p.RelatedTags))
	}
	if p.Featured != nil {
		v.Add("featured", strconv.FormatBool(*p.Featured))
	}
	if p.CYOM != nil {
		v.Add("cyom", strconv.FormatBool(*p.CYOM))
	}
	if p.Recurrence != "" {
		v.Add("recurrence", p.Recurrence)
	}
	if p.StartDateMin != nil {
		v.Add("start_date_min", p.StartDateMin.Format(time.RFC3339))
	}
	if p.StartDateMax != nil {
		v.Add("start_date_max", p.StartDateMax.Format(time.RFC3339))
	}
	if p.EndDateMin != nil {
		v.Add("end_date_min", p.EndDateMin.Format(time.RFC3339))
	}
	if p.EndDateMax != nil {
		v.Add("end_date_max", p.EndDateMax.Format(time.RFC3339))
	}
	if p.TagSlug != "" {
		v.Add("tag_slug", p.TagSlug)
	}
	return v
}

// Values encodes *GetEventParams as query string parameters.
func (p *GetEventParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.IncludeChat != nil {
		v.Add("include_chat", strconv.FormatBool(*p.IncludeChat))
	}
	if p.IncludeTemplate != nil {
		v.Add("include_template", strconv.FormatBool(*p.IncludeTemplate))
	}
	return v
}

// Values encodes *CommentsParams as query string parameters.
func (p *CommentsParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit > 0 {
		v.Add("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		v.Add("offset", strconv.Itoa(p.Offset))
	}
	if p.Order != "" {
		v.Add("order", p.Order)
	}
	v.Add("ascending", strconv.FormatBool(p.Ascending))
	if p.ParentEntityType != "" {
		v.Add("parent_entity_type", p.ParentEntityType)
	}
	if p.ParentEntityID != nil {
		v.Add("parent_entity_id", strconv.Itoa(*p.ParentEntityID))
	}
	if p.GetPositions != nil {
		v.Add("get_positions", strconv.FormatBool(*p.GetPositions))
	}
	if p.HoldersOnly != nil {
		v.Add("holders_only", strconv.FormatBool(*p.HoldersOnly))
	}
	return v
}

// Values encodes *SearchParams as query string parameters.
func (p *SearchParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Q != "" {
		v.Add("q", p.Q)
	}
	if p.Type != "" {
		v.Add("type", p.Type)
	}
	for _, preset := range p.Presets {
		if preset != "" {
			v.Add("presets", preset)
		}
	}
	if p.Page > 0 {
		v.Add("page", strconv.Itoa(p.Page))
	}
	if p.LimitPerType > 0 {
		v.Add("limit_per_type", strconv.Itoa(p.LimitPerType))
	}
	if p.Sort != "" {
		v.Add("sort", p.Sort)
	}
	// Omit when false so requests match browser/Gamma defaults (see search-v2 curl examples).
	if p.Ascending {
		v.Add("ascending", "true")
	}
	if p.Cache != nil {
		v.Add("cache", strconv.FormatBool(*p.Cache))
	}
	if p.EventsStatus != "" {
		v.Add("events_status", p.EventsStatus)
	}
	for _, tag := range p.EventsTag {
		v.Add("events_tag", tag)
	}
	if p.KeepClosedMarkets != nil {
		v.Add("keep_closed_markets", strconv.Itoa(*p.KeepClosedMarkets))
	}
	if p.SearchTags != nil {
		v.Add("search_tags", strconv.FormatBool(*p.SearchTags))
	}
	if p.SearchProfiles != nil {
		v.Add("search_profiles", strconv.FormatBool(*p.SearchProfiles))
	}
	if p.Recurrence != "" {
		v.Add("recurrence", p.Recurrence)
	}
	for _, tagID := range p.ExcludeTagID {
		v.Add("exclude_tag_id", strconv.Itoa(tagID))
	}
	if p.Optimized != nil {
		v.Add("optimized", strconv.FormatBool(*p.Optimized))
	}
	return v
}

// Values encodes *EventsPaginationParams as query string parameters.
func (p *EventsPaginationParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit > 0 {
		v.Add("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		v.Add("offset", strconv.Itoa(p.Offset))
	}
	if p.Order != "" {
		v.Add("order", p.Order)
	}
	v.Add("ascending", strconv.FormatBool(p.Ascending))
	if p.Active != nil {
		v.Add("active", strconv.FormatBool(*p.Active))
	}
	if p.Archived != nil {
		v.Add("archived", strconv.FormatBool(*p.Archived))
	}
	if p.Closed != nil {
		v.Add("closed", strconv.FormatBool(*p.Closed))
	}
	for _, tagSlug := range p.TagSlug {
		v.Add("tag_slug", tagSlug)
	}
	return v
}

// Values encodes *EventsKeysetParams as query string parameters for GET /events/keyset.
func (p *EventsKeysetParams) Values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit > 0 {
		v.Add("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		v.Add("order", p.Order)
	}
	v.Add("ascending", strconv.FormatBool(p.Ascending))
	if p.AfterCursor != "" {
		v.Add("after_cursor", p.AfterCursor)
	}
	if p.Locale != "" {
		v.Add("locale", p.Locale)
	}
	for _, id := range p.ID {
		v.Add("id", strconv.Itoa(id))
	}
	for _, slug := range p.Slug {
		v.Add("slug", slug)
	}
	if p.Closed != nil {
		v.Add("closed", strconv.FormatBool(*p.Closed))
	}
	if p.Live != nil {
		v.Add("live", strconv.FormatBool(*p.Live))
	}
	if p.Featured != nil {
		v.Add("featured", strconv.FormatBool(*p.Featured))
	}
	if p.CYOM != nil {
		v.Add("cyom", strconv.FormatBool(*p.CYOM))
	}
	if p.TitleSearch != "" {
		v.Add("title_search", p.TitleSearch)
	}
	if p.LiquidityMin != nil {
		v.Add("liquidity_min", strconv.FormatFloat(*p.LiquidityMin, 'f', -1, 64))
	}
	if p.LiquidityMax != nil {
		v.Add("liquidity_max", strconv.FormatFloat(*p.LiquidityMax, 'f', -1, 64))
	}
	if p.VolumeMin != nil {
		v.Add("volume_min", strconv.FormatFloat(*p.VolumeMin, 'f', -1, 64))
	}
	if p.VolumeMax != nil {
		v.Add("volume_max", strconv.FormatFloat(*p.VolumeMax, 'f', -1, 64))
	}
	if p.StartDateMin != nil {
		v.Add("start_date_min", p.StartDateMin.Format(time.RFC3339))
	}
	if p.StartDateMax != nil {
		v.Add("start_date_max", p.StartDateMax.Format(time.RFC3339))
	}
	if p.EndDateMin != nil {
		v.Add("end_date_min", p.EndDateMin.Format(time.RFC3339))
	}
	if p.EndDateMax != nil {
		v.Add("end_date_max", p.EndDateMax.Format(time.RFC3339))
	}
	if p.StartTimeMin != nil {
		v.Add("start_time_min", p.StartTimeMin.Format(time.RFC3339))
	}
	if p.StartTimeMax != nil {
		v.Add("start_time_max", p.StartTimeMax.Format(time.RFC3339))
	}
	for _, tagID := range p.TagID {
		v.Add("tag_id", strconv.Itoa(tagID))
	}
	if p.TagSlug != "" {
		v.Add("tag_slug", p.TagSlug)
	}
	for _, ex := range p.ExcludeTagID {
		v.Add("exclude_tag_id", strconv.Itoa(ex))
	}
	if p.RelatedTags != nil {
		v.Add("related_tags", strconv.FormatBool(*p.RelatedTags))
	}
	if p.TagMatch != "" {
		v.Add("tag_match", p.TagMatch)
	}
	for _, sid := range p.SeriesID {
		v.Add("series_id", strconv.Itoa(sid))
	}
	for _, gid := range p.GameID {
		v.Add("game_id", strconv.Itoa(gid))
	}
	if p.EventDate != nil {
		v.Add("event_date", p.EventDate.Format(time.RFC3339))
	}
	if p.EventWeek != nil {
		v.Add("event_week", strconv.Itoa(*p.EventWeek))
	}
	if p.FeaturedOrder != nil {
		v.Add("featured_order", strconv.FormatBool(*p.FeaturedOrder))
	}
	if p.Recurrence != "" {
		v.Add("recurrence", p.Recurrence)
	}
	for _, cb := range p.CreatedBy {
		v.Add("created_by", cb)
	}
	if p.ParentEventID != nil {
		v.Add("parent_event_id", strconv.Itoa(*p.ParentEventID))
	}
	if p.IncludeChildren != nil {
		v.Add("include_children", strconv.FormatBool(*p.IncludeChildren))
	}
	if p.PartnerSlug != "" {
		v.Add("partner_slug", p.PartnerSlug)
	}
	if p.IncludeChat != nil {
		v.Add("include_chat", strconv.FormatBool(*p.IncludeChat))
	}
	if p.IncludeTemplate != nil {
		v.Add("include_template", strconv.FormatBool(*p.IncludeTemplate))
	}
	if p.IncludeBestLines != nil {
		v.Add("include_best_lines", strconv.FormatBool(*p.IncludeBestLines))
	}
	return v
}
