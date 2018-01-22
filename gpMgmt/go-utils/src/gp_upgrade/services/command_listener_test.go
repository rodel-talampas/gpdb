package services_test

import (
	"io"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gp_upgrade/idl"
	"gp_upgrade/services"
)

var _ = Describe("CommandListenerManager", func() {
	var (
		manager    idl.CommandListenerServer
		server     *grpc.Server
		connCloser io.Closer
		clsClient  idl.CommandListenerClient
	)

	var startGRPCServer = func(cls idl.CommandListenerServer) (*grpc.Server, string) {
		lis, err := net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())
		s := grpc.NewServer()
		idl.RegisterCommandListenerServer(s, cls)
		go s.Serve(lis)

		return s, lis.Addr().String()
	}

	var establishClient = func(clsAddr string) (idl.CommandListenerClient, io.Closer) {
		conn, err := grpc.Dial(clsAddr, grpc.WithInsecure())
		Expect(err).ToNot(HaveOccurred())
		client := idl.NewCommandListenerClient(conn)

		return client, conn
	}

	BeforeEach(func() {
		var grpcAddr string
		manager = services.NewCommandListener("foo")
		server, grpcAddr = startGRPCServer(manager)
		clsClient, connCloser = establishClient(grpcAddr)
	})

	AfterEach(func() {
		server.Stop()
		connCloser.Close()
	})

	It("connect me", func() {
		request := idl.TransmitStateRequest{"transmit request"}
		reply, err := clsClient.TransmitState(context.Background(), &request)
		Expect(err).ToNot(HaveOccurred())
		Expect(reply.GetMessage()).To(Equal("Finished echo state request: transmit request foo"))
	})
})
