package trabbitmq

import (
    "log"
    "errors"
    
    "github.com/streadway/amqp"
)

type RabbitmqCharacter struct {
    ExchangeName string
    QueueName    string
    RouteKey     string
}

func setupQueueInternal(ch *amqp.Channel, setting *RabbitmqCharacter) error {
    
    err := ch.ExchangeDeclare (
            setting.ExchangeName,  // name
            "direct",              // kind
            true,                  // durable
            false,                 // autoDelete
            false,                 // internal
            false,                 // no-wait
            nil)                   // arguments
    if err != nil {
        log.Print("Error - RabbitMQ ExchangeDeclare: ", err)
        return err
    }
    
    _, err = ch.QueueDeclare(
            setting.QueueName,     // name
            true,                  // durable
            false,                 // delete when unused
            false,                 // exclusive
            false,                 // no-wait
            nil)                   // arguments
    if err != nil {
        log.Print("Error - RabbitMQ QueueDeclare: ", err)
        return err
    }
    
    err = ch.QueueBind(
            setting.QueueName,    // queue name
            setting.RouteKey,     // routing key
            setting.ExchangeName, // exchange
            false,                // no-wait
            nil)                  // arguments
    if err != nil {
        log.Print("Error - RabbitMQ QueueBind: ", err)
        return err
    }
    
    return nil
}

func SetupPublishQueue(ch *amqp.Channel, setting *RabbitmqCharacter) error {
  
    err := setupQueueInternal(ch, setting)
    if err != nil {
        return err
    }
    
    // Select Mode
    err = ch.Confirm( false ); // no-wait
    if err != nil {  
        log.Print("Error - RabbitMQ Publish Confirm: ", err)
        return err
    }
    
    return nil
}

func SetupGetQueue(ch *amqp.Channel, setting *RabbitmqCharacter) error {

    if err := setupQueueInternal(ch, setting); err != nil {
        return err
    }

    return nil
}

func SetupConsumeQueue(ch *amqp.Channel, setting *RabbitmqCharacter) (<-chan amqp.Delivery, error) {

    err := setupQueueInternal(ch, setting)
    if err != nil {
        return nil, err
    }
    
    err = ch.Qos(
            1,           // prefetchCount
            0,           // prefetchSize 
            true)        // global
    if err != nil {  
        log.Print("Error - RabbitMQ Consume Qos: ", err)
        return nil, err
    }
    
    msg_ch, err := ch.Consume(
            setting.QueueName,  // queue
            "",                 // consumer
            false,              // auto ack
            false,              // exclusive
            false,              // no local
            false,              // no wait
            nil)                // args
    if err != nil {  
        log.Print("Error - RabbitMQ Consume: ", err)
        return nil, err
    }
    
    return msg_ch, nil
}

// message

func PublisMessage(ch *amqp.Channel, setting *RabbitmqCharacter, msg string) error {

    err := ch.Publish(
            setting.ExchangeName,     // exchange
            setting.RouteKey,         // routing key
            true,                     // mandatory
            false,                    // immediate
            amqp.Publishing {
                DeliveryMode: amqp.Persistent,
                Body:        []byte(msg),
            })
    
    if err != nil {
        log.Print("Error - RabbitMQ Publish: ", err)
        return err
    }
    
    return nil
}

func GetMessage(ch *amqp.Channel, queue string, autoAck bool) ( []byte, uint64, error) {

    msg, ok, err := ch.Get(queue, autoAck)
    if err != nil {
        log.Print("Error - RabbitMQ GetMessage: ", err)    
        return nil, 0, err
    }
    
    if ok != true {
        log.Print("Error - RabbitMQ GetMessage failed")
        return nil, 0, errors.New("Empty Msg")
    }
    
    return msg.Body, msg.DeliveryTag, nil
}

// The consumer tag is local to a channel, so two clients can use the same consumer tags.
func Ack(ch *amqp.Channel, tag uint64) error {

    return ch.Ack(
            tag, 
            true)    // multiple
}

func Nack(ch *amqp.Channel, tag uint64, requeue bool) error {

    return ch.Nack(
            tag,
            true,    // multiple
            requeue)   
}