import socket
import threading
import logging
from time import sleep

# Most of it from http://stackoverflow.com/questions/23828264/how-to-make-a-simple-multithreaded-socket-server-in-python-that-remembers-client
class ThreadedServer(object):
    max_queue = 100
    client_polling_seconds = 5

    def _configure_socket(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)

    def __init__(self, host, port):
        self.host = host
        self.port = port
        self._configure_socket()

    def listen(self):
        self.sock.bind((self.host, self.port))
        self.sock.listen(self.max_queue)
        while True:
            # try:
                client, address = self.sock.accept()
                client.settimeout(self.client_polling_seconds)
                threading.Thread(target = self.listenToClient,args = (client,address)).start()

    def listenToClient(self, client, address):
        logging.info("Client connection from : {}".format(address))
        while True:
            try:
                data = client.recv(1)
                logging.info("Connection from : {} sending information".format(address))
            except socket.timeout as timeout:
                logging.info("Connection from : {} did not send any data, but alive".format(address))
            except Exception as error:
                logging.info("Connection from : {} was closed. Exception: {}".format(address, error))
                client.close()
                return False

def configure_logger():
    logging.basicConfig(
        level=logging.DEBUG,
        format='%(asctime)s %(message)s'
    )
# main()
if __name__ == "__main__":
    port_num = 8888
    host = 'localhost'
    configure_logger()
    logging.debug("Preparing main thread")
    server = ThreadedServer(host,port_num)
    threading.Thread(target = server.listen,daemon=True).start()
    logging.info('Daemon is listening on {}:{}... '.format(host, port_num))
    try:
        while True:
            logging.debug('Listening for connections. Hit ctrl+c to finish')
            sleep(10)
    except KeyboardInterrupt:
        logging.info("Ctrl+c detected")
    logging.info("Finishing our program")