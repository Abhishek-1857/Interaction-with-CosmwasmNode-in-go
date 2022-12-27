package main

import (
	"context"
	"fmt"
	query "interact/query"
	transaction "interact/transaction"
)

func main() {
	queryres:=query.QueryState()
	txres:=transaction.Transaction(context.Background())
	fmt.Println(queryres)
	fmt.Println(txres)
}
