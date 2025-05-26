package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Client структура, представляющая подключенного клиента.
type Client struct {
	conn net.Conn    // Сетевое соединение с клиентом.
	nick string      // Никнейм клиента.
	ch   chan string // Канал для отправки сообщений этому клиенту.  Буферизованный канал помогает избежать блокировок при отправке сообщений.
}

// clients - слайс, содержащий всех подключенных клиентов.
var clients []*Client

// clientsMutex - мьютекс для защиты слайса `clients` от одновременного доступа (race conditions).
var clientsMutex sync.Mutex

// wg - WaitGroup для ожидания завершения всех горутин перед завершением программы.
var wg sync.WaitGroup

// handleConnection - функция, обрабатывающая новое подключение клиента.
func handleConnection(conn net.Conn) {
	wg.Add(1)       // Увеличиваем счетчик WaitGroup на 1.
	defer wg.Done() // Уменьшаем счетчик WaitGroup на 1 при завершении функции (гарантированно выполнится).

	fmt.Println("New client connected:", conn.RemoteAddr()) // Логируем новое подключение.

	// Создаем нового клиента.
	client := Client{
		conn: conn,
		nick: "",                    // Инициализируем никнейм пустым значением.
		ch:   make(chan string, 10), // Создаем буферизованный канал для сообщений клиента (размер буфера 10).
	}

	reader := bufio.NewReader(conn) // Создаем читатель для чтения данных из соединения.

	// Читаем никнейм клиента из соединения.
	nickname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading nickname:", err) // Логируем ошибку, если не удалось прочитать никнейм.
		conn.Close()                                // Закрываем соединение.
		return                                      // Завершаем функцию.
	}

	nickname = strings.TrimSpace(nickname)         // Удаляем пробелы в начале и конце никнейма.
	client.nick = nickname                         // Устанавливаем никнейм клиента.
	fmt.Println("Received nickname:", client.nick) // Логируем полученный никнейм.

	// Отправляем приветственное сообщение клиенту.
	welcomeMessage := "Welcome to the chat, " + client.nick + "!\n"
	_, err = conn.Write([]byte(welcomeMessage))
	if err != nil {
		fmt.Println("Error sending welcome message:", err) // Логируем ошибку, если не удалось отправить приветственное сообщение.
		conn.Close()                                       // Закрываем соединение.
		return                                             // Завершаем функцию.
	}

	clientsMutex.Lock()                          // Захватываем мьютекс перед изменением слайса `clients`.
	fmt.Println("Adding client to clients list") // Логируем добавление клиента в слайс.
	clients = append(clients, &client)           // Добавляем клиента в слайс `clients`.
	clientsMutex.Unlock()                        // Освобождаем мьютекс после изменения слайса `clients`.

	// Горутина для чтения сообщений от клиента.
	go func() {
		wg.Add(1)       // Увеличиваем счетчик WaitGroup для горутины чтения.
		defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении горутины.
		defer func() {
			// Эта функция выполнится после завершения горутины (даже если произойдет panic).
			clientsMutex.Lock()                              // Захватываем мьютекс перед изменением слайса `clients`.
			fmt.Println("Removing client from clients list") // Логируем удаление клиента из слайса.
			// Удаляем клиента из слайса.
			for i, c := range clients {
				if c == &client {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			clientsMutex.Unlock()                                // Освобождаем мьютекс после изменения слайса `clients`.
			conn.Close()                                         // Закрываем соединение.
			fmt.Println("Connection closed from message reader") // Логируем закрытие соединения.
		}()

		// Бесконечный цикл для чтения сообщений от клиента.
		for {
			messageFromClient, err := reader.ReadString('\n') // Читаем сообщение от клиента.
			if err != nil {
				fmt.Println("Client disconnected:", conn.RemoteAddr(), "Error:", err) // Логируем отключение клиента.
				return                                                                // Завершаем горутину.
			}
			formattedMessage := fmt.Sprintf("%s: %s", client.nick, messageFromClient) // Форматируем сообщение, добавляя никнейм клиента.

			clientsMutex.Lock() // Захватываем мьютекс перед отправкой сообщения другим клиентам.
			// Отправляем сообщение всем остальным клиентам.
			for _, otherClient := range clients {
				if otherClient != &client { // Не отправляем сообщение самому себе.
					select {
					case otherClient.ch <- formattedMessage: // Отправляем сообщение в канал клиента.
						// Сообщение отправлено успешно.
					default:
						fmt.Println("Client channel is full, dropping message") // Логируем, если канал клиента заполнен (сообщение будет потеряно).
					}
				}
			}
			clientsMutex.Unlock() // Освобождаем мьютекс после отправки сообщения другим клиентам.
		}
	}()

	// Горутина для отправки сообщений клиенту.
	go func() {
		wg.Add(1)              // Увеличиваем счетчик WaitGroup для горутины отправки.
		defer wg.Done()        // Уменьшаем счетчик WaitGroup при завершении горутины.
		defer close(client.ch) // Закрываем канал клиента при завершении горутины (сигнал для завершения цикла чтения из канала).

		// Читаем сообщения из канала клиента и отправляем их в соединение.
		for msg := range client.ch {
			_, err := conn.Write([]byte(msg)) // Отправляем сообщение клиенту.
			if err != nil {
				fmt.Println("Error sending message to client:", err) // Логируем ошибку, если не удалось отправить сообщение.
				return                                               // Завершаем горутину.
			}
		}
	}()
}

// main - главная функция программы.
func main() {
	clients = make([]*Client, 0) // Инициализируем слайс `clients`.

	// Создаем listener для прослушивания входящих соединений на порту 8000.
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error creating listener:", err) // Логируем ошибку, если не удалось создать listener.
		return                                       // Завершаем программу.
	}
	defer listener.Close() // Закрываем listener при завершении программы (гарантированно выполнится).

	fmt.Println("Server is listening on :8000") // Логируем запуск сервера.

	// Бесконечный цикл для принятия входящих соединений.
	for {
		conn, err := listener.Accept() // Принимаем новое соединение.
		if err != nil {
			fmt.Println("Error accepting connection:", err) // Логируем ошибку, если не удалось принять соединение.
			continue                                        // Переходим к следующей итерации цикла.
		}
		go handleConnection(conn) // Запускаем горутину для обработки нового соединения.
	}
}
