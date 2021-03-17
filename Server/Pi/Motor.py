import sys
import threading
#import numpy as np # to work with numerical data efficiently

#import RPi.GPIO as GPIO 
from time import sleep

class Motor():

	global stop
	def __init__(self, stepPin):
		super().__init__()
		self.stepPin = stepPin
		
	def spinMotor(self, stop):

		while True:
			print("Spin")
			sleep(5);

			if stop():
				print("Done!")
				break

	def startSpinningMotor(self):

		self.stop_threads = False
		self.t1 = threading.Thread(target = self.spinMotor, args =(lambda : self.stop_threads,)) 
		self.t1.start() 
  
	def stopMotorJog(self):
		self.stop_threads = True
		print("Stopping spin")
		self.t1.join() 

		





