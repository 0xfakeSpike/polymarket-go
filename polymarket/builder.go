package polymarket

// BuilderSigner produces extra POLY_BUILDER_* headers for order flow / builder endpoints.
// Implement using @polymarket/builder-signing-sdk logic or your own HMAC scheme.
type BuilderSigner interface {
	SignBuilder(method, requestPath, body string) (headers map[string]string, err error)
}
