package sender

import (
	"app/internal/models"
	"bytes"
	"html/template"
	"net/smtp"
)

type Sender struct {
	SenderConfig
}

type SenderConfig struct {
	From     string `env:"SMTP_FROM" env-default:"noreply@example.com"`
	Password string `env:"SMTP_PASSWORD" env-default:""`
	SmtpHost string `env:"SMTP_HOST" env-default:"localhost"`
	SmtpPort string `env:"SMTP_PORT" env-default:"1025"`
}

func New(cfg SenderConfig) *Sender {
	return &Sender{cfg}
}

func (s *Sender) SendReminder(to string, event models.Event) error {
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            .reminder { font-family: Arial; padding: 20px; }
            .message { font-size: 18px; color: #333; }
            .date { color: #666; }
        </style>
    </head>
    <body>
        <div class="reminder">
            <h2>Напоминание о событии</h2>
            <p class="message">{{.Message}}</p>
            <p class="date">Создано: {{.CreatedAt.Format "02.01.2006 15:04"}}</p>
        </div>
    </body>
    </html>
    `

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, event); err != nil {
		return err
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	msg := "From: " + s.From + "\n" +
		"To: " + to + "\n" +
		"Subject: Напоминание о событии\n" +
		headers + "\n\n" +
		body.String()

	return smtp.SendMail(s.SmtpHost+":"+s.SmtpPort, nil, s.From, []string{to}, []byte(msg))
}
