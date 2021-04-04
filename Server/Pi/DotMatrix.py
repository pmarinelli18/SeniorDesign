import RPi.GPIO as GPIO
import random
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
		super().__init__()
		serial = spi(port=portNumber, device=0, gpio=noop())
		self.device = max7219(serial, width=32, height=8, block_orientation=-90)
		self.device.contrast(5)
		self.virtual = viewport(device, width=32, height=8)
		self.damagedLEDs = []


	def displayUserName(self, userName):
		with canvas(self.virtual) as draw:
			text(draw, (0, 1), userName, fill="white", font=proportional(CP437_FONT))

	def displayDamage(self):
		with canvas(self.virtual) as draw:
			draw.rectangle(virtual.bounding_box, outline="black", fill="black")
			for x in self.damagedLEDs:
				draw.point((x.col, x.row), fill="white")

	def updateHalth(self, health):
		onLights = int(50 - (50 * (health/100)))
		print(onLights)
		if len(self.damagedLEDs) - onLights > 0:
			print("Need to remove LEDs")
			for x in range(len(self.damagedLEDs) - onLights):
				#randomly remove one element at a time
				currentLength = len(self.damagedLEDs)
				self.damagedLEDs.pop(random.randint(0,currentLength-1))
		elif len(self.damagedLEDs) - onLights < 0:
			print("Need to add LEDs")
			for x in range(-1 * (len(self.damagedLEDs) - onLights)):
				self.damagedLEDs.append(BoatDamage())
		else:
			print("No changes!")
		self.displayDamage();


class BoatDamage():
	def __init__(self):
		self.col = random.randint(0,32)
		self.row = random.randint(0,8)
