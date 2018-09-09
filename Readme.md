# Clock

A "smart" clock written in [Go](https://golang.org) which runs on a [2.8inch screen](https://www.waveshare.com/2.8inch-RPi-LCD-A.htm) connected to a [Raspberry Pi Zero W](https://www.raspberrypi.org/)

## Goal

- Show time
- Alarm (for School)
- Show train delayes

## Setup

```sh
apt-get install libsdl2{,-mixer,-image,-ttf,-gfx}-dev
go get -v github.com/veandco/go-sdl2/{sdl,img,ttf}
go get -v github.com/bfanger/clock
```

## Architecture / Design

An abstraction on top of SDL to make an efficient event-based ui.

### Lazy execution

Work is deferred until the result is needed. This allows us to freely change individual properties of a layer without causing an updated texture per change.

The actual work is performed when a `Image(\*sdl.Renderer)` or `Compose(\*sdl.Renderer)` is called.
The result of that work is cached, so drawing the next frame will be even faster.

### Concepts

```go
type Imager interface {
  Image(r *sdl.Renderer) (*Image, error)
}
```

An `Imager` can generate a image/texture based on it's properties.
It doesn't have a position and can't be displayed on its own.

```go
type Composer interface {
  Compose(*sdl.Renderer) error
}
```

To display something in the renderer you'll need a Composer.
The Composer is responsible for rendering the texture(s) onto the screen

#### Engine

Composers are added to the Engine and are rendered automaticly.
All UI operation should be wrapped in a `engine.Go()` closure which are batched and executed in the main/ui thread.
A useful side-effect of calling engine.Go is that it will trigger a re-render.

# Sidenote

I use a custom display driver [fbcp-ili9341](https://github.com/juj/fbcp-ili9341) which i configured with:

```sh
cmake -DILI9341=ON -DSPI_BUS_CLOCK_DIVISOR=6 -DGPIO_TFT_RESET_PIN=27 -DGPIO_TFT_DATA_CONTROL=22 -DSTATISTICS=0 ..
```
