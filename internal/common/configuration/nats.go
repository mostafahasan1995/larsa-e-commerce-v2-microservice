package configuration

import (
	"errors"
	"fmt"
	"log"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/interval"
	"git.larsa.io/mahdawi/microservices-commons.git/nats"
	"github.com/google/uuid"
)

func InitNats() {
	env := EnvironmentVariables.NatsJetStream
	if env == "false" {
		InitNatsNor()
	} else {
		InitNatsJet()
	}

}

func InitNatsJet() {
	url := EnvironmentVariables.NatsUrl
	client := nats.ConnectJet("microservices", uuid.NewString(), url)
	if client == nil {
		err := interval.Do(2, func(attempt int) (bool, error) {
			client = nats.ConnectJet("microservices", uuid.NewString(), url)
			fmt.Println("trying to connect nats server", time.Now())
			if client == nil {
				return true, errors.New("could not connect to nats server")
			}
			return true, nil
		})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("connected to nats server")
		}
	} else {
		fmt.Println("connected to nats server")

	}
}
func InitNatsNor() {
	url := EnvironmentVariables.NatsUrl
	client := nats.Connect("microservices", uuid.NewString(), url)
	if client == nil {
		err := interval.Do(2, func(attempt int) (bool, error) {
			client = nats.Connect("microservices", uuid.NewString(), url)
			fmt.Println("trying to connect nats server", time.Now())
			if client == nil {
				return true, errors.New("could not connect to nats server")
			}
			return true, nil
		})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("connected to nats server")
		}
	} else {
		fmt.Println("connected to nats server")

	}
}
