package binance

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/models"
)

const nativeTransferTransaction = `
{
	"blockHeight": 7761368,
	"code": 0,
	"confirmBlocks": 2089441,
	"fromAddr": "tbnb1fhr04azuhcj0dulm7ka40y0cqjlafwae9k9gk2",
	"hasChildren": 0,
	"log": "Msg 0: ",
	"mappedTxAsset": "BNB",
	"timeStamp": 1555049867552,
	"toAddr": "tbnb1sylyjw032eajr9cyllp26n04300qzzre38qyv5",
	"txAge": 836729,
	"txAsset": "BNB",
	"txFee": 0.00125,
	"txHash": "1681EE543FB4B5A628EF21D746E031F018E226D127044A4F9BA5EE2542A44555",
	"txType": "TRANSFER",
	"value": 100000,
	"memo": "test"
}`

const nativeTokenTransferTransaction = `
{
	"blockHeight": 7928667,
	"code": 0,
	"confirmBlocks": 1922024,
	"fromAddr": "tbnb1ttyn4csghfgyxreu7lmdu3lcplhqhxtzced45a",
	"hasChildren": 0,
	"log": "Msg 0: ",
	"mappedTxAsset": "YLC",
	"timeStamp": 1555117625829,
	"toAddr": "tbnb12hlquylu78cjylk5zshxpdj6hf3t0tahwjt3ex",
	"txAge": 768924,
	"txAsset": "YLC-D8B",
	"txFee": 0.00125,
	"txHash": "95CF63FAA27579A9B6AF84EF8B2DFEAC29627479E9C98E7F5AE4535E213FA4C9",
	"txType": "TRANSFER",
	"value": 2.10572645,
	"memo": "test"
}`

var transferDst = models.Tx{
	ID:     "1681EE543FB4B5A628EF21D746E031F018E226D127044A4F9BA5EE2542A44555",
	Coin:   coin.BNB,
	From:   "tbnb1fhr04azuhcj0dulm7ka40y0cqjlafwae9k9gk2",
	To:     "tbnb1sylyjw032eajr9cyllp26n04300qzzre38qyv5",
	Fee:    "125000",
	Date:   1555049867,
	Block:  7761368,
	Status: models.StatusCompleted,
	Memo:   "test",
	Meta: models.Transfer{
		Value: "10000000000000",
	},
}

var nativeTransferDst = models.Tx{
	ID:     "95CF63FAA27579A9B6AF84EF8B2DFEAC29627479E9C98E7F5AE4535E213FA4C9",
	Coin:   coin.BNB,
	From:   "tbnb1ttyn4csghfgyxreu7lmdu3lcplhqhxtzced45a",
	To:     "tbnb12hlquylu78cjylk5zshxpdj6hf3t0tahwjt3ex",
	Fee:    "125000",
	Date:   1555117625,
	Block:  7928667,
	Status: models.StatusCompleted,
	Memo:   "test",
	Meta: models.NativeTokenTransfer{
		TokenID:  "YLC-D8B",
		Symbol:   "YLC",
		Value:    "210572645",
		Decimals: 8,
		From:     "tbnb1ttyn4csghfgyxreu7lmdu3lcplhqhxtzced45a",
		To:       "tbnb12hlquylu78cjylk5zshxpdj6hf3t0tahwjt3ex",
	},
}

type test struct {
	name        string
	apiResponse string
	expected    *models.Tx
	token       string
	address     string
}

func TestNormalize(t *testing.T) {
	testNormalize(t, &test{
		name:        "transfer",
		apiResponse: nativeTransferTransaction,
		expected:    &transferDst,
		token:       "",
		address:     "",
	})
	testNormalize(t, &test{
		name:        "native token transfer",
		apiResponse: nativeTokenTransferTransaction,
		expected:    &nativeTransferDst,
		token:       "YLC-D8B",
		address:     "",
	})
}

func testNormalize(t *testing.T, _test *test) {
	var srcTx Tx
	err := json.Unmarshal([]byte(_test.apiResponse), &srcTx)
	if err != nil {
		t.Error(err)
		return
	}

	tx, ok := Normalize(&srcTx, _test.token, _test.address)
	if !ok {
		t.Errorf("Binance: Transfer could not be normalized")
	}

	resJSON, err := json.Marshal(&tx)
	if err != nil {
		t.Fatal(err)
	}

	dstJSON, err := json.Marshal(_test.expected)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(resJSON, dstJSON) {
		println(string(resJSON))
		println(string(dstJSON))
		t.Error("Binance: transactions not equal")
	}
}
