#!/usr/bin/env python3

from sdl2.ext import init, Renderer, TEXTURE, SpriteFactory, World, Entity
from sdl2.sdlttf import TTF_Init, TTF_Quit
from window import Window
from klok import Klok
from eventloop import eventloop
from constants import RESOURCES
from nsapi import departures


def main():
    init()
    TTF_Init()

    window = Window()
    window.show()

    renderer = Renderer(window)
    factory = SpriteFactory(TEXTURE, renderer=renderer)
    world = World()

    spriteRenderer = factory.create_sprite_render_system()

    background = Entity(world)
    background.texture = factory.from_image(
        RESOURCES.get_path('background.png'))
    Klok(world, renderer=renderer)
    world.add_system(spriteRenderer)

    print(departures())  # todo show in interface

    try:
        eventloop(world, renderer=renderer)
    except KeyboardInterrupt:
        pass

    TTF_Quit()
    quit()
    return 0


if __name__ == '__main__':
    main()
