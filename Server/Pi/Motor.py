import sys
import threading
#import numpy as np # to work with numerical data efficiently

#import RPi.GPIO as GPIO 
from time import sleep

class Motor():

	global stop
	def __init__(self, stepPin):
		super().__init__()
		GPIO.setmode(GPIO.BOARD)
		GPIO.setup(stepPin, GPIO.OUT)
		self.pwm = GPIO.PWM(stepPin, 50)		#PWM is set to 50 HZ- from spec
		self.pwm.start(0)
		
	def spinMotor(self, stop):

		while True:
			for direction in range(0, 181, 10):
				setDirection(direction)
			for direction in range(181, 0, -10):
				setDirection(direction)
			print("Spin")
			#sleep(5);

			if stop():
				print("Done!")
				break

	def setDirection(direction):
		duty = a / 180 * direction + b
		self.pwm.ChangeDutyCycle(duty)
		time.sleep(1) # Change sleep time to change speed
	
	def startSpinningMotor(self):
		self.stop_threads = False
		self.t1 = threading.Thread(target = self.spinMotor, args =(lambda : self.stop_threads,)) 
		self.t1.start() 
  
	def stopMotorJog(self):
		self.stop_threads = True
		print("Stopping spin")
		self.t1.join() 

		





