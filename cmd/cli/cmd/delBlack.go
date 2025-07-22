// del_black.go
package cmd

import (
	"context"
	"log"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var noBlackNet string

var delBlackCmd = &cobra.Command{
	Use:   "del-black",
	Short: "Удалить сеть из черного списка",
	Run: func(_ *cobra.Command, _ []string) {
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("grpc dial: %v", err)
		}
		defer conn.Close()

		client := pb.NewAntiBruteforceClient(conn)
		_, err = client.RemoveFromBlacklist(context.Background(), &pb.NetworkRequest{Network: noBlackNet})
		if err != nil {
			log.Fatalf("del-black failed: %v", err)
		}
	},
}

func init() {
	delBlackCmd.Flags().StringVar(&noBlackNet, "network", "", "CIDR-сеть")
	delBlackCmd.MarkFlagRequired("network")
	rootCmd.AddCommand(delBlackCmd)
}
