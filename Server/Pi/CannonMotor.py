import sys
import threading
#import numpy as np # to work with numerical data efficiently

import RPi.GPIO as GPIO 
import time

class CannonMotor():

	global stop
	def __init__(self, stepPin):
		super().__init__()
		GPIO.setmode(GPIO.BOARD)
		GPIO.setup(stepPin, GPIO.OUT)
		self.pwm = GPIO.PWM(stepPin, 50)		#PWM is set to 50 HZ- from spec
		self.pwm.start(0)
		self.setDirection(10)
		self.pwm.ChangeDutyCycle(0)

	def setDirection(self, direction):
		a=10
		b=2
		duty = a / 180 * direction + b
		self.pwm.ChangeDutyCycle(duty)

	def spinMotor(self, stop):
		#Cannon recoils three times
		
		#self.pwm = GPIO.PWM(5, 50)
		#self.pwm.start(0)
		for direction in range(10, 50, 1):
			self.setDirection(direction)
			time.sleep(0.005) # Change sleep time to change speed
		for direction in range(50, 10, -1):
			self.setDirection(direction)
			time.sleep(0.01) # Change sleep time to change speed

		for direction in range(10, 50, 1):
			self.setDirection(direction)
			time.sleep(0.005) # Change sleep time to change speed
		for direction in range(50, 10, -1):
			self.setDirection(direction)
			time.sleep(0.01) # Change sleep time to change speed

		for direction in range(10, 50, 1):
			self.setDirection(direction)
			time.sleep(0.005) # Change sleep time to change speed
		for direction in range(50, 10, -1):
			self.setDirection(direction)
			time.sleep(0.01) # Change sleep time to change speed

		self.pwm.ChangeDutyCycle(0)
		print("Done!")

	
	def startSpinningMotor(self):
		self.stop_threads = False
		self.t1 = threading.Thread(target = self.spinMotor, args =(lambda : self.stop_threads,)) 
		self.t1.start() 
  
	def stopMotorJog(self):
		print("Can't do this for cannon, it stops by itself")

		
