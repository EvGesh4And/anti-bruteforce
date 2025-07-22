package cmd

import (
	"context"
	"log"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	resetLogin string
	resetIP    string
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Сбросить лимиты",
	Run: func(_ *cobra.Command, _ []string) {
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("grpc dial failed: %v", err)
		}
		defer conn.Close()

		client := pb.NewAntiBruteforceClient(conn)
		_, err = client.Reset(context.Background(), &pb.ResetRequest{
			Login: resetLogin,
			Ip:    resetIP,
		})
		if err != nil {
			log.Fatalf("reset failed: %v", err)
		}
	},
}

func init() {
	resetCmd.Flags().StringVar(&resetLogin, "login", "", "логин")
	resetCmd.Flags().StringVar(&resetIP, "ip", "", "IP-адрес")
	resetCmd.MarkFlagRequired("login")
	resetCmd.MarkFlagRequired("ip")
	rootCmd.AddCommand(resetCmd)
}
