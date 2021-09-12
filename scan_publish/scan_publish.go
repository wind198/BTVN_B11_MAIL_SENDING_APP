package scan_publish

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	u "example.com/rabbitmq/btvn_b11/utility"

	"github.com/robfig/cron/v3"
	"github.com/streadway/amqp"
)

var workChann = make(chan *Customer, 10) //All customer to sent email will be stored in this channel

var stopPublishing = make(chan bool) //This channel used to stop email publishing process

var stop = make(chan bool) //This channel used to stop the whole scanning and publish process

//these constant are used to build email from order information
const (
	EmailFrom                = "tuanbk1908@gmail.com"
	DefaultThankyouSubject   = "Thank you for purchasing from mystore.com"
	DefaultThankyouBodyPlain = "Thank you for purchasing from our store. Here's your order details:"
	DefaultThankyouBodyHtml  = "<strong>Thank you for purchasing from our store. Here's your order details:</strong>"
	DefaultFromName          = "My Store Owner"
	DefaultFromEmail         = "support@mystore.com"
)

//Publisher will scan the database for customer to send email

type Publisher struct {
	scheduler *cron.Cron    //this schedule and run jobs of scanning database
	chann     *amqp.Channel //channel to rabbitmq
	exchange  string        //exchange to publish message to
	database  *sql.DB       //database to scan

}

//Represent a customer information to send email
type Customer struct {
	OrderID  int     `json:"orderID"`
	Name     string  `json:"customerName"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	CreateAt string  `json:"createAt"`
	Email    string  `json:"customerEmail"`
}

//Represent an email that need to send
type Email struct {
	OrderID     int    `json:"orderID"`
	From        string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	HTMLcontent string `json:"htmlContent"`
}

//Create an email object from customer informatioin
func (c *Customer) email() *Email {
	var email = &Email{}
	log.Println(c.Email)
	email.To = c.Email
	email.From = EmailFrom
	email.Subject = fmt.Sprint("Order confirmation email from mystore.com")
	email.Content = fmt.Sprintf(DefaultThankyouBodyPlain+`
	Total bill: %v,
	Created at: %v,`, c.Price, c.CreateAt)
	email.HTMLcontent = fmt.Sprintf(DefaultThankyouBodyHtml+
		`
	<ul>
	<li>Total bill: <strong>%v</strong></li>
	<li>Created at: <strong>%v</strong></li>
	</ul>`, c.Price, c.CreateAt)
	email.OrderID = c.OrderID
	return email
}

//Create a new publisher
func NewPublisher(chann *amqp.Channel, exchange string, database *sql.DB) *Publisher {
	scheduler := cron.New(cron.WithSeconds())
	return &Publisher{
		scheduler: scheduler,
		chann:     chann,
		exchange:  exchange,
		database:  database,
	}
}

//start the publisher, start the scanning process and start publishing email to rabbit
func (p *Publisher) Start() {
	p.iterateScanDb()
	go p.publish()
	<-stop
}

//stop the publisher, stop scanning process and email publishing process
func (p *Publisher) Stop() {
	log.Printf("Publisher closing")
	p.stopIterateScanDb()
	p.stopPublishing()
	stop <- true

}

//Schedule database scanning frequency and start the scanning process
func (p *Publisher) iterateScanDb() {

	p.scheduler.AddFunc("*/5 * * * * *", func() { scanDb(p.database) })
	p.scheduler.Start()
}

//stop scanning process
func (p *Publisher) stopIterateScanDb() {
	p.scheduler.Stop()
	close(workChann)
	log.Println("stopeiterating scanning db at", time.Now().Format("2006-01-02T15:04:05"))
}

//stop email publishing process
func (p *Publisher) stopPublishing() {
	stopPublishing <- true
	log.Println("stoped puishing message to rabbitmq at", time.Now().Format("2006-01-02T15:04:05"))
}

//scan db for customer that need to send email
func scanDb(database *sql.DB) {
	q := "SELECT id,customer_name,total_price,currency,created_at,email from `order` where thankyou_email_sent=0 and cancelled_at is null"
	rows, err := database.Query(q)
	u.PrintError(err, "Error during query db")
	for rows.Next() {
		var aCustomer Customer
		err := rows.Scan(&aCustomer.OrderID, &aCustomer.Name, &aCustomer.Price, &aCustomer.Currency, &aCustomer.CreateAt, &aCustomer.Email)
		u.PrintError(err, "Error during scanning")
		if err == nil {
			workChann <- &aCustomer
			log.Println("sending 1 customer to chanel")
		}
	}
	rows.Close()
	u.PrintError(err, "Error during query and iteration")
}

func (p *Publisher) publish() {
	publishEmailToExchange := func() {

		for customer := range workChann {
			aEmail := customer.email()
			body, err := json.Marshal(aEmail)
			u.PrintError(err, "Error during maring")
			if err != nil {
				return
			}
			p.chann.Publish(
				p.exchange, "", false, false, amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				},
			)
			log.Println("senidng 1 email to exchange")

		}
	}
	go publishEmailToExchange()
	<-stopPublishing
}
