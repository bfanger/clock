
from ctypes import py_object, cast, POINTER, byref
import sdl2.ext as sdl2ext
from sdl2 import events as sdlevents
from sdl2.keycode import SDLK_ESCAPE


def eventloop(world, renderer):

    event = sdlevents.SDL_Event()
    while True:
        ret = sdlevents.SDL_WaitEvent(byref(event))
        if ret == 0:
            raise sdl2ext.SDLError()

        if event.type == sdlevents.SDL_QUIT:
            # Closed window
            break
        elif (event.type == sdlevents.SDL_KEYUP and
                event.key.keysym.sym == SDLK_ESCAPE):
            # Pressed ESC
            break
        elif event.type == sdlevents.SDL_USEREVENT:
            # Timer event (probably)
            entity = cast(event.user.data1, POINTER(
                py_object)).contents.value
            entity.update()
        else:
            print("Event:", event.type)

        renderer.clear()
        world.process()
