package utils

import (
	"strings"
	"wstrade/config"
)

func ConvertToOkxSpotInstID(okxFuturesInstID string) string {
	// BTC-USDT-SWAP => BTC-USDT
	return strings.Replace(okxFuturesInstID, "-SWAP", "", -1)
}

func ConvertToBinanceInstID(okxFuturesInstID string) string {
	// BTC-USDT-SWAP => BTCUSDT
	return strings.Replace(okxFuturesInstID, "-USDT-SWAP", "USDT", -1)
}

func ConvertBinanceFuturesInstIDToOkxFuturesInstID(binanceFuturesInstID string) string {
	// BTCUSDT	=> BTC-USDT-SWAP
	return strings.Replace(binanceFuturesInstID, "USDT", "-USDT-SWAP", -1)
}

func ConvertBinanceFuturesInstIDToOkxSpotInstID(binanceFuturesInstID string) string {
	// BTCUSDT	=> BTC-USDT
	return strings.Replace(binanceFuturesInstID, "USDT", "-USDT", -1)
}

func ConvertToStdInstType(exchange config.Exchange, instType string) config.InstrumentType {
	if exchange == config.OkxExchange {
		switch instType {
		case "SWAP":
			return config.FuturesInstrument
		case "SPOT":
			return config.SpotInstrument
		}
	}
	return config.UnknownInstrument
}

func GenOkxFuturesInstID(exchange config.Exchange, instType config.InstrumentType, instID string) string {
	if exchange == config.OkxExchange {
		switch instType {
		case "SWAP":
			return instID
		case "SPOT":
			return instID + "-SWAP"
		}
	}
	return instID
}
