package mailsender

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/AkronimBlack/stock/common"
	"github.com/valyala/bytebufferpool"
)

const (
	contentTypeHTML         = `text/html; charset=\"utf-8\"`
	mimeVer                 = "1.0"
	contentTransferEncoding = "base64"
)

//Email email structure expected from a file
type Email struct {
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	To      []string `json:"to"`
	Error   error
}

//CallOpts a list of flags
type CallOpts struct {
	Filename   string
	FromName   string
	FromAddr   string
	Username   string
	Password   string
	Hostname   string
	Port       int
	UseAuth    bool
	Concurrent bool
	Number     int
}

//Execute send emails
func Execute(opts CallOpts) {
	email := openAndRead(opts.Filename)
	config := Config{
		Host:     opts.Hostname,
		Username: opts.Username,
		Password: opts.Password,
		FromAddr: opts.FromAddr,
		FromName: opts.FromName,
		Port:     opts.Port,
		Auth:     opts.UseAuth,
	}
	e := make(chan Email)

	var s int
	var f int
	var counter int
	start := time.Now()

	for i := 0; i < opts.Number; i++ {
		fmt.Println("Sending email ", i)
		if opts.Concurrent {
			go sendMail(email, config, e)
		} else {
			rMail, err := sendMail(email, config, nil)
			if err != nil {
				fmt.Println("Sent email err: ")
				fmt.Println(rMail)
				f++
			} else {
				fmt.Println("Sent email : ")
				fmt.Println(rMail)
				s++
			}
			counter++
		}

	}
	if opts.Concurrent {
		for {
			select {
			case rEmail := <-e:
				fmt.Println("Sent email err: ")
				fmt.Println(rEmail)
				counter++
				if rEmail.Error != nil {
					fmt.Println("Sent email err: ")
					fmt.Println(rEmail)
					f++
				} else {
					fmt.Println("Sent email : ")
					fmt.Println(rEmail)
					s++
				}

				if counter == opts.Number {
					fmt.Printf("%d mails sent in %d ms. Successful: %d | Failed: %d", opts.Number, time.Now().Sub(start).Milliseconds(), s, f)
					return
				}
			}
		}
	}

	fmt.Printf("%d mails sent in %d ms. Successful: %d | Failed: %d", opts.Number, time.Now().Sub(start).Milliseconds(), s, f)

}

func openAndRead(message string) Email {
	file, err := os.Open(message)
	common.PanicOnError(err)
	byteValue, err := ioutil.ReadAll(file)
	fErr := file.Close()
	common.PanicOnError(fErr)
	common.PanicOnError(err)
	var email Email
	err = json.Unmarshal(byteValue, &email)
	common.PanicOnError(err)
	return email
}

func sendMail(email Email, config Config, e chan Email) (Email, error) {
	mailSender := NewMailer(config)
	email.Error = mailSender.Send(common.ReplacePlaceholder(email.Subject), []byte(common.ReplacePlaceholder(email.Body)), email.To)
	if e != nil {
		e <- email
	}
	return email, nil
}

// Config ....
// Host is the server mail host, IP or address. Port is the listening port.
// Username and Passwordare required if auth set to true.
// FromName and FromAddr are standard mailing values and are in the format of FromName: AkronimBlack FromAddr: akronimBlack@hostname.com
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	FromAddr string
	FromName string
	Auth     bool
}

// NewMailer creates and returns a new mail sender.
// if FromName not set the first part of FromAddr will be taken. If FromAddr empty and username is given (in the correct format, with an @) then
// username will be taken as FromAddr
func NewMailer(cfg Config) *Mailer {
	if cfg.FromAddr == "" && cfg.Username != "" && strings.Contains(cfg.Username, "@") {
		cfg.FromAddr = cfg.Username
	}

	if cfg.FromName == "" && cfg.FromAddr != "" {
		cfg.FromName = cfg.Username[0:strings.IndexByte(cfg.FromAddr, '@')]
	}

	if cfg.FromAddr == "" || cfg.FromName == "" {
		common.PanicOnError(fmt.Errorf("Missing from address or from alias"))
	}

	fromAddr := mail.Address{
		Name:    cfg.FromName,
		Address: cfg.FromAddr,
	}
	var auth smtp.Auth
	if cfg.Auth {
		if cfg.Username == "" || cfg.Password == "" {
			common.PanicOnError(fmt.Errorf("username or password missing"))
		}
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	return &Mailer{
		config:   cfg,
		useAuth:  cfg.Auth,
		fromAddr: fromAddr,
		auth:     auth,
	}
}

//Mailer struct
type Mailer struct {
	config   Config
	fromAddr mail.Address
	auth     smtp.Auth
	useAuth  bool
}

type stringWriter interface {
	WriteString(string) (int, error)
}

func writeHeaders(w stringWriter, subject string, body []byte, to []string) {
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "To", strings.Join(to, ",")))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "Subject", subject))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "MIME-Version", mimeVer))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", contentTypeHTML))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Transfer-Encoding", contentTransferEncoding))
	w.WriteString(fmt.Sprintf("\r\n%s", base64.StdEncoding.EncodeToString(body)))
}

var bufPool bytebufferpool.Pool

//Send send email through smtp
func (m *Mailer) Send(subject string, body []byte, to []string) error {
	buffer := bufPool.Get()
	defer bufPool.Put(buffer)

	if m.config.Host == "" || m.config.Port <= 0 {
		return fmt.Errorf("host or port missing")
	}

	buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "From", m.fromAddr.String()))
	writeHeaders(buffer, subject, body, to)

	return smtp.SendMail(
		fmt.Sprintf(fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)),
		m.auth,
		m.config.Username,
		to,
		buffer.Bytes(),
	)
}
