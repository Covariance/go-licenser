// Copyright (c) 2021 Yandex LLC. All rights reserved.
// Author: Martynov Pavel <covariance@yandex-team.ru>

package main

import (
	"flag"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use: "licenser",
	}

	logger, _ = zap.NewProduction()
	log       = logger.Sugar()
)

func main() {
	flag.Parse()

	if err := rootCmd.Execute(); err != nil {
		log.Errorf("execution failed: %v", err)
	}
}
