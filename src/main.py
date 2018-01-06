#!/usr/bin/env python3

from ctypes import py_object, cast, POINTER, byref
from sdl2.ext import init, Renderer, TEXTURE, SpriteFactory, World, Entity
from sdl2.sdlttf import TTF_Init, TTF_Quit
from sdl2 import (SDLK_ESCAPE, SDL_Event, SDL_WaitEvent, SDL_QUIT,
                  SDL_KEYUP, SDL_USEREVENT, SDL_BLENDMODE_BLEND)
from window import Window
from klok import Klok
from constants import SCREEN_WIDTH, SCREEN_HEIGHT
import sdl2.ext as sdl2ext

from constants import RESOURCES


def main():
    init()
    TTF_Init()

    window = Window()
    window.show()

    renderer = Renderer(window)
    renderer.blendmode = SDL_BLENDMODE_BLEND
    factory = SpriteFactory(TEXTURE, renderer=renderer)
    world = World()

    spriteRenderer = factory.create_sprite_render_system()

    # Background
    background = Entity(world)
    backgroundSprite = factory.from_image(
        RESOURCES.get_path('background.png'))
    background.sprite = backgroundSprite

    # Time
    Klok(world, renderer=renderer)

    # Brighness
    brighness = Entity(world)
    brighness.sprite = factory.from_color(
        color=0xAA000000,
        size=(SCREEN_WIDTH, SCREEN_HEIGHT),
        masks=(0xFF000000,
               0x00FF0000,
               0x0000FF00,
               0x000000FF))

    world.add_system(spriteRenderer)

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
