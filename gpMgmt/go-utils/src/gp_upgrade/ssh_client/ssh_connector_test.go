package ssh_client_test

import (
	"golang.org/x/crypto/ssh"

	"io"

	"bytes"

	"errors"

	"path"
	"runtime"

	"io/ioutil"

	"gp_upgrade/ssh_client"

	"bufio"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	param_network string
	param_addr    string
	param_config  *ssh.ClientConfig
)

var _ = Describe("SshConnector", func() {
	var (
		subject       *ssh_client.RealSshConnector
		test_key_path string
	)
	BeforeEach(func() {
		_, this_file_path, _, _ := runtime.Caller(0)
		test_key_path = path.Join(path.Dir(this_file_path), "../integrations/sshd/private_key.pem")
		subject = &ssh_client.RealSshConnector{
			SshDialer:      FakeDialer{},
			SshKeyParser:   FakeKeyParser{},
			PrivateKeyPath: test_key_path,
		}

	})

	Describe("#New", func() {
		It("populates the private key correctly", func() {
			const PRIVATE_KEY_FILE_PATH = "/tmp/testPrivateKeyFile.key"
			ioutil.WriteFile(PRIVATE_KEY_FILE_PATH, []byte("----TEST PRIVATE KEY ---"), 0600)
			sshConnector, err := ssh_client.NewSshConnector(PRIVATE_KEY_FILE_PATH)
			Expect(err).ToNot(HaveOccurred())
			Expect(sshConnector.(*ssh_client.RealSshConnector).PrivateKeyPath).To(Equal(PRIVATE_KEY_FILE_PATH))
		})
		It("returns an error when private key is missing", func() {
			_, err := ssh_client.NewSshConnector("pathThatDoesNotExist")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("#Connect", func() {
		Context("Happy path connection", func() {
			It("returns a session", func() {
				_, err := subject.Connect("localhost", 22, "gpadmin")

				Expect(err).ToNot(HaveOccurred())
			})
			It("calls dial with correct parameters", func() {
				_, err := subject.Connect("localhost", 22, "gpadmin")

				Expect(err).ToNot(HaveOccurred())
				Expect(param_network).To(Equal("tcp"))
				Expect(param_addr).To(Equal("localhost:22"))
				Expect(param_config.User).To(Equal("gpadmin"))
				// docker container has ssh client library that requires a callback
				Expect(param_config.HostKeyCallback).ToNot(Equal(nil))
				Expect(len(param_config.Auth)).To(Equal(1))
			})
		})

		Context("errors", func() {
			Context("private key file cannot be opened", func() {
				It("returns an error message", func() {
					subject.PrivateKeyPath = "invalid_private_key"
					_, err := subject.Connect("", 0, "")

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("open invalid_private_key: no such file or directory"))
				})
			})

			Context("private key file cannot be parsed", func() {
				It("returns an error message", func() {
					subject.SshKeyParser = ThrowingKeyParser{}

					_, err := subject.Connect("", 0, "")

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("test parsing failure"))
				})
			})

			Context("dialing connection returns error", func() {
				It("returns an error message", func() {
					subject.SshDialer = ThrowingDialer{}

					_, err := subject.Connect("", 0, "")

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("test dialing failure"))
				})
			})
			Context("new session returns error", func() {
				It("returns an error message", func() {
					subject.SshDialer = ThrowingBadClientDialer{}

					_, err := subject.Connect("", 0, "")

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("test NewSession failure"))
				})
			})
		})
	})

	Describe("#ConnectAndExecute", func() {
		Context("happy: when command runs successfully", func() {
			XIt("it returns the output from a command", func() {
				result, err := subject.ConnectAndExecute("localhost", 22, "gpadmin", "foo")

				// todo have to control the ssh session that results...  see SshClient.NewSession... should mock out but fails somehow
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal("foobar"))
			})
		})
	})

})

type FakeSigner struct{}

func (signer FakeSigner) PublicKey() ssh.PublicKey {
	return nil
}
func (signer FakeSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return nil, nil
}

type FakeKeyParser struct{}

func (parser FakeKeyParser) ParsePrivateKey(pemBytes []byte) (ssh.Signer, error) {
	return FakeSigner{}, nil
}

type FakeSshClient struct {
	//sshSession ssh.Session
	buf bytes.Buffer
}

//func (fakeSshClient FakeSshClient) NewSession() (ssh_client.Session, error) {
func (fakeSshClient FakeSshClient) NewSession() (*ssh.Session, error) {
	fakeSshClient.buf.Write([]byte("test session"))
	reader := bufio.NewReader(&fakeSshClient.buf)
	result := &ssh.Session{Stdin: reader}
	return result, nil

	//fakeSession := FakeSession{}
	//return &fakeSession, nil
}

//func (fakeSshClient FakeSshClient) NewSession() (Session, error) {
//	fakeSshClient.buf.Write([]byte("test session"))
//	//reader := bufio.NewReader(&fakeSshClient.buf)
//	//result := &ssh.Session{Stdin: reader}
//	fakeSession := &FakeSession{bufio.NewReader(&fakeSshClient.buf)}
//	return fakeSession, nil
//}

type FakeSession struct{}

func (fakeSession FakeSession) Close() error {
	return nil
}

func (fakeSession FakeSession) Output(string) ([]byte, error) {
	return []byte("fake session output"), nil
}

type FakeDialer struct{}

//func (dialer FakeDialer) Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
func (dialer FakeDialer) Dial(network, addr string, config *ssh.ClientConfig) (ssh_client.SshClient, error) {
	param_network = network
	param_addr = addr
	param_config = config
	return &FakeSshClient{}, nil
}

type ThrowingKeyParser struct{}

func (parser ThrowingKeyParser) ParsePrivateKey(pemBytes []byte) (ssh.Signer, error) {
	return nil, errors.New("test parsing failure")
}

type ThrowingDialer struct{}

func (dialer ThrowingDialer) Dial(network, addr string, config *ssh.ClientConfig) (ssh_client.SshClient, error) {
	return nil, errors.New("test dialing failure")
}

type ThrowingClient struct{}

//func (fakeSshClient ThrowingClient) NewSession() (ssh_client.Session, error) {
func (fakeSshClient ThrowingClient) NewSession() (*ssh.Session, error) {
	return nil, errors.New("test NewSession failure")
}

type ThrowingBadClientDialer struct{}

//func (badClientDialer ThrowingBadClientDialer) Dial(network, addr string, config *ssh.ClientConfig) (ssh_client.SshClient, error) {
func (badClientDialer ThrowingBadClientDialer) Dial(network, addr string, config *ssh.ClientConfig) (ssh_client.SshClient, error) {
	return new(ThrowingClient), nil
}
