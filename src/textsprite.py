
import sdl2.ext as sdl2ext
from sdl2 import render, surface
from sdl2.sdlttf import (TTF_OpenFont,
                         TTF_RenderText_Shaded,
                         TTF_GetError
                         )
from constants import RESOURCES, WHITE, BLACK


class TextSprite(sdl2ext.TextureSprite):
    def __init__(self, renderer, text="", fontFile="Oswald-Regular.ttf", fontSize=16,
                 textColor=WHITE,
                 backgroundColor=BLACK):
        if isinstance(renderer, sdl2ext.Renderer):
            self.renderer = renderer.renderer
        elif isinstance(renderer, render.SDL_Renderer):
            self.renderer = renderer
        else:
            raise TypeError("unsupported renderer type")

        self.font = TTF_OpenFont(
            bytes(RESOURCES.get_path(fontFile), 'utf-8'),
            fontSize)
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
        position = self.position
        super(TextSprite, self).__init__(texture)
        self.x = position[0]
        self.y = position[1]
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
