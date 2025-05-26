# TCP-Chat
Go-based chat server implementation, showcasing concurrent communication between multiple clients. Demonstrates the use of:

*   `net` package for handling TCP connections:  The server utilizes the `net` package in Go's standard library to establish TCP connections with clients. This allows for reliable, two-way communication between the server and connected clients.

*   `bufio` for buffered I/O: The `bufio` package is used for efficient reading and writing of data to and from the network connections. Buffered I/O improves performance by reducing the number of system calls required to read and write data.

*   Goroutines for concurrency:  Concurrency is achieved through the use of goroutines, lightweight, concurrent functions in Go. Each client connection is handled in its own goroutine, allowing the server to handle multiple clients simultaneously without blocking.

*   Channels for message passing: Channels provide a safe and efficient way for goroutines to communicate and exchange data. In this project, channels are used to pass messages between the client handling goroutine and the writing goroutine, ensuring messages are delivered to the clients in the correct order.

*   `sync.Mutex` for protecting shared resources:  A mutex (`sync.Mutex`) is used to protect shared resources, such as the list of connected clients (`clients` slice), from concurrent access. This prevents race conditions and ensures data integrity.

*   Direct message passing for simpler architecture: Instead of using a global message channel, the server directly sends messages to each client's channel. This approach simplifies the architecture, reduces the risk of bottlenecks, and improves overall stability.

This project utilizes goroutines and channels for handling multiple connections and message broadcasting efficiently. It avoids global state by using direct message passing.

**What does this project do?**

This project implements a basic chat server that allows multiple clients to connect and exchange messages in real-time.  When a client connects, it is prompted to enter a nickname. Once the nickname is provided, the client can send messages, which are then broadcasted to all other connected clients, prepended with the sender's nickname.

**How to run this project:**

1.  **Install Go:** Make sure you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

2.  **Clone the repository:** Clone this repository to your local machine using `git clone [repository URL]`.

3.  **Navigate to the project directory:**  Use the `cd` command to navigate to the directory where you cloned the repository.

4.  **Run the server:** Execute the command `go run server.go` in your terminal to start the chat server. The server will listen for incoming connections on port 8000.

5.  **Run the client:** Open one or more additional terminal windows. In each window, navigate to the project directory and execute the command `go run client.go`. Each client will connect to the server.

6.  **Chat:** In each client window, enter your nickname when prompted. After that, you can start typing messages and pressing Enter. Your messages will be displayed in all other connected client windows, prepended with your nickname.
