package microservice

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/hashicorp/consul/api"

	consulsd "github.com/go-kit/kit/sd/consul"
)

// registration
func Register(consulAddress string,
	consulPort string,
	advertiseAddress string,
	advertisePort string) (registrar sd.Registrar) {

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		consulConfig.Address = consulAddress + ":" + consulPort
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	check := api.AgentServiceCheck{
		HTTP:     "http://" + advertiseAddress + ":" + advertisePort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Basic health checks",
	}

	port, _ := strconv.Atoi(advertisePort)
	num := rand.Intn(100) // to make service ID unique
	asr := api.AgentServiceRegistration{
		ID:      "microservice" + strconv.Itoa(num), //unique service ID
		Name:    "microservice",                     // 高亮行
		Address: advertiseAddress,
		Port:    port,
		Tags:    []string{"microservice"}, // 高亮行
		Check:   &check,
	}
	registrar = consulsd.NewRegistrar(client, &asr, logger)
	return
}
