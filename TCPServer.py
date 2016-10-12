import socketserver
from time import sleep

class MyTCPHandler(socketserver.BaseRequestHandler):

    def handle(self):
        print("Got connection from: {}".format(self.client_address[0]))
        sleep(5)

if __name__ == "__main__":
    HOST, PORT = "localhost", 8888
    server = socketserver.TCPServer((HOST, PORT), MyTCPHandler)
    server.serve_forever()