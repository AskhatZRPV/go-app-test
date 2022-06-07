package notifications

import (
	"fmt"
	"log"
	"os"
)

type Notification struct {
	username string
	spec     string
	time     string
}

func CreateNotification(username string, spec string, time string) *Notification {
	n := Notification{
		username: username,
		spec:     spec,
		time:     time,
	}
	return &n
}

func Write(s string) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println(s)
}

func (n *Notification) NotifyDay() {
	Write(fmt.Sprintf("| Привет %s! Напоминаем что вы записаны к %s завтра в %s", n.username, n.spec, n.time))
}

func (n *Notification) NotifyTwoHours() {
	Write(fmt.Sprintf("| Привет %s! Вам через 2 часа к %s в %s", n.username, n.spec, n.time))
}
