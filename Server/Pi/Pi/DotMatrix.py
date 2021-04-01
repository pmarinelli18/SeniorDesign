import RPi.GPIO as GPIO
from time import sleep, strftime
from datetime import datetime
from luma.core.interface.serial import spi, noop
from luma.core.render import canvas
from luma.core.virtual import viewport
from luma.led_matrix.device import max7219
from luma.core.legacy import text, show_message
from luma.core.legacy.font import proportional, CP437_FONT, LCD_FONT

class DotMatrix():

	device = None
	virtual = None
	#portNumber was 0 during testing
	def __init__(self, portNumber):
		global virtual
		super().__init__()
		serial = spi(port=portNumber, device=0, gpio=noop())
		device = max7219(serial, width=32, height=8, block_orientation=-90)
		device.contrast(5)
		virtual = viewport(device, width=32, height=16)


	def displayUserName(self, userName):
		global virtual
		with canvas(virtual) as draw:
            		text(draw, (0, 1), userName, fill="white", font=proportional(CP437_FONT))
