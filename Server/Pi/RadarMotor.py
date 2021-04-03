import sys
import threading
#import numpy as np # to work with numerical data efficiently

import RPi.GPIO as GPIO 
import time

class RadarMotor():

	global stop
	def __init__(self, stepPin):
		super().__init__()
		GPIO.setmode(GPIO.BOARD)
		GPIO.setup(stepPin, GPIO.OUT)
		self.pwm = GPIO.PWM(stepPin, 50)		#PWM is set to 50 HZ- from spec
		self.pwm.start(0)
		self.t1 = None
		self.isSpinning = False
		self.setDirection(10)

	def setDirection(self, direction):
		a=10
		b=2
		duty = a / 180 * direction + b
		self.pwm.ChangeDutyCycle(duty)
		time.sleep(.1) # Change sleep time to change speed

	def spinMotor(self, stop):
		#self.pwm = GPIO.PWM(5, 50)
		#self.pwm.start(0)
		while True:
			for direction in range(10, 170, 1):
				self.setDirection(direction)
			for direction in range(170, 10, -1):
				self.setDirection(direction)
			print("Spin")
			#sleep(5);

			if stop():
				self.pwm.ChangeDutyCycle(0)
				#self.pwm.stop()
				print("Done!")
				break

	def startSpinningMotor(self):
		self.stop_threads = False
		if (!self.isSpinning):
			self.isSpinning = True
			self.t1 = threading.Thread(target = self.spinMotor, args =(lambda : self.stop_threads,)) 
			self.t1.start() 
  
	def stopMotorJog(self):
		if (self.t1 != None):
			self.stop_threads = True
			print("Stopping spin")
			self.t1.join() 
			self.isSpinning = False

		
