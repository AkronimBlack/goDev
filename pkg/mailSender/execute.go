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

	"github.com/AkronimBlack/dev-tools/common"
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

//Execute send emails
func Execute(message, from, username, password, hostname string, auth, concurrent bool, num int) {
	email := openAndRead(message)
	config := Config{
		Host:     hostname,
		Username: username,
		Password: password,
		FromAddr: from,
		Port:     25,
		Auth:     auth,
	}
	e := make(chan Email)

	var s int
	var f int
	var counter int
	start := time.Now()

	for i := 0; i < num; i++ {
		fmt.Println("Sending email ", i)
		if concurrent {
			go sendMail(from, email, config, e)
		} else {
			rMail, err := sendMail(from, email, config, nil)
			if err != nil {
				fmt.Println("Sent email err: ")
				fmt.Println(rMail)
				counter = counter + 1
				f = f + 1
			} else {
				fmt.Println("Sent email : ")
				fmt.Println(rMail)
				counter = counter + 1
				s = s + 1
			}
		}

	}
	if concurrent {
		for {
			select {
			case rEmail := <-e:
				fmt.Println("Sent email err: ")
				fmt.Println(rEmail)
				counter = counter + 1
				if rEmail.Error != nil {
					fmt.Println("Sent email err: ")
					fmt.Println(rEmail)
					counter = counter + 1
					f = f + 1
				} else {
					fmt.Println("Sent email : ")
					fmt.Println(rEmail)
					counter = counter + 1
					s = s + 1
				}

				if counter == num {
					fmt.Printf("%d mails sent in %d ms. Successful: %d | Failed: %d", num, time.Now().Sub(start).Milliseconds(), s, f)
					return
				}
			}
		}
	}

	fmt.Printf("%d mails sent in %d ms. Successful: %d | Failed: %d", num, time.Now().Sub(start).Milliseconds(), s, f)

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

func sendMail(from string, email Email, config Config, e chan Email) (Email, error) {
	mailSender := NewMailer(config)
	email.Error = mailSender.Send(common.ReplacePlaceholder(email.Subject), []byte(common.ReplacePlaceholder(email.Body)), email.To)
	if e != nil {
		e <- email
	}
	return email, nil
}

//Config ....
type Config struct {
	// Host is the server mail host, IP or address.
	Host string
	// Port is the listening port.
	Port int
	// Username is the auth username@domain.com for the sender.
	Username string
	// Password is the auth password for the sender.
	Password string
	// FromAddr is the 'from' part of the mail header, it overrides the username.
	FromAddr string
	// FromAlias is the from part, if empty this is the first part before @ from the Username field.
	FromAlias string
	Auth      bool
}

// NewMailer creates and returns a new mail sender.
func NewMailer(cfg Config) *Mailer {
	m := &Mailer{config: cfg}
	addr := cfg.FromAddr
	if addr == "" {
		addr = cfg.Username
	}

	if cfg.FromAlias == "" {
		if cfg.Username != "" && strings.Contains(cfg.Username, "@") {
			m.fromAddr = mail.Address{Name: cfg.Username[0:strings.IndexByte(cfg.Username, '@')], Address: addr}
		}
	} else {
		m.fromAddr = mail.Address{Name: cfg.FromAlias, Address: addr}
	}
	m.useAuth = cfg.Auth

	return m
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

	if m.useAuth {
		cfg := m.config
		if cfg.Host == "" || cfg.Port <= 0 {
			return fmt.Errorf("username, password, host or port missing")
		}
		m.auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	fullhost := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "From", m.fromAddr.String()))
	writeHeaders(buffer, subject, body, to)

	return smtp.SendMail(
		fmt.Sprintf(fullhost),
		m.auth,
		m.config.Username,
		to,
		buffer.Bytes(),
	)
}
