package main

import (
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	command *string
	verbose *bool
	str     string
)

func main() {
	command = flag.String("c", "X21", "the command which is send to CULFW")
	verbose = flag.Bool("v", false, "verbose output")

	flag.Parse()

	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 38400}
	s := initializeCul(c)

	if *command == "X21" && !*verbose {
		readAndParse(s)
	} else {
		read(s)
	}
}

func read(s *serial.Port) {
	for {
		log.Println(readRaw(s))
	}
}

func readAndParse(s *serial.Port) {
	// CULFW always returns a "K", followed by 8 digits, followed by a hex number
	r, _ := regexp.Compile(`^K\d{8}[A-Z0-9]{2}`)

	// read from serial as long as we didn't receive something already
	// or it didn't end with \n and isn't a full value yet
	for strings.Count(str, "") <= 1 || !(strings.Contains(str, "\n") && r.MatchString(str)) {
		str += readRaw(s)
		raw := parseRaw(str)

		fmt.Printf("{\"raw\": \"%v\", \"temp\": %v, \"hum\": %v, \"created_at\": \"%v\"}\n", raw,
			parseValue(raw, 6, 3, 4), parseValue(raw, 7, 8, 5), time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	}
}

func readRaw(s *serial.Port) string {
	buf := make([]byte, 128)
	n, _ := s.Read(buf)
	return string(buf[:n])
}

// open the serial connection and send appropriate command
// for more commands see http://culfw.de/commandref.html
func initializeCul(c *serial.Config) *serial.Port {
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println(err)
	}

	if *verbose {
		s.Write([]byte("V\n"))
	}

	_, err = s.Write([]byte(*command + "\n"))
	if err != nil {
		fmt.Println(err)
	}

	return s
}

// Sometimes the CULFW returns up to ten values if it didn't get read for a longer period.
// In that case, take the last value and return it
func parseRaw(str string) string {
	values := strings.SplitAfterN(str, "\r\n", 1)
	return strings.Replace(values[len(values)-1], "\r\n", "", -1)
}

// stolen from how fhem parses the data
// https://github.com/mhop/fhem-mirror/blob/master/fhem/FHEM/14_CUL_WS.pm#L146-L156
func parseValue(str string, c1 int, c2 int, c3 int) float64 {
	floatStr := fmt.Sprintf("%v%v.%v", string(str[c1]), string(str[c2]), string(str[c3]))
	float, _ := strconv.ParseFloat(floatStr, 64)
	return float
}
