/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/nexi-intra/koksmat-emit/cmd"
	"github.com/nexi-intra/koksmat-emit/config"
)

func main() {
	config.Setup()
	cmd.Execute()
}
