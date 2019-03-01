package sdk

import (
	"context"
	"fmt"
	h3c "github.com/h3copen/comwaresdk/cmwproto/grpc_service"
	"google.golang.org/grpc"
	md "google.golang.org/grpc/metadata"
	"log"
	"time"
)

// type GrpcClient interface {
// 	Close()
// 	AddUnicastRoute(m map[string]MapValueIntface)
// 	AddUnicastRouteV6(m map[string]MapValueIntface)
// 	AddUnicastRoutes(ms []map[string]MapValueIntface)
// 	AddUnicastRoutesV6(ms []map[string]MapValueIntface)
// 	DelUnicastRoute(delRoute map[string]string)
// 	DelUnicastRoutes(delRoutes []map[string]string)
// 	DelUnicastRouteV6(delRoute map[string]string)
// 	DelUnicastRoutesV6(delRoutes []map[string]string)
// 	GetDeviceInfo() string
// }

type GrpcSession struct {
	Client h3c.GrpcServiceClient
	Conn   *grpc.ClientConn
	Token  string
}

func NewClient(addr string, port uint, username string, password string) (*GrpcSession, error) {
	address := fmt.Sprintf("%s:%d", addr, port)

	log.Printf("Server address: %v, UserName: %v, Password: %v\n", address, username, password)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}

	//create grpc_service client
	c := h3c.NewGrpcServiceClient(conn)

	//prepare context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var token string

	loginReply, err := c.Login(ctx, &h3c.LoginRequest{UserName: &username, Password: &password})
	if err != nil {
		log.Fatalf("could not login: %v", err)
		conn.Close()
		return nil, err
	}

	token = loginReply.GetTokenId()
	log.Printf("Token: %s", token)

	s := &GrpcSession{Client: c,
		Conn:  conn,
		Token: token}

	return s, nil
}

func (s *GrpcSession) Close() {
	var logoutReq = h3c.LogoutRequest{TokenId: &s.Token}
	ctx, cancel := CtxWithToken(s.Token, time.Second)
	defer cancel()
	s.Client.Logout(ctx, &logoutReq)
	s.Conn.Close()
	return
}

func CtxWithToken(tk string, timeout time.Duration) (context.Context, context.CancelFunc) {
	//Add token to meta data
	var mdata = md.Pairs("token_id", tk)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	var ctx_with_token = md.NewOutgoingContext(ctx, mdata)
	return ctx_with_token, cancel
}
