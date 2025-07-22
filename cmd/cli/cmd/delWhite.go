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

var noWhiteNet string

var delWhiteCmd = &cobra.Command{
	Use:   "del-white",
	Short: "Удалить сеть из белого списка",
	Run: func(_ *cobra.Command, _ []string) {
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("grpc dial: %v", err)
			return
		}
		defer conn.Close()

		client := pb.NewAntiBruteforceClient(conn)
		_, err = client.RemoveFromWhitelist(context.Background(), &pb.NetworkRequest{Network: noWhiteNet})
		if err != nil {
			log.Printf("del-white failed: %v", err)
			return
		}
	},
}

func init() {
	delWhiteCmd.Flags().StringVar(&noWhiteNet, "network", "", "CIDR-сеть")
	delWhiteCmd.MarkFlagRequired("network")
	rootCmd.AddCommand(delWhiteCmd)
}
