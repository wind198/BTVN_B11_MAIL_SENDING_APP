package worker

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"example.com/rabbitmq/btvn_b11/scan_publish"
	u "example.com/rabbitmq/btvn_b11/utility"

	"example.com/rabbitmq/btvn_b11/mailer"
	"github.com/streadway/amqp"
)

var stopWorking = make(chan bool)

type Worker struct {
	queueName  string
	chann      *amqp.Channel
	database   *sql.DB
	sendClient *mailer.SendGridClient
}

func NewWorker(queueName string, chann *amqp.Channel, database *sql.DB) *Worker {
	sendClient := mailer.NewSendGridClient(os.Getenv("SENDGRID_API"))
	return &Worker{
		queueName:  queueName,
		chann:      chann,
		database:   database,
		sendClient: sendClient,
	}
}

func (w *Worker) Start() {
	deliveryChannle := w.register()
	w.sendEmail(deliveryChannle)

}
func (w *Worker) Stop() {
	log.Println("Stopping the worker")
	stopWorking <- true
	log.Println("Stopped the worker at", time.Now().Format("2006-01-02T15:04:05"))
}

func (w *Worker) register() <-chan amqp.Delivery {
	msgs, err := w.chann.Consume(
		w.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	u.PrintError(err, "Error during register customer")
	if err == nil {
		log.Println("Resigter suceed")
	}
	return msgs

}

func (w *Worker) sendEmail(msgs <-chan amqp.Delivery) {
	q, err := w.database.Prepare("update `order`set thankyou_email_sent = 1 where id=?")
	u.PrintError(err, "error preparing statement")
	defer q.Close()
	sendEmail := func() {
		log.Println("hello, sending email")
		for msg := range msgs {
			log.Println("found one msg")
			body := bytes.NewReader(msg.Body)
			var aEmail scan_publish.Email
			json.NewDecoder(body).Decode(&aEmail)
			err := w.sendClient.Send(aEmail)
			if err != nil {
				log.Printf("Error sending email: %v", err)
			}
			w.updateDb(q, aEmail.OrderID)
		}
	}
	go sendEmail()
	<-stopWorking
	log.Println("Stoped email sending process at", time.Now().Format("2006-01-02T15:04:05"))
}

func (w *Worker) updateDb(stm *sql.Stmt, id int) {
	res, err := stm.Exec(id)
	if err != nil {
		log.Println("Error updating DB", err)
	} else {
		rowAffected, err := res.RowsAffected()
		u.PrintError(err, "Error calculating rows affected")
		log.Printf("update succeed, %v affected", rowAffected)
	}

}
