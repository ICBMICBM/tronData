package main

import (
	"encoding/csv"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"strconv"
	"time"
	"tronData/trongrid"
)

var Logger = log.New(os.Stdout, "tronData:", log.Lshortfile)

func main() {
	var addr string
	Logger.Printf("Input a tron address:")
	_, _ = fmt.Scanln(&addr)
	data := trongrid.GetAllTrc20Tx(addr)
	t := trongrid.MergeTx(data)
	mapToCsv(addr, t)
}

func tsToString(ts int64) string {
	t := time.Unix(ts/1000, 0)
	return t.Format(time.DateTime)
}

func mapToCsv(addr string, tx map[string]map[string]trongrid.TxDetail) {
	txInFile, err := os.OpenFile(fmt.Sprintf("%s_in.csv", addr), os.O_WRONLY|os.O_CREATE, 0777)
	txOutFile, err := os.OpenFile(fmt.Sprintf("%s_out.csv", addr), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		Logger.Fatalln(err)
	}
	writerIn := csv.NewWriter(txInFile)
	writerOut := csv.NewWriter(txOutFile)
	writerIn.Write([]string{"from", "to", "value", "count", "minTimestamp", "maxTimestamp"})
	writerOut.Write([]string{"from", "to", "value", "count", "minTimestamp", "maxTimestamp"})
	for from, tos := range tx {
		for to := range tos {
			var detail trongrid.TxDetail
			detail = tx[from][to]
			var minTS = lo.MinBy(detail.TxTimestamp, func(a int64, b int64) bool {
				return a < b
			})
			var maxTS = lo.MaxBy(detail.TxTimestamp, func(a int64, b int64) bool {
				return a > b
			})
			if from == addr {
				writerOut.Write([]string{from, to, strconv.FormatFloat(detail.Total, 'f', 6, 64), strconv.FormatUint(detail.TxCount, 10), tsToString(minTS), tsToString(maxTS)})
			} else {
				writerIn.Write([]string{from, to, strconv.FormatFloat(detail.Total, 'f', 6, 64), strconv.FormatUint(detail.TxCount, 10), tsToString(minTS), tsToString(maxTS)})
			}
		}
	}
	writerIn.Flush()
	writerOut.Flush()
	defer txOutFile.Close()
	defer txInFile.Close()
	Logger.Println("done")
}
