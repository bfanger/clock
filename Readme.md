# Clock

A "smart" clock written in [Go](https://golang.org) which runs on a [4.0 inch screen](https://shop.pimoroni.com/products/hyperpixel-4) connected to a [Raspberry Pi Zero W](https://www.raspberrypi.org/)

## Goal

- Show time
- Alarm (for School)
- Show train delayes

## Setup

```sh
go install github.com/mitranim/gow@latest
apt-get install libsdl2{,-mixer,-image,-ttf,-gfx}-dev
go get -v github.com/veandco/go-sdl2/{sdl,img,ttf}
go get -v github.com/bfanger/clock
```

## MacOS dev setup

```
brew install pkg-config sdl2 sdl2_image sdl2_ttf
go install github.com/mitranim/gow@latest
```

## Architecture / Design

An abstraction on top of SDL to make an efficient event-based ui.

### Lazy execution

Allow work to be deferred until the result are needed. This allows freely changing individual properties of a layer without causing an updated texture for every change.

The actual work is performed when a `Image(\*sdl.Renderer)` via `Compose(\*sdl.Renderer)` is called.
The resulting texture is cached, so drawing the next frame will be even faster.

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

# Getting latest SDL on Raspberry OS

Download latest [libsdl2 source](https://www.libsdl.org/download-2.0.php)
(2.0.14 at the time of writing)

```
unzip SDL2-2.0.14.zip
cd SDL2-2.0.14
./configure
make
sudo chmod -R pi /usr/local
make install
```

Download [SDL_image](https://www.libsdl.org/projects/SDL_image/)

```
sudo apt install libjpeg-dev libtiff-dev
tar -xvf ./SDL2_image-2.0.5.tar.gz
cd SDL2_image-2.0.5
./configure
make
make install
```

and [SDL_ttf](https://www.libsdl.org/projects/SDL_ttf/)

```
tar -xvf ./SDL2_ttf-2.0.15.tar.gz
cd SDL2_ttf-2.0.15
./configure
make
make install
```

Remove the packaged version:

```
sudo apt-get remove libsdl2{,-mixer,-image,-ttf,-gfx}-dev
sudo apt-get remove libsdl2{,-mixer,-image,-ttf}-2.0
```
