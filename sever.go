package main

import (
	"fmt"
	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var (
	cfgFile string
)

type ServerPort struct {
	Addr         string
	Password     string
	Ciph         core.Cipher
	InFlowMeter  int64
	OutFlowMeter int64
}

var ServerControl Control

type Control struct {
	servers map[string]*ServerPort
}

func init() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("'config.yaml' file read error:", err)
		os.Exit(0)
	}
	info := viper.GetStringMap("server")
	ServerControl.servers = make(map[string]*ServerPort, 10)
	for _, v := range info {
		var serverport ServerPort
		portpass := v.(string)
		out := strings.Split(portpass, "/")
		serverport.Addr = ":" + out[0]
		serverport.Password = out[1]
		ServerControl.servers[out[0]] = &serverport
	}
}
func NewServerPort(cipher, password string, key []byte) core.Cipher {
	ciph, err := core.PickCipher(cipher, key, password)
	if err != nil {
		log.Fatal(err)
	}
	return ciph
}

func StarUpServers() {
	for _, v := range ServerControl.servers {
		ciph := NewServerPort("aes-128-gcm", v.Password, []byte{})
		v.Ciph = ciph
		log.Printf("star port %s server\n", v.Addr)
		go udpRemote(v.Addr, ciph.PacketConn)
		go tcpRemote(v.Addr, ciph.StreamConn)
	}

}
