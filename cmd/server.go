package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jerray/chatroom/controllers"
	"github.com/jerray/chatroom/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var serverCmdPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start up chatroom server",
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCmdPort))
		handleInitError(err, "net")

		gs := grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Time: 10 * time.Minute,
			}),

			// Register stream middleware.
			grpc.StreamInterceptor(controllers.ClientIDSetter),
		)

		api := controllers.NewRoomController()
		pb.RegisterChatroomServer(gs, api)
		go gs.Serve(lis)

		log.Println("server started")

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		// Finish after all clients disconnected
		gs.GracefulStop()
		log.Println("server exited")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().IntVarP(&serverCmdPort, "port", "p", 3000, "Port to listen")
}
