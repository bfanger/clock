## Clock

A "smart" clock written in [Go](https://golang.org) which run on a [2.8inch screen](https://www.waveshare.com/2.8inch-RPi-LCD-A.htm) connected to a [Raspberry Pi Zero](https://www.raspberrypi.org/)

Shows visual timers for school.

## Setup

```sh
apt-get install libsdl2{,-mixer,-image,-ttf,-gfx}-dev
go get -v github.com/veandco/go-sdl2/{sdl,img,ttf}
```

## Compile & Run

```sh
go build cmd/clock/clock.go
./clock
```
