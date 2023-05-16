import socket
import json
import time

# Replace with your server address and port
server_address = "localhost"
server_port = 55555

def send_message(sock, message):
    sock.sendall(message)

def receive_message(sock, buffer_size=1024):
    try:
        data = sock.recv(buffer_size)
    except BlockingIOError:
        print("No data available to read.")
        return None

    return data

def main():
    # Create a TCP/IP socket
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        # Connect to the Go server
        sock.connect((server_address, server_port))

        # Send our city name
        send_message(sock, "PythonCity".encode())

        # Receive the remote city name
        response = receive_message(sock).decode()
        remote_city_name = response.strip()
        print(f"Connected to city {remote_city_name}")
        
        # make sure they add us to their city before sending any other messages
        time.sleep(1.00)

        for i in range(10):
            time.sleep(0.00)
            print("Sending merchant...")

            merchant_data = {
                "Money": 1211.4803370633176,
                "BuysSells": "chair",
                "CarryingCapacity": 20,
                "Owned": 3,
                "ExpectedPrices": {
                    "bed": {
                            "PORTSVILLE": 31.896729121741835,
                            "RIVERWOOD": 16.95381916897244,
                            "SEASIDE": 30.22688570620929,
                            "WINTERHOLD": 32.65932633911733,
                            "PythonCity": 31.0
                    },
                    "chair": {
                            "PORTSVILLE": 11.045778022141988,
                            "RIVERWOOD": 20.342182353944178,
                            "SEASIDE": 10.709200882363625,
                            "WINTERHOLD": 11.756341277296725,
                            "PythonCity": 12.0
                    },
                    "thread": {
                            "PORTSVILLE": 2.421718986326459,
                            "RIVERWOOD": 4.585140611165299,
                            "SEASIDE": 2.301383393344175,
                            "WINTERHOLD": 2.4909610238149353,
                            "PythonCity": 3.0
                    },
                    "wood": {
                            "PORTSVILLE": 2.371161055748651,
                            "RIVERWOOD": 4.541456107859949,
                            "SEASIDE": 2.288907176392255,
                            "WINTERHOLD": 2.4714546798485775,
                            "PythonCity": 3.0
                    }
                }
            }
            print(len(json.dumps(merchant_data).encode()))
            message = json.dumps(merchant_data) + "\n"
            send_message(sock, message.encode())

        # send back any received merchants
        sock.setblocking(False)
        while True:
            message = receive_message(sock)
            if message:
                send_message(sock, message)
                print("sending back merchant")
                time.sleep(0.1)
            else:
                break

        time.sleep(1)
        # Close the connection
        sock.close()

if __name__ == "__main__":
    main()