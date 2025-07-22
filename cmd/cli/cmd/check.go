package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	login    string
	password string
	ip       string
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Проверить авторизацию",
	Run: func(_ *cobra.Command, _ []string) {
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("gRPC connect error: %v", err)
			return
		}
		defer conn.Close()

		client := pb.NewAntiBruteforceClient(conn)
		resp, err := client.Check(context.Background(), &pb.CheckRequest{
			Login:    login,
			Password: password,
			Ip:       ip,
		})
		if err != nil {
			log.Printf("check failed: %v", err)
			return
		}

		fmt.Println("OK:", resp.Ok)
	},
}

func init() {
	checkCmd.Flags().StringVar(&login, "login", "", "логин")
	checkCmd.Flags().StringVar(&password, "password", "", "пароль")
	checkCmd.Flags().StringVar(&ip, "ip", "", "IP-адрес")
	checkCmd.MarkFlagRequired("login")
	checkCmd.MarkFlagRequired("password")
	checkCmd.MarkFlagRequired("ip")
	rootCmd.AddCommand(checkCmd)
}
