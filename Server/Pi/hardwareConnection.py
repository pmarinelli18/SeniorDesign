import socket
import json
#import RPi.GPIO as GPIO
import time

#from RadarMotor import *
#from CannonMotor import *
#from DotMatrix import *


import threading
#GPIO.setwarnings(False)
#1 is the pin it is connected to
#p1RadarMotor = RadarMotor(5)
#p2RadarMotor = RadarMotor(??)

#p1CannonMotor = CannonMotor(??)
#p2CannonMotor = CannonMotor(??)

#p1DotMatrix = DotMatrix(??)
#p2DotMatrix = DotMatrix(??)




s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_address = ("localhost", 80)
s.connect(server_address)
s.send("hardwareConnection".encode())

while 1:
	data = s.recv(1024)
	if not data:
		break

	dataToParse = json.loads(data)
	if (dataToParse["id"] == "FinishedMiniGame"):
		p1Temp = json.dumps(dataToParse["player1"])
		p2Temp = json.dumps(dataToParse["player2"])

		player1 = json.loads(p1Temp)
		player2 = json.loads(p2Temp)
		
		print("\n")

		if (player1["radarState"] == "Enabled"):
			print("Player 1 Start spinning motor")
			#p1RadarMotor.startSpinningMotor()
		else:
			print("Player 1 Stop spinning")
			#p1RadarMotor.stopMotorJog()
		if (player2["radarState"] == "Enabled"):
			print("Player 2 Start spinning motor")
			#p2RadarMotor.startSpinningMotor()
		else:
			print("Player 2 Stop spinning")
			#p2RadarMotor.stopMotorJog()

		if (player1["shotCanon"] == True):
			print("Player 1 shot canon")
			#p1CannonMotor.startSpinningMotor()

		if (player2["shotCanon"] == True):
			print("Player 2 shot canon")
			#p2CannonMotor.startSpinningMotor()



	if (dataToParse["id"] == "displayUserNames"):
		print("Player 1 name: " + dataToParse["player1"])
		print("Player 2 name: " + dataToParse["player2"])
		#p1DotMatrix.displayUserName(dataToParse["player1"])
		#p2DotMatrix.displayUserName(dataToParse["player2"])

s.close()


#effects:
#Damage on LED
#Radar spin or not 
#Cannon shoot back and forth a few times
#Torpedo LEDs

