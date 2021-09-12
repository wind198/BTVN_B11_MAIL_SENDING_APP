# BTVN_B11_MAIL_SENDING_APP

## This module include following package
#### scan_publish - It scan database on a schedule and publish jobs to rabbitMQ queue when it find an order for which it need to send an email for confirmation
#### worker - It wait for jobs from rabbitMQ queue, and execute them (send email)
#### mailer - It contain the logic to handle email sending process
#### uitlity - some helper function
#### main function - It connect application to rabbitmq, declare exchange, queue and call start the scan, publish, and sending email process

## Note

### You need to re-create mySQL database to run this app.
#### Open mySQL cli, create a database call "btvn_b11"
#### Run 
#### Set 2 ENV variable DBUSER & DBPASS as user name and password to access mySQL sever

### You need to set an ENV SENDGRID_API with value ="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"
