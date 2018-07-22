package main

import (
	"log"
	"net"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot/platforms/raspi"

	ini "gopkg.in/ini.v1"
)

func main() {

	log.Println("Carregando arquivo de configuração...")
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
	log.Println("Iniciando a configuração do Raspberry...")
	raspiAdaptor := raspi.NewAdaptor()
	relay := gpio.NewRelayDriver(raspiAdaptor, "11")
	log.Println("Iniciando a conexão com o servidor MQTT...")
	mqttAdaptor := mqtt.NewAdaptor(serverMQTT, "luminaria-jeff")
	mqttAdaptor.SetAutoReconnect(true)
	mqttAdaptor.Publish(filaMQTT, []byte("Iniciando lumunaria em: "+getIPAddresses()))
	work := func() {
		mqttAdaptor.On(filaMQTT, func(msg mqtt.Message) {
			msgText := string(msg.Payload())
			var errRelay error
			switch msgText {
			case "1":
				log.Println("Ligando relay...")
				if relayComandoInverso == "0" {
					errRelay = relay.On()
				} else {
					errRelay = relay.Off()
				}
				if errRelay != nil {
					log.Printf("Erro ao ligar o relay: %+v\n", errRelay)
				}
			case "0":
				log.Println("Desligando relay...")
				if relayComandoInverso == "0" {
					errRelay = relay.Off()
				} else {
					errRelay = relay.On()
				}
				if errRelay != nil {
					log.Printf("Erro ao desligar o relay: %+v\n", errRelay)
				}
			}
		})
	}
	robot := gobot.NewRobot("Luminaria",
		[]gobot.Connection{raspiAdaptor, mqttAdaptor},
		[]gobot.Device{relay},
		work,
	)
	err = robot.Start()
	if err != nil {
		log.Fatalln("Erro ao iniciar a luminaria")
	}
	gobot.Every(2*time.Second, func() {
		if !robot.Running() {
			log.Println("Tentando reiniciar a luminaria...")
			err = robot.Start()
			if err != nil {
				log.Fatalln("Erro ao reiniciar a luminaria")
			}
			return
		}
		mqttAdaptor.Publish(filaMQTT, []byte("Luminaria OK em "+getIPAddresses()))
		log.Println("Luminaria ok")
	})
}

func getIPAddresses() (ips string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("getIPAddresses ", err.Error())
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips += ipnet.IP.String() + " - "
			}
		}
	}
	return
}
