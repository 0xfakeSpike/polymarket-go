package polymarket

import "encoding/json"

// clobWSMaxReadBytes is the max size of one WebSocket read (see e.g. go-polymarket-ws WSClient).
const clobWSMaxReadBytes = 10 * 1024 * 1024

// unwrapJSONArrayWSMessage returns data unchanged, or the first element when data is a JSON array.
// Polymarket sometimes wraps a market event in a one-element array (same pattern as github.com/dragonhuntr/go-polymarket-ws).
func unwrapJSONArrayWSMessage(data []byte) []byte {
	if len(data) == 0 || data[0] != '[' {
		return data
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil || len(arr) == 0 {
		return data
	}
	return arr[0]
}
