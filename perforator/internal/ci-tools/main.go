package main

import (
	"context"
	"fmt"
	"os"
)


func main() {
	err := mainImpl(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
