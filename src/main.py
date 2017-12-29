#!/usr/bin/env python3

from ctypes import py_object, pointer, cast, c_void_p, POINTER, byref

import sdl2.ext as sdl2ext
from sdl2 import (pixels, render, events as sdlevents, surface,
                  timer)
from sdl2.sdlttf import (TTF_OpenFont,
                         TTF_RenderText_Shaded,
                         TTF_GetError,
                         TTF_Init,
                         TTF_Quit
                         )
from datetime import datetime

SCREEN_WIDTH = 320
SCREEN_HEIGHT = 240
RESOURCES = sdl2ext.Resources(__file__, "../resources")


class TextSprite(sdl2ext.TextureSprite):
    def __init__(self, renderer, text="", fontSize=16,
                 textColor=pixels.SDL_Color(255, 255, 255),
                 backgroundColor=pixels.SDL_Color(0, 0, 0)):
        if isinstance(renderer, sdl2ext.Renderer):
            self.renderer = renderer.renderer
        elif isinstance(renderer, render.SDL_Renderer):
            self.renderer = renderer
        else:
            raise TypeError("unsupported renderer type")

        self.font = TTF_OpenFont(
            bytes(RESOURCES.get_path("tuffy.ttf"), 'utf-8'), fontSize)
        if self.font is None:
            raise TTF_GetError()
        self._text = text
        self.fontSize = fontSize
        self.textColor = textColor
        self.backgroundColor = backgroundColor
        texture = self._createTexture()

        super(TextSprite, self).__init__(texture)

    def _createTexture(self):
        textSurface = TTF_RenderText_Shaded(
            self.font,
            bytes(self._text, 'utf-8'), self.textColor, self.backgroundColor)
        if textSurface is None:
            raise TTF_GetError()
        texture = render.SDL_CreateTextureFromSurface(
            self.renderer, textSurface)
        if texture is None:
            raise sdl2ext.SDLError()
        surface.SDL_FreeSurface(textSurface)
        return texture

    def _updateTexture(self):
        textureToDelete = self.texture

        texture = self._createTexture()
        super(TextSprite, self).__init__(texture)

        render.SDL_DestroyTexture(textureToDelete)

    @property
    def text(self):
        return self._text

    @text.setter
    def text(self, value):
        if self._text == value:
            return
        self._text = value
        self._updateTexture()


def get_time():
    return datetime.now().strftime("%H:%M:%S")


class Klok(sdl2ext.Entity):
    def __init__(self, world, *args, **kwargs):
        if "renderer" not in kwargs:
            raise ValueError("you have to provide a renderer= argument")
        renderer = kwargs['renderer']
        # super(Klok, self).__init__(world, *args, **kwargs)
        textSprite = TextSprite(renderer, get_time(), fontSize=72)
        self.textSprite = textSprite
        print(textSprite.size)
        SCREEN_WIDTH
        object.__setattr__(self, 'callback', self.getCallBackFunc())
        object.__setattr__(self, 'timerId', timer.SDL_AddTimer(
            1000, self.callback, None))

    def getCallBackFunc(self):
        def oneSecondElapsed(time, userdata):
            event = sdlevents.SDL_Event()
            user_event = sdlevents.SDL_UserEvent()

            user_event.type = sdlevents.SDL_USEREVENT
            user_event.code = 0
            user_event.data1 = cast(pointer(py_object(self)), c_void_p)
            user_event.data2 = 0

            event.type = sdlevents.SDL_USEREVENT
            event.user = user_event

            sdlevents.SDL_PushEvent(event)

            return time
        return timer.SDL_TimerCallback(oneSecondElapsed)


def main():
    sdl2ext.init()
    TTF_Init()

    window = sdl2ext.Window("Klok", size=(SCREEN_WIDTH, SCREEN_HEIGHT))
    window.show()

    renderer = sdl2ext.Renderer(window)
    factory = sdl2ext.SpriteFactory(sdl2ext.TEXTURE, renderer=renderer)
    world = sdl2ext.World()

    Klok(world, renderer=renderer)

    spriteRenderer = factory.create_sprite_render_system()
    world.add_system(spriteRenderer)

    def eventloop():

        event = sdlevents.SDL_Event()
        while True:
            ret = sdlevents.SDL_WaitEvent(byref(event))
            if ret == 0:
                raise sdl2ext.SDLError()

            if event.type == sdlevents.SDL_QUIT:
                break
            elif event.type == sdlevents.SDL_USEREVENT:
                entity = cast(event.user.data1, POINTER(
                    py_object)).contents.value
                entity.textsprite.text = get_time()
            else:
                print("Event:", event.type)
            renderer.clear()
            world.process()

    try:
        eventloop()
    except KeyboardInterrupt:
        pass

    TTF_Quit()
    sdl2ext.quit()
    return 0


if __name__ == '__main__':
    main()
