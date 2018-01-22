package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jhunt/vcaptive"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

func cfg(deflt, env string) string {
	if s := os.Getenv(env); s != "" {
		return s
	}
	return deflt
}

func main() {
	broker := &Broker{}
	broker.Service.ID = cfg("mariadb-43e22be8-5a3a-496c-b502-02079483f6dd", "SERVICE_ID")
	broker.Service.Name = cfg("mariadb", "SERVICE_NAME")
	broker.Plan.ID = cfg("mariadb-43e22be8-5a3a-496c-b502-02079483f6dd", "PLAN_ID")
	broker.Plan.Name = cfg("shared", "PLAN_NAME")
	broker.Description = cfg("A shared MariaDB database", "DESCRIPTION")
	broker.Tags = strings.Split(cfg("shared,mariadb,tinsmith", "TAGS"), ",")

	app, err := vcaptive.ParseApplication(os.Getenv("VCAP_APPLICATION"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "VCAP_APPLICATION: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("running v%s of %s at http://%s\n", app.Version, app.Name, app.URIs[0])
	services, err := vcaptive.ParseServices(os.Getenv("VCAP_SERVICES"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "VCAP_SERVICES: %s\n", err)
		os.Exit(1)
	}

	var (
		found    bool
		instance vcaptive.Instance
	)
	if name := os.Getenv("USE_SERVICE"); name != "" {
		instance, found = services.Named(name)
		if !found {
			fmt.Fprintf(os.Stderr, "VCAP_SERVICES: no service named '%s' found\n", name)
			os.Exit(2)
		}

	} else {
		instance, found = services.Tagged("mariadb")
		if !found {
			fmt.Fprintf(os.Stderr, "VCAP_SERVICES: no 'mariadb' service found\n")
			os.Exit(2)
		}
	}

	if s, ok := instance.GetString("username"); ok {
		broker.Username = s
	} else {
		fmt.Fprintf(os.Stderr, "VCAP_SERVICES: '%s' service has no 'username' credential\n", instance.Label)
		os.Exit(3)
	}
	if s, ok := instance.GetString("password"); ok {
		broker.Password = s
	} else {
		fmt.Fprintf(os.Stderr, "VCAP_SERVICES: '%s' service has no 'password' credential\n", instance.Label)
		os.Exit(3)
	}
	if s, ok := instance.GetString("host"); ok {
		broker.Host = s
	} else {
		fmt.Fprintf(os.Stderr, "VCAP_SERVICES: '%s' service has no 'host' credential\n", instance.Label)
		os.Exit(3)
	}
	if u, ok := instance.GetUint("port"); ok {
		broker.Port = fmt.Sprintf("%d", u)
	} else {
		fmt.Fprintf(os.Stderr, "VCAP_SERVICES: '%s' service has no 'port' credential; using default of 3306\n", instance.Label)
		broker.Port = "3306"
	}

	if err := broker.Init(); err != nil {
		panic(err)
	}

	http.Handle("/", brokerapi.New(
		broker,
		lager.NewLogger("mariadb-tinsmith"),
		brokerapi.BrokerCredentials{
			Username: cfg("b-mariadb", "SB_BROKER_USERNAME"),
			Password: cfg("mariadb", "SB_BROKER_PASSWORD"),
		},
	))
	err = http.ListenAndServe(":"+cfg("3000", "PORT"), nil)
	fmt.Fprintf(os.Stderr, "http server exited: %s\n", err)
}
