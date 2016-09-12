# cul.go
reading CULFW devices with Go and returning parsed values as JSON string

## Installation

    glide install
    go build -o cul
    ./cul
    # Output: {"raw": "K0100327506", "temp": 20.0, "hum": 75.3, "created_at": "2016-09-12T07:16:42+0000"}

## More Info

I implemented it in Ruby also, its here: https://github.com/benben/cul.rb

see http://culfw.de/culfw.html

## License
see LICENSE file
