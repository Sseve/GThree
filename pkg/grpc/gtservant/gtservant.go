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
	"os"
	"os/exec"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type gtservantServer struct {
	service.UnimplementedZoneServer
}

func addZone(gmPath string, in *service.ZoneRequest) (string, error) {
	if err := os.MkdirAll(gmPath, os.ModePerm); err != nil {
		log.Println("create game path failed: ", err)
		return "", err
	}
	cmd := fmt.Sprintf(`cd %v && svn --username %v --password %v co %v -r %v . |grep "Checked out revision"`,
		gmPath, viper.GetString("svn_username"), viper.GetString("svn_password"), viper.GetString("svn_address"),
		in.SvnVersion)
	return runCommand(cmd)
}

func manageZone(gmPath string, in *service.ZoneRequest) (string, error) {
	cmd := fmt.Sprintf("cd %v && sh %v %v", gmPath, viper.GetString("zone_script"), in.Target)
	return runCommand(cmd)
}

func binZone(gmPath string, in *service.ZoneRequest) (string, error) {
	cmd := fmt.Sprintf(`cd %v && wget %v gameserv && chmod +x gameserv`, gmPath, viper.GetString("zone_bin_addr"))
	return runCommand(cmd)
}

func conZone(gmPath string, in *service.ZoneRequest) (string, error) {
	cmd := fmt.Sprintf(`cd %v && svn --username %v --password %v update -r %v | grep "Updated to revision" `,
		gmPath, viper.GetString("svn_username"), viper.GetString("svn_password"), in.SvnVersion)
	return runCommand(cmd)
}

func (g *gtservantServer) OptZone(ctx context.Context, in *service.ZoneRequest) (*service.ZoneReply, error) {
	// 业务逻辑处理
	var (
		cmdOut string
		err    error
		gmPath string = viper.GetString("zone_path") + in.Name + "_" + in.Zid
	)
	if in.Target == "add" {
		// 开服
		cmdOut, err = addZone(gmPath, in)
		// cmdOut = "开服成功"
	} else if in.Target == "bin" {
		// 更新bin程序文件
		cmdOut, err = binZone(gmPath, in)
		// cmdOut = "更新bin文件成功"
	} else if in.Target == "con" {
		// 更新配置文件
		cmdOut, err = conZone(gmPath, in)
		// cmdOut = "更新配置文件成功"
	} else {
		cmdOut, err = manageZone(gmPath, in)
		// cmdOut = fmt.Sprintf("%v %v %v 成功", in.Name, in.Zid, in.Target)
	}

	if err != nil {
		return &service.ZoneReply{Zid: in.Zid, Name: in.Name, Result: fmt.Sprintf("%v", err)}, nil
	} else {
		return &service.ZoneReply{Zid: in.Zid, Name: in.Name, Result: fmt.Sprintf("%v", cmdOut)}, nil
	}

}

// 运行命令
func runCommand(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(out), err)
		return "", err
	}
	return string(out), nil
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
