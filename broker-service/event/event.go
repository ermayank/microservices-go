package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //Name
		"topic",      //Type
		true,         // Durable ?
		false,        //auto-deleted?
		false,        //Internal ?
		false,        //no-wait ?
		nil,          //Arguments ?
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //Name
		false, //Durable ?
		false, //delete when unused ?
		true,  //exclusive ?
		false, //no-wait ?
		nil,   // Arguments ?
	)
}
