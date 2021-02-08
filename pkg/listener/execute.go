package listener

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Azure/go-amqp"
)

var wg sync.WaitGroup

//LaunchListener ...
func LaunchListener(hostname, port, username, password string, topics []string) {
	log.Println("Starting")
	// Create client
	client, err := amqp.Dial("amqp://"+hostname+":"+port,
		amqp.ConnSASLPlain(username, password),
	)
	if err != nil {
		log.Fatal("Dialing AMQP server:", err)
	}
	defer client.Close()

	log.Println("Connection made ", "amqp://"+hostname+":"+port)
	// Open a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Creating AMQP session:", err)
	}

	log.Println("New Session")
	ctx := context.Background()

	for _, topic := range topics {
		receiver, err := session.NewReceiver(
			amqp.LinkSourceAddress(topic),
			amqp.LinkCredit(10),
		)
		if err != nil {
			log.Fatal("Failed createing link:", err)
		}
		defer func() {
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			receiver.Close(ctx)
			cancel()
		}()
		wg.Add(1)
		go Listen(ctx, receiver, topic)
	}
	wg.Wait()
}

//Listen ...
func Listen(ctx context.Context, receiver *amqp.Receiver, topicName string) {
	defer wg.Done()
	log.Println("Listening on ", receiver.Address())
	for {
		// Receive next message
		msg, err := receiver.Receive(ctx)
		if err != nil {
			log.Fatal("Reading message from AMQP:", err)
		}

		// Accept message
		if err := msg.Accept(ctx); err != nil {
			log.Println("Error: ", err.Error())
		}
		log.Println("----------------------------------------------------------------------------------------------")
		log.Println("-------------------------    GOT NEW MESSAGE -------------------------------------------------")
		log.Println("----------------------------------------------------------------------------------------------\n\n")
		log.Println("Receiver: ", receiver.Address(), "\n\n")
		log.Println("----------------------------------------------------------------------------------------------")
		log.Println("Application Properties")
		log.Println(msg.ApplicationProperties, "\n\n")
		log.Println("----------------------------------------------------------------------------------------------")
		log.Println("Properties")
		log.Println(msg.Properties, "\n\n")
		log.Println("----------------------------------------------------------------------------------------------")
		log.Println("Annotations")
		log.Println(msg.Annotations, "\n\n")
		log.Println("----------------------------------------------------------------------------------------------")
		log.Println("DeliveryAnnotations")
		log.Println(msg.DeliveryAnnotations, "\n\n")
		log.Println("----------------------------------------------------------------------------------------------")
		log.Println("Text payload")
		log.Println(msg.Value, "\n\n")
		log.Println("----------------------------------------------------------------------------------------------")
		log.Println("Binary payload")
		for _, data := range msg.Data {
			log.Println(string(data))
		}

	}
}
