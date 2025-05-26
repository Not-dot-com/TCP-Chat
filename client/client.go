package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// main - главная функция программы.
func main() {
	// Устанавливаем соединение с сервером по адресу localhost:8000.
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Error connecting to server:", err) // Логируем ошибку, если не удалось установить соединение.
		return                                          // Завершаем программу.
	}
	defer conn.Close() // Закрываем соединение при завершении программы (гарантированно выполнится).

	reader := bufio.NewReader(os.Stdin) // Создаем читатель для чтения данных из стандартного ввода (клавиатуры).

	// Запрашиваем никнейм у пользователя.
	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n') // Читаем никнейм из стандартного ввода.
	nickname = strings.TrimSpace(nickname) // Удаляем пробелы в начале и конце никнейма.

	// Отправляем никнейм на сервер.
	_, err = conn.Write([]byte(nickname + "\n"))
	if err != nil {
		fmt.Println("Error sending nickname:", err) // Логируем ошибку, если не удалось отправить никнейм.
		return                                      // Завершаем программу.
	}
	fmt.Println("Nickname sent") // Логируем успешную отправку никнейма.

	// Горутина для чтения сообщений от сервера.
	go func() {
		serverReader := bufio.NewReader(conn) // Создаем читатель для чтения данных из соединения с сервером.
		// Бесконечный цикл для чтения сообщений от сервера.
		for {
			message, err := serverReader.ReadString('\n') // Читаем сообщение от сервера.
			if err != nil {
				fmt.Println("Server disconnected:", err) // Логируем отключение сервера.
				return                                   // Завершаем горутину.
			}
			fmt.Print(message) // Выводим полученное сообщение в стандартный вывод (консоль).
		}
	}()

	// Основной цикл для отправки сообщений на сервер.
	for {
		fmt.Print("Enter message: ")          // Запрашиваем сообщение у пользователя.
		message, _ := reader.ReadString('\n') // Читаем сообщение из стандартного ввода.
		message = strings.TrimSpace(message)  // Удаляем пробелы в начале и конце сообщения.

		// Отправляем сообщение на сервер.
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err) // Логируем ошибку, если не удалось отправить сообщение.
			return                                     // Завершаем программу.
		}
	}
}
