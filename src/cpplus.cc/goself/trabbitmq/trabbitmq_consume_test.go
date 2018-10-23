package trabbitmq

import (
    "log"
    "testing"

    "github.com/streadway/amqp"
)

// something already defined in trabbitmq_test.go   

func consume_message() (string, error) {

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
    
    msg_ch, err := SetupConsumeQueue(ch, &rabbitCharacter)
    if err != nil {
        log.Print("Setup RabbitMQ fail: ", err)
        return "", err
    }
    
    select {
        case dt, ok := <- msg_ch:
        if !ok {
            log.Println("Read from channel failed.")
            break
        }

        msg_content := dt.Body;
        msg_tag     := dt.DeliveryTag
        
        log.Printf("Recv: %s", msg_content)
        Ack(ch, msg_tag)
        
        return string(msg_content), nil
    }
        
    return "", nil
}   


func TestConsume(t *testing.T) {
    
    msg := "TestMessageForConsume!"
    
    err := send_message(msg)
    if err != nil {
        t.Fatal(err)
    }
    
    expected_msg, err := consume_message()
    if err != nil {
        t.Fatal(err)
    }
    
    if msg != expected_msg {
        t.Errorf("Expected %s, but get: %s!", msg, expected_msg)
    }

}
