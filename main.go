package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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

func mapToCsv(addr string, tx map[string]map[string][]float64) {
	txInFile, err := os.OpenFile(fmt.Sprintf("%s_in.csv", addr), os.O_WRONLY|os.O_CREATE, 0777)
	txOutFile, err := os.OpenFile(fmt.Sprintf("%s_out.csv", addr), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		Logger.Fatalln(err)
	}
	writerIn := csv.NewWriter(txInFile)
	writerOut := csv.NewWriter(txOutFile)
	writerIn.Write([]string{"from", "to", "value", "count"})
	writerOut.Write([]string{"from", "to", "value", "count"})
	for k, v := range tx {
		if k != addr {
			for i, j := range v {
				writerIn.Write([]string{k, i, strconv.FormatFloat(j[0], 'f', 6, 64), strconv.FormatFloat(j[1], 'f', 3, 64)})
			}
		} else {
			for i, j := range v {
				writerOut.Write([]string{k, i, strconv.FormatFloat(j[0], 'f', 6, 64), strconv.FormatFloat(j[1], 'f', 3, 64)})
			}
		}
	}
	writerIn.Flush()
	writerOut.Flush()
	defer txOutFile.Close()
	defer txInFile.Close()
	Logger.Println("done")
}

//func drawImage(tx map[string]map[string][]float64) {
//
//}
