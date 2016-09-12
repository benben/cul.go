package main

import (
	"fmt"
	"github.com/tarm/serial"
	"regexp"
	"strings"
	"time"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 38400}

	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println(err)
	}

	n, err := s.Write([]byte("X21\n"))
	if err != nil {
		fmt.Println(err)
	}

	var str string
	r, _ := regexp.Compile(`^K\d{8}[A-Z0-9]{2}`)

	for strings.Count(str, "") <= 1 || !(strings.Contains(str, "\n") && r.MatchString(str)) {
		buf := make([]byte, 128)
		n, _ = s.Read(buf)
		str += string(buf[:n])
		raw := parseRaw(str)

		fmt.Printf("{\"raw\": \"%v\", \"temp\": %v, \"hum\": %v, \"created_at\": \"%v\"}\n", raw,
			parseTemp(raw), parseHum(raw), time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	}
}

func parseRaw(str string) string {
	values := strings.SplitAfterN(str, "\r\n", 1)
	return strings.Replace(values[len(values)-1], "\r\n", "", -1)
}

func parseTemp(str string) string {
	return fmt.Sprintf("%v%v.%v", string(str[6]), string(str[3]), string(str[4]))
}

func parseHum(str string) string {
	return fmt.Sprintf("%v%v.%v", string(str[7]), string(str[8]), string(str[5]))
}
