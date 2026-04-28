// Package polymarket implements Polymarket Gamma (markets, search) and CLOB trading APIs,
// aligned with @polymarket/clob-client. [NewClient] obtains L2 API credentials by default;
// use [WithSkipL2APIKeyBootstrap] or [WithAPIKeyCredentials] to override.
//
// Deprecated: for new applications, import the root package
// "github.com/0xfakeSpike/polymarket-go". This subpackage remains as a compatibility path.
// Default order signature type is POLY_GNOSIS_SAFE (wire value 2), matching polymarket.com POST /order.
// You must set the Safe maker with [WithPolymarketSafeMaker] or [WithFunderAddress]; pure EOA trading
// uses [WithSignatureType] with go-order-utils model.EOA.
//
// To always sign against the neg-risk exchange contract (fixed verifyingContract), use
// [WithForceNegRiskSigning] or [Client.SetForceNegRiskSigning].
//
// Non-CLOB HTTP: use [Client.GammaGET]/[Client.GammaPOST] for the Gamma API host and
// [Client.DataGET]/[Client.DataPOST] for the Data API host (paths must start with "/").
package polymarket
