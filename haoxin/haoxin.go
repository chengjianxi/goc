package haoxin

import (
	"fmt"
	"math/rand"
	"time"
)

func NewSid(moduleId string) string {
	seq := rand.Intn(1000000)
	t := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%06d", moduleId, t, seq)
}

func NewTid(moduleId string) string {
	seq := rand.Intn(1000000)
	t := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%06d", moduleId, t, seq)
}
