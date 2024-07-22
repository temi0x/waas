package wallet

func GetContractAddress(chain string, tokenName string) (contractAddress string) {
	switch chain {
	case "ETH":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0xdac17f958d2ee523a2206206994597c13d831ec7"
		case "USDC":
			return "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"

		default:
			return ""
		}
	case "BASE":
		switch tokenName {
		case "ETH":
			return "0x"
		// case "USDT":
		// 	return "0x55d398326f99059ff775485246999027b3197955"
		case "USDC":
			return "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
		default:
			return ""
		}
	case "BSC":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0x55d398326f99059ff775485246999027b3197955"
		case "USDC":
			return "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d"

		default:
			return ""
		}
	case "POLYGON":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0xc2132d05d31c914a87c6611c10748aeb04b58e8f"
		case "USDC":
			return "0x"

		default:
			return ""
		}
	case "FANTOM":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0x049d68029688eabf473097a2fc38ef61633a3c7a"
		case "USDC":
			return "0x04068DA6C83AFCFA0e13ba15A6696662335D5B75"

		default:
			return ""
		}
	case "XDAI":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0x4ECaBa5870353805a9F068101A40E0f32ed605C6"
		case "USDC":
			return "0x"

		default:
			return ""
		}
	case "AVALANCHE":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0x9702230A8Ea53601f5cD2dc00fDBc13d4dF4A8c7"
		case "USDC":
			return "0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E"

		default:
			return ""

		}
	case "ARBITRUM":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9"
		case "USDC":
			return "0xaf88d065e77c8cC2239327C5EDb3A432268e5831"

		default:
			return ""
		}

	case "OPTIMISM":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0x94b008aa00579c1307b0ef2c499ad98a8ce58e58"
		case "USDC":
			return "0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85"

		default:
			return ""

		}

	case "HECO":
		switch tokenName {
		case "ETH":
			return "0x"
		case "USDT":
			return "0xa71edc38d189767582c38a3145b5873052c3e47a"
		case "USDC":
			return "0x985458e523db3d53125813ed68c274899e9dfab4"

		default:
			return ""
		}
	}

	return ""
}
