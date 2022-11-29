package gtservant

import (
	"GThree/pkg/grpc/service"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type gtservantServer struct {
	service.UnimplementedZoneServer
}

func (g *gtservantServer) OptZone(ctx context.Context, in *service.ZoneRequest) (*service.ZoneReply, error) {
	// 业务逻辑处理
	switch in.Target {
	case "add":
		// 开服
		fmt.Println(in.Ip, in.Name, in.Zid, in.Target)
	case "start":
		// 启动
	case "stop":
		// 关闭
	case "check":
		// 检查
	default:
		// 未知操作
	}
	return &service.ZoneReply{Zid: in.Zid, Name: in.Name, Result: fmt.Sprintf("xx== %d ==xx", in.Zid)}, nil
}

func Start() {

	// 加入证书
	cert, _ := tls.LoadX509KeyPair(viper.GetString("app_pem_file"), viper.GetString("app_key_file"))
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile(viper.GetString("app_ca_pem"))
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
		ClientCAs:    certPool,
	})

	listen, err := net.Listen("tcp", viper.GetString("app_addr"))
	if err != nil {
		log.Println("start listen failed: ", err)
	}

	serve := grpc.NewServer(grpc.Creds(creds))

	log.Println("server start on: ", viper.GetString("app_addr"))
	service.RegisterZoneServer(serve, &gtservantServer{})
	if err := serve.Serve(listen); err != nil {
		log.Println("server start failed: ", err)
		return
	}
}
