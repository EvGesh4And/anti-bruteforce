// add_black.go
package cmd

import (
	"context"
	"log"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var blackNet string

var addBlackCmd = &cobra.Command{
	Use:   "add-black",
	Short: "Добавить сеть в черный список",
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
		_, err = client.AddToBlacklist(context.Background(), &pb.NetworkRequest{Network: blackNet})
		if err != nil {
			log.Printf("add-black failed: %v", err)
			return
		}
	},
}

func init() {
	addBlackCmd.Flags().StringVar(&blackNet, "network", "", "CIDR-сеть")
	addBlackCmd.MarkFlagRequired("network")
	rootCmd.AddCommand(addBlackCmd)
}
