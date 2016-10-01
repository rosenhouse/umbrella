package accumulate

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
)

type AcquirerRPC struct {
	dispensary *dispensary
}

func (a *AcquirerRPC) AcquireOne(args struct{}, reply *string) error {
	path, err := a.dispensary.AcquireOne()
	if err != nil {
		return err
	}
	*reply = path
	return nil
}

type Server interface {
	Address() string
	GoListenAndServe() error
	Close() error
	ListAll() []string
}

type server struct {
	dispensary *dispensary
	address    string
	listener   net.Listener
}

func NewServer() (Server, error) {
	serverDir, err := ioutil.TempDir("", "umbrella.accumulate.dispensary")
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(serverDir, "data")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}

	serverAddress := filepath.Join(serverDir, "socket")

	srv := &server{
		dispensary: &dispensary{
			dir: dataDir,
		},
		address: serverAddress,
	}

	return srv, nil
}

func (s *server) Address() string   { return s.address }
func (s *server) ListAll() []string { return s.dispensary.ListAll() }

func (s *server) GoListenAndServe() error {
	objectToServe := &AcquirerRPC{s.dispensary}

	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("Dispensary", objectToServe)

	var err error
	s.listener, err = net.Listen("unix", s.address)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				return // likely just a cancelation
			}

			go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()

	cc, err := net.Dial("unix", s.address)
	if err != nil {
		return fmt.Errorf("validating new server: %s", err)
	}
	return cc.Close()
}

func (s *server) Close() error {
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}
	return s.dispensary.Cleanup()
}

func AcquireRemote(testServerAddr string) (string, error) {
	conn, err := net.Dial("unix", testServerAddr)
	if err != nil {
		return "", fmt.Errorf("dial umbrella server: %s", err)
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	var path string
	err = client.Call("Dispensary.AcquireOne", struct{}{}, &path)
	if err != nil {
		return "", fmt.Errorf("rpc AcquireOne: %s", err)
	}

	return path, nil
}
