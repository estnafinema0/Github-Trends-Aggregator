package email

import (
	"log"
	"time"
	"text/template"
	"bytes"

	"gopkg.in/gomail.v2"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/store"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/config"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
)

func StartEmail(st *store.Store, l *log.Logger) {
	ticker := time.NewTicker(config.NotifyInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		l.Println("Starting email notification...")
		st.Lock()
		d := gomail.NewDialer(config.SMTPServer, config.SMTPPort, config.SMTPUsername, config.SMTPKey)
		s, err := d.Dial()

		if err != nil {
			l.Printf("Error connecting to dialer: %v\n", err)
		} else {
			t := template.Must(template.ParseFiles("static/email.html"))
			m := gomail.NewMessage()
			for _, e := range st.NotifsList {
				// m.SetHeader("Message-ID", "<" + uuid.New().String() + "@mail.gmail.com>")
				m.SetHeader("From", config.NotifyEmail)
				m.SetHeader("To", e)
				m.SetHeader("Subject", "Trending repos")
				var tpl bytes.Buffer
				st.Unlock()
				rs := st.GetReposFiltered("", "stars")
				st.Lock()
				data := struct { Repos []models.Repository } { rs }
				if err := t.Execute(&tpl, data); err != nil {
					l.Printf("Could execute template for %q: %v", e, err)
				}
				m.SetBody("text/html", tpl.String())
				if err := gomail.Send(s, m); err != nil {
					l.Printf("Could not send email to %q: %v", e, err)
				}
				m.Reset()
			}
		}
		st.Unlock()
		l.Printf("Sent %d notifications\n", len(st.NotifsList))
	}
}
