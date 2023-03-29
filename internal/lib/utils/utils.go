package utils

import (
	"fmt"
	"os"
	"strconv"
)

func EnvVar(k string, d any) (res any) {
	if val, ok := os.LookupEnv(k); ok {
		switch fmt.Sprintf("%T", d) {
		case "string":
			res = val
		case "int":
			i, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("Incorrect parameter!")
				os.Exit(1)
			}
			res = i
		default:
			fmt.Println("Incorrect parameter!")
			os.Exit(1)
		}
	} else {
		res = d
	}
	return
}
