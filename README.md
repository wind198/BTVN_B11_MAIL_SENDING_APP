# BTVN_B11_MAIL_SENDING_APP

## This module include following package
#### scan_publish - It scan database on a schedule and publish jobs to rabbitMQ queue when it find an order for which it need to send an email for confirmation
#### worker - It wait for jobs from rabbitMQ queue, and execute them (send email)
#### mailer - It contain the logic to handle email sending process
#### uitlity - some helper function
#### main function - It connect application to rabbitmq, declare exchange, queue and start the scan, publish, and sending email process

## Note

### You need to re-create mySQL database to run this app.
#### Open mySQL cli, create a database call "btvn_b11"
#### Run the createdb.sql file in the repo to reproduce the db
#### Set 2 ENV variable DBUSER & DBPASS as user name and password to access mySQL sever

### You will need an  API key and an authenticated email in order to use Sendgrid service, the API key need to be set as SENDGRID_API env variable, and the FromEmail constant in scan_publish package need to be updated to be your authenticated email
