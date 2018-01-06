from os.path import dirname
from os import environ
from sdl2.ext import Resources
from sdl2.pixels import SDL_Color
from dotenv import load_dotenv

dotenv_path = dirname(dirname(__file__)) + '/.env'
load_dotenv(dotenv_path)

SCREEN_WIDTH = 320
SCREEN_HEIGHT = 240
RESOURCES = Resources(__file__, "../resources")
GREEN = SDL_Color(52, 252, 3, 1)
ORANGE = SDL_Color(152, 119, 1)  # 223, 174, 1
WHITE = SDL_Color(248, 253, 219)  # 255, 255, 255
BLACK = SDL_Color(0, 0, 0)
NS_USERNAME = environ.get('NS_USERNAME')
NS_PASSWORD = environ.get('NS_PASSWORD')
