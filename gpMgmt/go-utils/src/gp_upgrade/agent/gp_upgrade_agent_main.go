package main

import (
	"os"

	gpbackupUtils "github.com/greenplum-db/gpbackup/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gp_upgrade/agent/services"
	pb "gp_upgrade/idl"
	"net"
)

const (
	port = ":6416"
)

func main() {
	//debug.SetTraceback("all")
	//parser := flags.NewParser(&AllServices, flags.HelpFlag|flags.PrintErrors)
	//
	//_, err := parser.Parse()
	//if err != nil {
	//	os.Exit(utils.GetExitCodeForError(err))
	//}
	var logdir string
	var RootCmd = &cobra.Command{
		Use:   "command_listener --log-directory [path]",
		Short: "Start the Command Listener (blocks)",
		Long:  `Start the Command Listener (blocks)`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if logdir == "" {
				return errors.New("the required flag '--log-directory' was not specified")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do Stuff Here
			gpbackupUtils.InitializeLogging("gp_upgrade_agent", logdir)
			errorChannel := make(chan error)
			defer close(errorChannel)
			lis, err := net.Listen("tcp", port)
			if err != nil {
				gpbackupUtils.GetLogger().Fatal(err, "failed to listen")
				return err
			}

			server := grpc.NewServer()
			myImpl := services.NewCommandListener()
			pb.RegisterCommandListenerServer(server, myImpl)
			reflection.Register(server)
			go func(myListener net.Listener) {
				if err := server.Serve(myListener); err != nil {
					gpbackupUtils.GetLogger().Fatal(err, "failed to serve", err)
					errorChannel <- err
				}

				close(errorChannel)
			}(lis)

			select {
			case err := <-errorChannel:
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	RootCmd.PersistentFlags().StringVar(&logdir, "log-directory", "", "command_listener log directory")

	if err := RootCmd.Execute(); err != nil {
		gpbackupUtils.GetLogger().Error(err.Error())
		os.Exit(1)
	}
}
