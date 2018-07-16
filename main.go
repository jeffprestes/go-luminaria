package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot/platforms/raspi"

	ini "gopkg.in/ini.v1"
)

func main() {

	fmt.Println("Carregando arquivo de configuração...")
	cfg, err := ini.Load("config.ini")
	if err != nil {
		panic("Erro ao carregar arquivo de configuração: " + err.Error())
	}
	filaMQTT := cfg.Section("").Key("fila").Value()
	serverMQTT := cfg.Section("").Key("servidor_url").Value()
	relayComandoInverso := cfg.Section("").Key("relay_inverso").Value()
	if len(filaMQTT) < 10 || len(serverMQTT) < 10 {
		panic("Erro ao carregar arquivo de configuração: não foi possivel carregar os valores do server ou da fila MQTT")
	}
	fmt.Println("Iniciando a configuração do Raspberry...")
	raspiAdaptor := raspi.NewAdaptor()
	relay := gpio.NewRelayDriver(raspiAdaptor, "11")
	fmt.Println("Iniciando a conexão com o servidor MQTT...")
	mqttAdaptor := mqtt.NewAdaptor(serverMQTT, "luminaria-jeff")
	work := func() {
		mqttAdaptor.On(filaMQTT, func(msg mqtt.Message) {
			msgText := string(msg.Payload())
			switch msgText {
			case "1":
				fmt.Println("Ligando relay...")
				if relayComandoInverso == "0" {
					err = relay.On()
				} else {
					err = relay.Off()
				}
				if err != nil {
					fmt.Printf("Erro ao ligar o relay: %+v\n", err)
				}
			case "0":
				fmt.Println("Desligando relay...")
				if relayComandoInverso == "0" {
					err = relay.Off()
				} else {
					err = relay.On()
				}
				if err != nil {
					fmt.Printf("Erro ao desligar o relay: %+v\n", err)
				}
			}
		})
	}
	robot := gobot.NewRobot("Luminaria",
		[]gobot.Connection{raspiAdaptor, mqttAdaptor},
		[]gobot.Device{relay},
		work,
	)
	robot.Start()
}
