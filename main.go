package main

import (
	"log"
	"os"
	"os/signal"

	"example.com/rabbitmq/btvn_b11/scan_publish"
	"example.com/rabbitmq/btvn_b11/todb"
	u "example.com/rabbitmq/btvn_b11/utility"
	"example.com/rabbitmq/btvn_b11/worker"
	"github.com/streadway/amqp"
)

func main() {
	stopWorking := make(chan bool)
	database := todb.ConnectDB()
	//declare an exchange
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	u.PrintError(err, "failed to connect rabbitmq")
	chann, err := conn.Channel()
	exchangeName := "mail exchange"
	u.PrintError(err, "error during generate channel")

	err = chann.ExchangeDeclare(
		exchangeName, "direct", false, false, false, false, nil,
	)
	u.PrintError(err, "error during exchange declaration")
	//declare the queue
	queue, err := chann.QueueDeclare(
		"msg queue", false, false, false, false, nil)
	u.PrintError(err, "error during declare queue")

	//bind queue to exchange
	err = chann.QueueBind(queue.Name, "", exchangeName, false, nil)
	u.PrintError(err, "error during binding queue")

	//declare the publisher and bind it to the exchange
	publisher := scan_publish.NewPublisher(
		chann,
		exchangeName,
		database,
	)

	//declare the worker and bind them to the queu
	worker := worker.NewWorker(
		queue.Name, chann, database,
	)

	//let the main routine run until interuption signal
	go publisher.Start()
	go worker.Start()
	signalChann := make(chan os.Signal)
	signal.Notify(signalChann, os.Interrupt)
	listeningForSignal := func() {
		sig := <-signalChann
		log.Printf("Got %s signal. Exitting...", sig)
		publisher.Stop()
		// worker.Stop()
		stopWorking <- true
	}
	go listeningForSignal()
	<-stopWorking
	//End main
}
