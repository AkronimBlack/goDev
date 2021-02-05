package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/AkronimBlack/dev-tools/common"
	"github.com/Azure/go-amqp"
	"github.com/bxcodec/faker"
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
			newMap[i] = convertString(newMap[i].(string))
		case map[string]string:
			castValue := newMap[i].(map[string]string)
			for j, k := range castValue {
				castValue[j] = convertString(k)
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

func trimExcessFat(value string, left string, right string) string {
	stringValue := strings.TrimLeft(value, left)
	stringValue = strings.TrimRight(stringValue, right)
	return stringValue
}

func convertString(value string) string {
	if strings.Contains(value, "{{") {
		val := trimExcessFat(value, "{{", "}}")
		splitData := strings.Split(val, ".")
		if splitData[0] == "faker" {
			return newFaked(splitData[1])
		}
	}
	return value
}

func newFaked(key string) string {
	v := FakerOptions{}
	err := faker.FakeData(&v)
	if err != nil {
		log.Println("Problem faking data")
		fmt.Println(err)
	}
	rv := reflect.ValueOf(v)
	fv := rv.FieldByName(key)
	return fv.String()
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

type FakerOptions struct {
	Email           string `faker:"email"`
	PhoneNumber     string `faker:"phone_number"`
	URL             string `faker:"url"`
	UserName        string `faker:"username"`
	TitleMale       string `faker:"title_male"`
	TitleFemale     string `faker:"title_female"`
	FirstName       string `faker:"first_name"`
	FirstNameMale   string `faker:"first_name_male"`
	FirstNameFemale string `faker:"first_name_female"`
	LastName        string `faker:"last_name"`
	Name            string `faker:"name"`
	Date            string `faker:"date"`
	Time            string `faker:"time"`
	MonthName       string `faker:"month_name"`
	Year            string `faker:"year"`
	DayOfWeek       string `faker:"day_of_week"`
	DayOfMonth      string `faker:"day_of_month"`
	Timestamp       string `faker:"timestamp"`
	Century         string `faker:"century"`
	TimeZone        string `faker:"timezone"`
	TimePeriod      string `faker:"time_period"`
	Word            string `faker:"word"`
	Sentence        string `faker:"sentence"`
	Paragraph       string `faker:"paragraph"`
	Currency        string `faker:"currency"`
	UUID            string `faker:"uuid_digit"`
}
