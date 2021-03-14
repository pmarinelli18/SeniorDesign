import socket

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_address = ("localhost", 80)
s.connect(server_address)
s.send("hardwareConnection".encode())

while 1:
    data = s.recv(1024)
    if not data:
        break
    print(data)
s.close()