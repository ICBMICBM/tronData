package trongrid

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var Logger = log.New(os.Stdout, "trongrid:", log.Lshortfile)

type TxDetail struct {
	Total       float64
	TxCount     uint64
	TxTimestamp []int64
}

var apiKeys = []string{
	"92f3bce0-8bd6-4679-934b-89bf8cceedc6",
	"4e7eec01-3806-41d8-8bae-77510dbe2922",
	"8fa8391d-ed31-4ac7-bb8c-26a15beec487",
	"6d4a070b-3e53-42ce-aae1-f1bd63fbd472",
	"56318519-9b06-4e46-a281-748118034749",
	"1c92b3f9-d274-44fa-878d-8d4c7e1b549e",
	"7c248b6d-d2a3-47a5-aa29-eec7c571e077",
	"a6f8c516-8e5a-42f4-ac5c-f28b7d103cf7",
	"0f50155b-1506-4c74-9e0d-f6b02473b19b",
	"7ff2590e-6f3a-48fe-bd4d-5905911bfa0f",
	"2c5b7dbe-c983-41e3-8d9e-6004f70f18fd",
	"2b660cdd-9951-4436-8664-8ec745021f63",
}

func randomApiKey() string {
	return apiKeys[rand.Int()%len(apiKeys)]
}

func GetTrc20Tx(addr string, timestamp uint64) (string, uint64, error) {
	url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions/trc20?only_confirmed=true&contract_address=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t&limit=200", addr)
	if timestamp > 0 {
		url = fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions/trc20?only_confirmed=true&contract_address=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t&limit=200&max_timestamp=%d", addr, timestamp)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("TRON-PRO-API-KEY", randomApiKey())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer func(Body io.ReadCloser) error {
		err := Body.Close()
		if err != nil {
			return err
		}
		return nil
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}
	jsonString := string(body)
	Logger.Println(gjson.Get(jsonString, "data.@reverse.0.block_timestamp").Uint())
	return jsonString, gjson.Get(jsonString, "data.@reverse.0.block_timestamp").Uint(), nil
}

func GetAllTrc20Tx(addr string) []string {
	var res []string
	var nxt uint64
	for {
		j, n, err := GetTrc20Tx(addr, nxt)

		if err != nil {
			continue
		} else {
			if nxt == n {
				break
			} else {
				res = append(res, j)
				nxt = n
			}
		}
	}
	return res
}

func MergeTx(data []string) map[string]map[string]TxDetail {
	var used = map[string]bool{}
	var txs = map[string]map[string]TxDetail{}
	for _, d := range data {
		dataBlock := gjson.Get(d, "data").Array()
		for _, t := range dataBlock {
			singleTx := t.Map()
			txHash, from, to, value := singleTx["transaction_id"].String(), singleTx["from"].String(), singleTx["to"].String(), singleTx["value"].Float()*0.000001
			if b := used[txHash]; b {
				continue
			}
			if singleTx["type"].String() == "Transfer" {
				if _, b := txs[from]; !b {
					txs[from] = map[string]TxDetail{to: {
						Total:       value,
						TxCount:     1,
						TxTimestamp: []int64{singleTx["block_timestamp"].Int()},
					}}
				} else {
					if _, b := txs[from][to]; !b {
						txs[from][to] = TxDetail{
							Total:       value,
							TxCount:     1,
							TxTimestamp: []int64{singleTx["block_timestamp"].Int()},
						}
					} else {
						var detail = txs[from][to]
						var newDetail = TxDetail{
							Total:       detail.Total + value,
							TxCount:     detail.TxCount + 1,
							TxTimestamp: append(detail.TxTimestamp, singleTx["block_timestamp"].Int()),
						}
						txs[from][to] = newDetail
					}
				}
				used[txHash] = true
			} else {
				continue
			}
		}
	}
	return txs
}
