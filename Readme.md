# Clock

A "smart" clock written in [Go](https://golang.org) which runs on a [4.0 inch screen](https://shop.pimoroni.com/products/hyperpixel-4) connected to a [Raspberry Pi Zero W](https://www.raspberrypi.org/)

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

Allow work to be deferred until the result are needed. This allows freely changing individual properties of a layer without causing an updated texture for every change.

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

Composers are added to Scene in the Engine and are rendered automaticly.
All UI operation should be wrapped in a `engine.Go()` closure which are batched and executed in the main/ui thread.
A useful side-effect of calling engine.Go is that it will trigger a re-render.
