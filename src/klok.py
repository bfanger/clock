from ctypes import py_object, pointer, cast, c_void_p
from datetime import datetime
import sdl2.ext as sdl2ext
from sdl2 import (timer,
                  SDL_Event,
                  SDL_UserEvent,
                  SDL_USEREVENT,
                  SDL_PushEvent)
from textsprite import TextSprite


def get_time():
    return datetime.now().strftime("%H:%M")


class Klok(sdl2ext.Entity):
    def __init__(self, world, *args, **kwargs):
        if "renderer" not in kwargs:
            raise ValueError("you have to provide a renderer= argument")
        renderer = kwargs['renderer']
        # super(Klok, self).__init__()
        textSprite = TextSprite(renderer, '00:00', fontSize=120)
        self.textSprite = textSprite
        textSprite.x = 120
        textSprite.y = 96
        self.update()

        object.__setattr__(self, 'callback', self.getCallBackFunc())
        object.__setattr__(self, 'timerId', timer.SDL_AddTimer(
            15000, self.callback, None))

    def update(self):
        pass
        self.textsprite.text = get_time()

    def getCallBackFunc(self):
        def oneSecondElapsed(time, userdata):
            event = SDL_Event()
            user_event = SDL_UserEvent()

            user_event.type = SDL_USEREVENT
            user_event.code = 0
            user_event.data1 = cast(pointer(py_object(self)), c_void_p)
            user_event.data2 = 0

            event.type = SDL_USEREVENT
            event.user = user_event

            SDL_PushEvent(event)

            return time
        return timer.SDL_TimerCallback(oneSecondElapsed)
