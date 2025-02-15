package config

import (
	"time"
	"os"
	"strconv"
)

// FetchInterval is the interval at which trends are fetched
const FetchInterval = time.Minute * 1
const NotifyInterval = time.Hour * 24
const HistoryLength = 100
var NotifyEmail string
var SMTPServer string
var SMTPPort int
var SMTPUsername string
var SMTPKey string

func LoadSecrets() bool {
	var exists bool
	SMTPServer, exists = os.LookupEnv("SMTP_SERVER")
	if !exists { return false }
	var port_string string
	port_string, exists = os.LookupEnv("SMTP_PORT")
	SMTPPort, _ = strconv.Atoi(port_string)
	if !exists { return false }
	SMTPUsername, exists = os.LookupEnv("SMTP_USERNAME")
	if !exists { return false }
	SMTPKey, exists = os.LookupEnv("SMTP_KEY")
	if !exists { return false }
	NotifyEmail, exists = os.LookupEnv("FROM_EMAIL")
	if !exists { return false }
	return true
}
