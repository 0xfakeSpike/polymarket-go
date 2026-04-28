package polymarket

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/polymarket/go-order-utils/pkg/model"
)

// ROUNDING_CONFIG mirrors clob-client/src/order-builder/helpers.ts
var roundingConfig = map[string]roundConfig{
	"0.1":    {price: 1, size: 2, amount: 3},
	"0.01":   {price: 2, size: 2, amount: 4},
	"0.001":  {price: 3, size: 2, amount: 5},
	"0.0001": {price: 4, size: 2, amount: 6},
}

type roundConfig struct {
	price, size, amount int
}

func decimalPlaces(num float64) int {
	if num == float64(int64(num)) {
		return 0
	}
	s := strconv.FormatFloat(num, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i < 0 {
		return 0
	}
	return len(s) - i - 1
}

func roundNormal(num float64, decimals int) float64 {
	if decimalPlaces(num) <= decimals {
		return num
	}
	return math.Round((num+1e-12)*math.Pow10(decimals)) / math.Pow10(decimals)
}

func roundDown(num float64, decimals int) float64 {
	if decimalPlaces(num) <= decimals {
		return num
	}
	p := math.Pow10(decimals)
	return math.Floor(num*p+1e-12) / p
}

func roundUp(num float64, decimals int) float64 {
	if decimalPlaces(num) <= decimals {
		return num
	}
	p := math.Pow10(decimals)
	return math.Ceil(num*p-1e-12) / p
}

func getOrderRawAmounts(side Side, size, price float64, cfg roundConfig) (utilsSide model.Side, rawMakerAmt, rawTakerAmt float64) {
	rawPrice := roundNormal(price, cfg.price)

	if side == SideBuy {
		rawTakerAmt = roundDown(size, cfg.size)
		rawMakerAmt = rawTakerAmt * rawPrice
		if decimalPlaces(rawMakerAmt) > cfg.amount {
			rawMakerAmt = roundUp(rawMakerAmt, cfg.amount+4)
			if decimalPlaces(rawMakerAmt) > cfg.amount {
				rawMakerAmt = roundDown(rawMakerAmt, cfg.amount)
			}
		}
		return model.BUY, rawMakerAmt, rawTakerAmt
	}

	rawMakerAmt = roundDown(size, cfg.size)
	rawTakerAmt = rawMakerAmt * rawPrice
	if decimalPlaces(rawTakerAmt) > cfg.amount {
		rawTakerAmt = roundUp(rawTakerAmt, cfg.amount+4)
		if decimalPlaces(rawTakerAmt) > cfg.amount {
			rawTakerAmt = roundDown(rawTakerAmt, cfg.amount)
		}
	}
	return model.SELL, rawMakerAmt, rawTakerAmt
}

const collateralTokenDecimals = 6

func parseUnitsHuman(amount float64, decimals int) (string, error) {
	s := strconv.FormatFloat(amount, 'f', -1, 64)
	rat := new(big.Rat)
	if _, ok := rat.SetString(s); !ok {
		return "", fmt.Errorf("parse amount %q", s)
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	rat.Mul(rat, new(big.Rat).SetInt(scale))
	if !rat.IsInt() {
		num := rat.Num()
		den := rat.Denom()
		q := new(big.Int).Quo(num, den)
		return q.String(), nil
	}
	return rat.Num().String(), nil
}

func getMarketOrderRawAmounts(side Side, amount, price float64, cfg roundConfig) (model.Side, float64, float64) {
	rawPrice := roundDown(price, cfg.price)
	if side == SideBuy {
		rawMakerAmt := roundDown(amount, cfg.size)
		rawTakerAmt := rawMakerAmt / rawPrice
		if decimalPlaces(rawTakerAmt) > cfg.amount {
			rawTakerAmt = roundUp(rawTakerAmt, cfg.amount+4)
			if decimalPlaces(rawTakerAmt) > cfg.amount {
				rawTakerAmt = roundDown(rawTakerAmt, cfg.amount)
			}
		}
		return model.BUY, rawMakerAmt, rawTakerAmt
	}
	rawMakerAmt := roundDown(amount, cfg.size)
	rawTakerAmt := rawMakerAmt * rawPrice
	if decimalPlaces(rawTakerAmt) > cfg.amount {
		rawTakerAmt = roundUp(rawTakerAmt, cfg.amount+4)
		if decimalPlaces(rawTakerAmt) > cfg.amount {
			rawTakerAmt = roundDown(rawTakerAmt, cfg.amount)
		}
	}
	return model.SELL, rawMakerAmt, rawTakerAmt
}

type orderAmounts struct {
	MakerAmount string
	TakerAmount string
}

func buildMarketOrderCreationArgs(
	req MarketOrderRequest,
	price float64,
	cfg roundConfig,
) (*orderAmounts, error) {
	side, rawMaker, rawTaker := getMarketOrderRawAmounts(req.Side, req.Amount, price, cfg)
	_ = side
	makerAmt, err := parseUnitsHuman(rawMaker, collateralTokenDecimals)
	if err != nil {
		return nil, err
	}
	takerAmt, err := parseUnitsHuman(rawTaker, collateralTokenDecimals)
	if err != nil {
		return nil, err
	}
	return &orderAmounts{MakerAmount: makerAmt, TakerAmount: takerAmt}, nil
}

func buildOrderCreationArgs(
	user OrderRequest,
	cfg roundConfig,
) (*orderAmounts, error) {
	side, rawMaker, rawTaker := getOrderRawAmounts(user.Side, user.Size, user.Price, cfg)
	_ = side
	makerAmt, err := parseUnitsHuman(rawMaker, collateralTokenDecimals)
	if err != nil {
		return nil, err
	}
	takerAmt, err := parseUnitsHuman(rawTaker, collateralTokenDecimals)
	if err != nil {
		return nil, err
	}
	return &orderAmounts{MakerAmount: makerAmt, TakerAmount: takerAmt}, nil
}
