// add_white.go
package cmd

import (
	"context"
	"log"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var whiteNet string

var addWhiteCmd = &cobra.Command{
	Use:   "add-white",
	Short: "Добавить сеть в белый список",
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
		_, err = client.AddToWhitelist(context.Background(), &pb.NetworkRequest{Network: whiteNet})
		if err != nil {
			log.Fatalf("add-white failed: %v", err)
		}
	},
}

func init() {
	addWhiteCmd.Flags().StringVar(&whiteNet, "network", "", "CIDR-сеть")
	addWhiteCmd.MarkFlagRequired("network")
	rootCmd.AddCommand(addWhiteCmd)
}
