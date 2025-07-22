package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var addr string

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI клиент для gRPC сервиса Anti-Bruteforce",
	Long:  "CLI для взаимодействия с gRPC-сервисом защиты от перебора.",
}

func Execute() {
	// Глобальный флаг --addr
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "localhost:8081", "адрес gRPC сервера")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
