package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Response struct {
	Amount Amount `json:"amount"`
}

func (sm *SmokeTestRunner) CheckZRC20ReserveAndSupply() {
	fmt.Println("Checking ZRC20 Reserve and Supply")
	sm.checkEthTSSBalance()
	sm.checkBtcTSSBalance()
	sm.checkUsdtTSSBalance()
	sm.checkZetaTSSBalance()
}

func (sm *SmokeTestRunner) checkEthTSSBalance() {
	tssBal, err := sm.GoerliClient.BalanceAt(context.Background(), sm.TSSAddress, nil)
	if err != nil {
		panic(err)
	}
	zrc20Supply, err := sm.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if tssBal.Cmp(zrc20Supply) < 0 {
		panic(fmt.Sprintf("ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ", tssBal, zrc20Supply))
	}
	fmt.Printf("ETH: TSS balance (%d) >= ZRC20 TotalSupply (%d)\n", tssBal, zrc20Supply)
}

func (sm *SmokeTestRunner) checkBtcTSSBalance() {
	utxos, err := sm.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	var btcBalance float64
	for _, utxo := range utxos {
		if utxo.Address == sm.BTCTSSAddress.EncodeAddress() {
			btcBalance += utxo.Amount
		}
	}
	zrc20Supply, err := sm.BTCZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	// #nosec G701 smoketest - always in range
	if int64(btcBalance*1e8) < zrc20Supply.Int64() {
		// #nosec G701 smoketest - always in range
		panic(fmt.Sprintf("BTC: TSS Balance (%d) < ZRC20 TotalSupply (%d) ", int64(btcBalance*1e8), zrc20Supply))
	}
	// #nosec G701 smoketest - always in range
	fmt.Printf("BTC: Balance (%d) >= ZRC20 TotalSupply (%d)\n", int64(btcBalance*1e8), zrc20Supply)
}

func (sm *SmokeTestRunner) checkUsdtTSSBalance() {
	usdtBal, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, sm.ERC20CustodyAddr)
	if err != nil {
		panic(err)
	}
	zrc20Supply, err := sm.USDTZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if usdtBal.Cmp(zrc20Supply) < 0 {
		panic(fmt.Sprintf("USDT: TSS balance (%d) < ZRC20 TotalSupply (%d) ", usdtBal, zrc20Supply))
	}
	fmt.Printf("USDT: TSS balance (%d) >= ZRC20 TotalSupply (%d)\n", usdtBal, zrc20Supply)
}

func (sm *SmokeTestRunner) checkZetaTSSBalance() {
	zetaLocked, err := sm.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	resp, err := http.Get("http://zetacore0:1317/cosmos/bank/v1beta1/supply/by_denom?denom=azeta")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	zetaSupply, _ := big.NewInt(0).SetString(result.Amount.Amount, 10)
	if zetaLocked.Cmp(zetaSupply) < 0 {
		fmt.Printf(fmt.Sprintf("ZETA: TSS balance (%d) < ZRC20 TotalSupply (%d) \n", zetaLocked, zetaSupply))
	} else {
		fmt.Printf("ZETA: TSS balance (%d) >= ZRC20 TotalSupply (%d)\n", zetaLocked, zetaSupply)
	}
}
