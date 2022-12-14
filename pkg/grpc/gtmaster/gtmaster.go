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

type ZoneResponse struct {
	Zid    string
	Ip     string
	Name   string
	Target string
	Result string
}

// rpc区服管理
func ZoneServant(zone models.Zone, ZoneResult chan ZoneResponse) {
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
	ctx1, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*3))
	defer cancel()
	conn, err := grpc.DialContext(ctx1, address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Println("Call gtservant failed: ", err)
	}
	defer conn.Close()

	c := service.NewZoneClient(conn)

	ctx2, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 请求grpc
	resp, err := c.OptZone(ctx2, &service.ZoneRequest{Ip: zone.Ip, Zid: zone.Zid, Name: zone.Name, Target: zone.Targt, SvnVersion: zone.SvnVersion})
	if err != nil {
		log.Println("Request gtservant failed: ", err)
	}
	r := ZoneResponse{
		Zid:    zone.Zid,
		Ip:     zone.Ip,
		Name:   zone.Name,
		Target: zone.Targt,
		Result: resp.GetResult(),
	}
	// 获取响应
	ZoneResult <- r
}
