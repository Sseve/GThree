package gtservant

import (
	"GThree/pkg/dto"
	"GThree/pkg/grpc/service"
	"GThree/pkg/utils"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type gtservantServer struct {
	service.UnimplementedZoneServer
}

type ZoneResponse struct {
	Zid    string
	Ip     string
	Name   string
	Target string
	Result string
}

// 添加区服
func addZone(gmPath string, in *service.ZoneRequest) {
	if err := os.MkdirAll(gmPath, 0755); err != nil {
		log.Println("create game path failed: ", err)
		return
	}
	cmd := fmt.Sprintf(`cd %v && svn --username %v --password %v co %v -r %v . |grep "Checked out revision"`,
		gmPath, viper.GetString("svn_username"), viper.GetString("svn_password"), viper.GetString("svn_address"),
		in.SvnVersion)
	fmt.Println(cmd)
	go runCommand(context.Background(), cmd, in)
}

// 区服管理
func manageZone(gmPath string, in *service.ZoneRequest) {
	cmd := fmt.Sprintf("cd %v && sh %v %v", gmPath, viper.GetString("zone_script"), in.Target)
	go runCommand(context.Background(), cmd, in)
}

// 更新程序
func binZone(gmPath string, in *service.ZoneRequest) {
	cmd := fmt.Sprintf(`cd %v && wget %v gameserv && chmod +x gameserv`, gmPath, viper.GetString("zone_bin_addr"))
	go runCommand(context.Background(), cmd, in)
}

// 更新配置
func conZone(gmPath string, in *service.ZoneRequest) {
	cmd := fmt.Sprintf(`cd %v && svn --username %v --password %v update -r %v | grep "At revision" `,
		gmPath, viper.GetString("svn_username"), viper.GetString("svn_password"), in.SvnVersion)
	go runCommand(context.Background(), cmd, in)
}

// 获取配置版本信息
func infoZone(gmPath string, in *service.ZoneRequest) {
	cmd := fmt.Sprintf(`cd %v && svn info|grep "Revision:"`, gmPath)
	go runCommand(context.Background(), cmd, in)
}

// 区服操作
func (g *gtservantServer) OptZone(ctx context.Context, in *service.ZoneRequest) (*service.ZoneReply, error) {
	// 业务逻辑处理
	gmPath := viper.GetString("zone_path") + in.Name + "_" + in.Zid
	switch in.Target {
	case "add":
		addZone(gmPath, in)
	case "bin":
		binZone(gmPath, in)
	case "con":
		conZone(gmPath, in)
	case "info":
		infoZone(gmPath, in)
	default:
		manageZone(gmPath, in)
	}
	return &service.ZoneReply{Zid: in.Zid, Name: in.Name, Result: "正在操作区服"}, nil
}

// 运行命令
func runCommand(ctx context.Context, command string, in *service.ZoneRequest) {
	cmd := exec.Command("/bin/bash", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		utils.Logger.Error("运行指定命令出错: ", err)
		dto.SetZResToRds(in.Name+"_"+in.Zid, err)
	}
	dto.SetZResToRds(in.Name+"_"+in.Zid, string(out))
}

// 开启rpc服务
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
		utils.Logger.Error("远程调用服务监听失败: ", err)
	}
	serve := grpc.NewServer(grpc.Creds(creds))
	service.RegisterZoneServer(serve, &gtservantServer{})
	if err := serve.Serve(listen); err != nil {
		utils.Logger.Error("远程调用服务启动失败: ", err)
		return
	}
	utils.Logger.Info("远程调用服务启动: ", viper.GetString("app_addr"))
}
