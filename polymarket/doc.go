// Package polymarket implements a Go client for the Polymarket CLOB API.
//
// Prefer importing the module root:
//
//	import "github.com/0xfakeSpike/polymarket-go"
//
// The polymarket/ path provides the same [Client] for existing importers.
//
// [NewClient] obtains L2 API credentials by default; use [WithSkipL2APIKeyBootstrap] or
// [WithAPIKeyCredentials] to change that. Default order signature type is POLY_GNOSIS_SAFE
// (wire value 2). Set the Safe maker with [WithPolymarketSafeMaker] or [WithFunderAddress];
// EOA flows use [WithSignatureType] with go-order-utils model.EOA. For a fixed neg-risk
// verifying contract, use [WithForceNegRiskSigning] or [Client.SetForceNegRiskSigning].
package polymarket
