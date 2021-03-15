import socket
import json


s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_address = ("localhost", 80)
s.connect(server_address)
s.send("hardwareConnection".encode())

while 1:
	data = s.recv(1024)
	if not data:
		break

	dataToParse = json.loads(data)
	p1Temp = json.dumps(dataToParse["player1"])
	p2Temp = json.dumps(dataToParse["player2"])

	player1 = json.loads(p1Temp)
	player2 = json.loads(p2Temp)
	
	#Replace the print statements with actual code that uses the GPIO ports to change how the hardware is moving
	print("\n\nPlayer 1: \nHealth: " + player1["health"]+"\nRadar State: " + player1["radarState"])
	print("Player 2: \nHealth: " + player2["health"]+"\nRadar State: " + player2["radarState"])
s.close()

#effects:
#Damage on LED
#Radar spin or not
#Cannon shoot back and forth a few times
#Torpedo LEDs