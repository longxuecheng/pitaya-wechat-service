package main

// Now time this server is embedded in this project, later on may be extracted as a independent service

import (
	context "context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"pitaya-wechat-service/rpc"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/muesli/cache2go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 50052, "The server port")
)

type authServiceServer struct {
	tokenCache *cache2go.CacheTable
}

var hmacSampleSecret = []byte("my_secret_key")

func (s *authServiceServer) Authorize(ctx context.Context, request *rpc.AuthRequest) (*rpc.TokenResponse, error) {
	log.Printf("authorization request is %v", request)
	// following token
	jwtMap := jwt.MapClaims{
		"name": "lxc",
		"id":   1,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtMap)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return nil, err
	}
	tokenResponse := &rpc.TokenResponse{
		Token: tokenString,
		Ttl:   10000,
	}
	s.tokenCache.Add(tokenString, 5*time.Second, jwtMap)
	return tokenResponse, nil
}

func newServer() *authServiceServer {
	server := new(authServiceServer)
	server.tokenCache = cache2go.Cache("token_cache")
	return server
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterAuthorizationServiceServer(grpcServer, newServer())
	log.Println("rpc server is listening on port ", *port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("rpc server start failed")
	}
}
