from sdl2.ext import Entity, TEXTURE, SpriteFactory
from constants import SCREEN_WIDTH, SCREEN_HEIGHT, DEPTH_BRIGHTNESS


class Brightness(Entity):
    """Brightness implemented in software, using a black transparent overlay"""

    def __init__(self, world, *args, **kwargs):
        if "renderer" not in kwargs:
            raise ValueError("you have to provide a renderer= argument")
        renderer = kwargs['renderer']
        factory = SpriteFactory(TEXTURE, renderer=renderer)
        sprite = factory.from_color(
            color=0xAA000000,
            size=(SCREEN_WIDTH, SCREEN_HEIGHT),
            masks=(0xFF000000,
                   0x00FF0000,
                   0x0000FF00,
                   0x000000FF)
        )
        sprite.depth = DEPTH_BRIGHTNESS
        self.sprite = sprite
