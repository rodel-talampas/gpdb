package commanders_test

import (
	"errors"
	"gp_upgrade/cli/commanders"
	mockpb "gp_upgrade/mock_idl"
	"testing"

	pb "gp_upgrade/idl"

	"github.com/golang/mock/gomock"
	"github.com/greenplum-db/gpbackup/testutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("reporter", func() {

	var (
		client *mockpb.MockCliToHubClient
		t      *testing.T
		ctrl   *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(t)
		client = mockpb.NewMockCliToHubClient(ctrl)
	})

	AfterEach(func() {
		defer ctrl.Finish()
	})

	Describe("ConvertMaster", func() {
		It("Reports success when pg_upgrade started", func() {
			_, testStdout, _, _ := testutils.SetupTestLogger()
			client.EXPECT().UpgradeConvertMaster(
				gomock.Any(),
				&pb.UpgradeConvertMasterRequest{},
			).Return(&pb.UpgradeConvertMasterReply{}, nil)
			err := commanders.NewUpgrader(client).ConvertMaster("", "", "", "")
			Expect(err).To(BeNil())
			Eventually(testStdout).Should(gbytes.Say("Kicked off pg_upgrade request"))
		})

		It("reports failure when command fails to connect to the hub", func() {
			_, _, testStderr, _ := testutils.SetupTestLogger()
			client.EXPECT().UpgradeConvertMaster(
				gomock.Any(),
				&pb.UpgradeConvertMasterRequest{},
			).Return(&pb.UpgradeConvertMasterReply{}, errors.New("something bad happened"))
			err := commanders.NewUpgrader(client).ConvertMaster("", "", "", "")
			Expect(err).ToNot(BeNil())
			Eventually(testStderr).Should(gbytes.Say("ERROR - Unable to connect to hub"))

		})
	})

})
