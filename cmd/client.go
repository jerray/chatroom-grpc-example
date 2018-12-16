package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jerray/chatroom/pb"
	"github.com/jerray/chatroom/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var clientCmdHost string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start up chatroom client",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cc, err := grpc.DialContext(ctx, clientCmdHost, grpc.WithInsecure())
		handleInitError(err, "connect")
		defer cc.Close()

		client := pb.NewChatroomClient(cc)
		stream, err := client.CheckIn(context.Background())
		handleInitError(err, "client checkin")
		defer stream.CloseSend()

		waitc := make(chan struct{})

		// Receiving message from server
		go func() {
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					close(waitc)
					return
				}
				if err != nil {
					log.Fatalf("Failed to receive a note : %v", err)
				}

				if login := in.GetLogin(); login != nil {
					log.Println(login.GetName(), "logged in")
				} else if message := in.GetMessage(); message != nil {
					log.Println(message.GetFrom().GetName(), ":", message.GetContent())
				}
			}
		}()

		// Reading message from stdin and send to server
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())
				if text == "" {
					continue
				}

				event := utils.ParseInput(text)
				if event == nil {
					continue
				}
				err := stream.Send(event)
				if err != nil {
					log.Fatalf("Failed to send message to server: %v", err)
				}
			}
		}()
		log.Println("client started")

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	loop:
		for {
			select {
			case <-waitc:
				break loop
			case <-quit:
				break loop
			}
		}
		log.Println("client exited")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().StringVarP(&clientCmdHost, "host", "s", "127.0.0.1", "Server host address")
}
