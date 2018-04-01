# Clock

A "smart" clock written in [Go](https://golang.org) which runs on a [2.8inch screen](https://www.waveshare.com/2.8inch-RPi-LCD-A.htm) connected to a [Raspberry Pi Zero W](https://www.raspberrypi.org/)

## Goal

* Show time
* Alarm (for School)
* Show train delayes

## Setup

```sh
apt-get install libsdl2{,-mixer,-image,-ttf,-gfx}-dev
go get -v github.com/veandco/go-sdl2/{sdl,img,ttf}
go get -v github.com/bfanger/clock
```

## Architecture / Design

An abstraction on top of SDL to make an efficient ui without input.

### Laxy execution

Work is deferred until the result is needed.

The properties are just data, only when a Paint(\*sdl.Renderer) is called actual work is performed.
From loading images and fonts to rendering text.

The result of that work is cached so if no properties are changed drawing the next frame wil be fast.

### Concepts

```go
type Painter interface {
  Paint(r *sdl.Renderer) (*Texure, error)
  Destroy() error
}
```

A painter type can generate a texture based on it's properties.
It doesn't have a position and can't be displayed on its own.

```go
type Layer  {
  Name() string
  Render(*sdl.Renderer) error
}
```

To display something in the renderer you'll need a Layer.
The layer is responsible for rendering the texture(s)

When you add a layer to the Renderer, you'll also provide a zIndex which determines the ordering of the layers.
layers with the same zIndex are rendered in the order they are added.
