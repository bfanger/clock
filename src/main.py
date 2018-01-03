#!/usr/bin/env python3

from ctypes import py_object, pointer, cast, c_void_p, POINTER, byref
import sdl2.ext as sdl2ext
from sdl2 import events as sdlevents, timer
from sdl2.sdlttf import TTF_Init, TTF_Quit
from datetime import datetime
from window import Window, SCREEN_WIDTH, SCREEN_HEIGHT
from textsprite import TextSprite


def get_time():
    return datetime.now().strftime("%H:%M:%S")


class Klok(sdl2ext.Entity):
    def __init__(self, world, *args, **kwargs):
        if "renderer" not in kwargs:
            raise ValueError("you have to provide a renderer= argument")
        renderer = kwargs['renderer']
        # super(Klok, self).__init__(world, *args, **kwargs)
        textSprite = TextSprite(renderer, get_time(), fontSize=92)
        self.textSprite = textSprite
        textSprite.x = SCREEN_WIDTH // 2 - textSprite.size[0] // 2 - 1
        textSprite.y = SCREEN_HEIGHT // 2 - textSprite.size[1] // 2 - 1
        # print(textSprite.size[1])

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

    window = Window()
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
