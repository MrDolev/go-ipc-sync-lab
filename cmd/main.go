package main

import (
	pd "go-ipc/pkg/prodcons"
)

func main() {
	var input_data []any = []any{10, 10, 10, 31, 41}
	var pd pd.ServiceRunnerI = pd.NewProdCons(input_data)
	pd.Runner()
}
