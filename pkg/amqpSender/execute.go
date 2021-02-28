package amqpsender

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/AkronimBlack/stock/common"
	"github.com/Azure/go-amqp"
)

//SendMessages ....
func SendMessages(hostname, port, filename, username, password string, topics []string, num int) {

	client, err := amqp.Dial("amqp://"+hostname+":"+port,
		amqp.ConnSASLPlain(username, password),
	)
	common.PanicOnError(err)

	log.Println("Connection made: ", "amqp://"+hostname+":"+port)

	session, err := client.NewSession()
	common.PanicOnError(err)

	file, err := os.Open(filename)
	common.PanicOnError(err)

	byteValue, err := ioutil.ReadAll(file)
	common.PanicOnError(err)

	var config EventMessage
	err = json.Unmarshal(byteValue, &config)
	common.PanicOnError(err)

	for _, topic := range topics {
		log.Println("New Session")
		ctx := context.Background()
		sender, err := session.NewSender(amqp.LinkTargetAddress(topic))
		for i := 0; i < num; i++ {
			sendJSON, subErr := json.Marshal(interateAndBuild(config.Payload))
			if subErr != nil {
				log.Println(err.Error())
				return
			}
			sendProps := interateAndBuild(config.ApplicationProperties)
			sendPropsJSON, subErr := json.Marshal(sendProps)
			if subErr != nil {
				log.Println(err.Error())
				return
			}
			log.Println("----- Headers -----")
			log.Println(string(sendPropsJSON))
			log.Println("----- Payload -----")
			log.Println(string(sendJSON))
			log.Printf("Sending message: %d to topic %s", i+1, topic)
			err = sender.Send(ctx, NewMessage(sendJSON, sendProps))
			if err != nil {
				log.Fatal("Error :", err)
			}
			log.Println("Message sent")
		}
		sender.Close(ctx)
	}
	err = client.Close()
	err = file.Close()
	common.PanicOnError(err)
}

func interateAndBuild(data map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(data))

	for i, v := range data {
		newMap[i] = v
		switch newMap[i].(type) {
		case string:
			newMap[i] = common.ReplacePlaceholder(newMap[i].(string))
		case map[string]string:
			castValue := newMap[i].(map[string]string)
			for j, k := range castValue {
				castValue[j] = common.ReplacePlaceholder(k)
			}
			newMap[i] = castValue
		case map[string]interface{}:
			interateAndBuild(newMap[i].(map[string]interface{}))
		default:
			log.Printf("%v is unknown \n ", v)
		}
	}
	return newMap
}

func NewMessage(data []byte, props map[string]interface{}) *amqp.Message {
	return &amqp.Message{
		Data:                  [][]byte{data},
		ApplicationProperties: props,
	}
}

//EventMessage message bluprint to send to bus
type EventMessage struct {
	Payload               map[string]interface{} `json:"payload"`
	ApplicationProperties map[string]interface{} `json:"properties"`
}
