from sdl2.ext import Window as SDL_Window
from sdl2.ext.common import SDLError
from sdl2 import (SDL_GetNumVideoDisplays,
                  SDL_DisplayMode,
                  SDL_GetCurrentDisplayMode,
                  SDL_WINDOWPOS_CENTERED_DISPLAY,
                  SDL_WINDOW_FULLSCREEN,
                  SDL_ShowCursor,
                  SDL_DISABLE)
from constants import SCREEN_WIDTH, SCREEN_HEIGHT


class Window(SDL_Window):
    def __init__(self):
        displays = SDL_GetNumVideoDisplays()
        flags = 0

        if displays > 1:
            position = (SDL_WINDOWPOS_CENTERED_DISPLAY(1),
                        SDL_WINDOWPOS_CENTERED_DISPLAY(1))
        else:
            position = (0, 0)
            mode = SDL_DisplayMode()
            if SDL_GetCurrentDisplayMode(0, mode) != 0:
                raise SDLError()
            if mode.w == 320:
                flags += SDL_WINDOW_FULLSCREEN

        SDL_ShowCursor(SDL_DISABLE)
        super(Window, self).__init__(
            "Klok",
            flags=flags,
            size=(SCREEN_WIDTH, SCREEN_HEIGHT),
            position=position)
