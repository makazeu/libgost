package libgost

import (
	"crypto/tls"
	"errors"
	"os"

	"github.com/ginuerzh/gost"
	"github.com/go-log/log"
)

type StringList []string

type GostConfig baseConfig

func init() {
	gost.SetLogger(&gost.LogLogger{})
}

func NewGost(chainNodes, serveNodes StringList) *GostConfig {
	baseCfg := &GostConfig{}
	baseCfg.ChainNodes = stringList(chainNodes)
	baseCfg.ServeNodes = stringList(serveNodes)
	return baseCfg
}

func StartServing(config *GostConfig) {
	// generate random self-signed certificate.
	cert, err := gost.GenCertificate()
	if err != nil {
		log.Log(err)
		os.Exit(1)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	gost.DefaultTLSConfig = tlsConfig

	if err := start(config); err != nil {
		log.Log(err)
		os.Exit(1)
	}

	select {}
}

func start(config *GostConfig) error {
	gost.Debug = config.Debug

	var routers []router
	rts, err := config.route.GenRouters()
	if err != nil {
		return err
	}
	routers = append(routers, rts...)

	if len(routers) == 0 {
		return errors.New("invalid config")
	}
	for i := range routers {
		go routers[i].Serve()
	}

	return nil
}
