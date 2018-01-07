from sdl2.ext import Entity, TEXTURE, SpriteFactory
# from sdl2 import
from constants import RESOURCES, DEPTH_BACKGROUND


class Background(Entity):
    def __init__(self, world, *args, **kwargs):
        if "renderer" not in kwargs:
            raise ValueError("you have to provide a renderer= argument")
        renderer = kwargs['renderer']
        factory = SpriteFactory(TEXTURE, renderer=renderer)
        sprite = factory.from_image(
            RESOURCES.get_path('background.png'))
        sprite.depth = DEPTH_BACKGROUND
        self.sprite = sprite
