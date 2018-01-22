package commands

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gp_upgrade/hub/configutils"
	pb "gp_upgrade/idl"
	"gp_upgrade/shellParsers"

	gpbackupUtils "github.com/greenplum-db/gpbackup/utils"
	"github.com/pkg/errors"
)

const (
	// todo generalize to any host
	address = "localhost"
	port    = "6416"
)

type MonitorCommand struct {
	Host       string
	Port       int
	User       string
	PrivateKey string
	SegmentID  int
}

func (cmd MonitorCommand) Execute([]string) error {
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure())
	if err != nil {
		gpbackupUtils.GetLogger().Fatal(err, "did not connect")
	}
	if err != nil {
		return errors.New(err.Error())
	}
	client := pb.NewCommandListenerClient(conn)
	defer conn.Close()

	return cmd.execute(client, &shellParsers.RealShellParser{}, os.Stdout)
}

func (cmd MonitorCommand) execute(client pb.CommandListenerClient, shellParser shellParsers.ShellParser, writer io.Writer) error {

	/* Use as ssh reference for later use? */
	//user := cmd.User
	//if user == "" {
	//	user, _, _ = utils.GetUser() // todo last arg is for error--bubble up that error here? with what message?
	//}
	//output, err := connector.ConnectAndExecute(cmd.Host, cmd.Port, user, "ps auxx | grep pg_upgrade")

	targetPort, err := readConfigForSegmentPort(cmd.SegmentID)
	if err != nil {
		return errors.New(err.Error())
	}

	reply, err := client.CheckUpgradeStatus(context.Background(), &pb.CheckUpgradeStatusRequest{})
	if err != nil {
		return errors.New(err.Error())
	}

	gpbackupUtils.GetLogger().Info("Command Listener responded: %s", reply.ProcessList)

	status := "active"

	if !shellParser.IsPgUpgradeRunning(targetPort, reply.ProcessList) {
		status = "inactive"
	}
	msg := fmt.Sprintf(`pg_upgrade state - %s {"segment_id":%d,"host":"%s"}`, status, cmd.SegmentID, cmd.Host)
	fmt.Fprintf(writer, "%s\n", msg)

	return nil
}

func readConfigForSegmentPort(segmentID int) (int, error) {
	var err error
	reader := configutils.Reader{}
	err = reader.Read()
	if err != nil {
		return -1, errors.New(err.Error())
	}
	targetPort := reader.GetPortForSegment(segmentID)
	if targetPort == -1 {
		return -1, fmt.Errorf("segment_id %d not known in this cluster configuration", segmentID)
	}

	return targetPort, nil
}
