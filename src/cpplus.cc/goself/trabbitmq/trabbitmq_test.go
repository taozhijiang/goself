package trabbitmq

import (
    "log"
    "testing"

    "github.com/streadway/amqp"
)


// Cfg and parameters define here
var connStr string = "amqp://rabbituser:qj6RCT4GW9cG@119.23.68.224:5672/demo"
var rabbitCharacter RabbitmqCharacter = 
        RabbitmqCharacter {"paybank_exchange", "notify_general", "notify_general_key" }


func send_message(content string) error {

    conn, err := amqp.Dial(connStr)
    if err != nil {
        log.Print("Failed to connect to RabbitMQ:", err)
        return err
    }
    defer conn.Close()
    
    ch, err := conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()
    
    if err = SetupPublishQueue(ch, &rabbitCharacter); err != nil {
        log.Print("Setup RabbitMQ fail:", err)
        return err
    }
    
    if PublisMessage(ch, &rabbitCharacter, content); err != nil {
        log.Printf("Failed to send message.")
        return err
    }
    
    log.Printf(" [x] Sent %s OK", content)
    
    return nil
}     

func get_message() (string, error) {

    conn, err := amqp.Dial(connStr)
    if err != nil {
        log.Print("Failed to connect to RabbitMQ: ", err)
        return "", err
    }
    defer conn.Close()
    
    ch, err := conn.Channel()
    if err != nil {
        return "", err
    }
    defer ch.Close()
    
    if err = SetupGetQueue(ch, &rabbitCharacter); err != nil {
        log.Print("Setup RabbitMQ fail: ", err)
        return "", err
    }
    
    msg, tag, err := GetMessage(ch, rabbitCharacter.QueueName, false)
    if err != nil {
        log.Print("GetMessage RabbitMQ fail: ", err)
        return "", err
    }
    
    log.Printf("Recv: %s", msg)
    Ack(ch, tag)
    
    return string(msg), nil
}   


func TestGet(t *testing.T) {
    
    msg := "TestMessageForGet"
    
    err := send_message(msg)
    if err != nil {
        t.Fatal(err)
    }
    
    expected_msg, err := get_message()
    if err != nil {
        t.Fatal(err)
    }
    
    if msg != expected_msg {
        t.Errorf("Expected %s, but get: %s!", msg, expected_msg)
    }

}
