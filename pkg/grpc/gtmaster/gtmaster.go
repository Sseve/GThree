package gtmaster

import (
	"GThree/pkg/grpc/service"
	"GThree/pkg/models"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var ZoneResult = make(chan ZoneResponse, 2)

type ZoneResponse struct {
	Zid    uint32
	Ip     string
	Name   string
	Target string
	Result string
}

// 该函数的参数为区服信息
func CallServant(zone models.Zone) {
	// 加入认证
	cert, _ := tls.LoadX509KeyPair(viper.GetString("app_pem_file"), viper.GetString("app_key_file"))
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile(viper.GetString("app_ca_pem"))
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   viper.GetString("app_serv_name"),
		RootCAs:      certPool,
	})

	address := zone.Ip + viper.GetString("app_rpc_port")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Println("Call gtservant failed: ", err)
	}
	defer conn.Close()

	c := service.NewZoneClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 请求grpc
	resp, err := c.OptZone(ctx, &service.ZoneRequest{Ip: zone.Ip, Zid: zone.ZId, Name: zone.Name, Target: zone.Targt})
	if err != nil {
		log.Println("Request gtservant failed: ", err)
	}
	r := ZoneResponse{
		Zid:    zone.ZId,
		Ip:     zone.Ip,
		Name:   zone.Name,
		Target: zone.Targt,
		Result: resp.GetResult(),
	}
	// 获取响应
	ZoneResult <- r
}
