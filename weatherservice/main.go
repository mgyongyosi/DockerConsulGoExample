package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/uuid"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		fmt.Fprint(ctx, "The weather is sunny today!")
	case "/health":
		ctx.SetStatusCode(fasthttp.StatusOK)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func osSignal(err chan<- error) {
	osc := make(chan os.Signal)
	signal.Notify(osc, syscall.SIGINT, syscall.SIGTERM)
	err <- fmt.Errorf("%s", <-osc)
}

func main() {
	config := api.DefaultConfig()
	config.Address = "node03:8500"

	fmt.Println("Trying to connect to Consul server")

	client, err := api.NewClient(config)

	if err != nil {
		fmt.Printf("Couldn't connect to Consul on %s: %s", config.Address, err)
		return
	}

	host, _ := os.Hostname()

	serviceRegistration := &api.AgentServiceRegistration{
		ID:   fmt.Sprintf("weather-service-%s", uuid.NewV4()),
		Name: "weather-service",
		Checks: api.AgentServiceChecks{
			&api.AgentServiceCheck{
				HTTP:     fmt.Sprintf("http://%s:8080/health", host),
				Timeout:  "2s",
				Interval: "10s",
			},
		},
	}

	err = client.Agent().ServiceRegister(serviceRegistration)

	if err != nil {
		fmt.Printf("Couldn't register the service on %s: %s", config.Address, err)
		return
	}

	fmt.Println("Service is registered in Consul")

	errors := make(chan error)

	go osSignal(errors)

	go func() {
		errors <- fasthttp.ListenAndServe(":8080", fastHTTPHandler)
	}()

	<-errors

	fmt.Println("Stopping service")

	client.Agent().ServiceDeregister(serviceRegistration.ID)

	fmt.Println("Service is deregistered from Consul")

}
