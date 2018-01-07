#!/usr/bin/env python3

from ctypes import py_object, cast, POINTER, byref
from sdl2.ext import init, Renderer, TextureSpriteRenderSystem, World
from sdl2.sdlttf import TTF_Init, TTF_Quit
from sdl2 import (SDLK_ESCAPE, SDL_Event, SDL_WaitEvent, SDL_QUIT,
                  SDL_KEYUP, SDL_USEREVENT, SDL_BLENDMODE_BLEND)
from window import Window
from entities import Time, Background, Brightness
import sdl2.ext as sdl2ext


def main():
    init()
    TTF_Init()

    window = Window()
    window.show()

    renderer = Renderer(window)
    renderer.blendmode = SDL_BLENDMODE_BLEND
    world = World()

    Background(world, renderer=renderer)
    Time(world, renderer=renderer)
    Brightness(world, renderer=renderer)

    world.add_system(TextureSpriteRenderSystem(renderer))

    try:
        event = SDL_Event()
        while True:
            ret = SDL_WaitEvent(byref(event))
            if ret == 0:
                raise sdl2ext.SDLError()

            if event.type == SDL_QUIT:
                # Closed window
                break
            elif (event.type == SDL_KEYUP and
                  event.key.keysym.sym == SDLK_ESCAPE):
                # Pressed ESC
                break
            elif event.type == SDL_USEREVENT:
                # Timer event (probably)
                entity = cast(event.user.data1, POINTER(
                    py_object)).contents.value
                entity.update()
            else:
                print("Event:", event.type)

            renderer.clear()
            world.process()

    except KeyboardInterrupt:
        pass

    TTF_Quit()
    quit()
    return 0


if __name__ == '__main__':
    main()
