// Package polymarket implements the Polymarket Gamma (markets, search), Data API, and CLOB client.
//
// Prefer importing the module root:
//
//	import "github.com/0xfakeSpike/polymarket-go"
//
// The polymarket/ path provides the same [Client] for existing importers.
//
// [NewClient] obtains L2 API credentials by default; use [WithSkipL2APIKeyBootstrap] or
// [WithAPIKeyCredentials] to change that. Default order signature type is POLY_GNOSIS_SAFE (wire value 2).
// Set the Safe maker with [WithPolymarketSafeMaker] or [WithFunderAddress]; EOA flows use [WithSignatureType]
// with go-order-utils model.EOA. For a fixed neg-risk verifying contract, use [WithForceNegRiskSigning] or
// [Client.SetForceNegRiskSigning].
//
// Gamma and Data hosts: [Client.GammaGET], [Client.GammaPOST], [Client.DataGET], [Client.DataPOST] (paths start with "/").
package polymarket
