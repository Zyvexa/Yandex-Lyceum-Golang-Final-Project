package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TimeOperation struct {
	TimeAddition       int `json:"timeAddition"`
	TimeSubtraction    int `json:"timeSubtraction"`
	TimeMultiplication int `json:"timeMultiplication"`
	TimeDivision       int `json:"timeDivision"`
	TimeExponentiation int `json:"timeExponentiation"`
}

type Expression struct {
	Expression string `json:"expression"`
}

type Data struct {
	Num1               float64 `json:"num1"`
	Num2               float64 `json:"num2"`
	Operation          string  `json:"operation"`
	MaxCalculations    int     `json:"maxCalculations"`
	TimeAddition       int     `json:"timeAddition"`
	TimeSubtraction    int     `json:"timeSubtraction"`
	TimeMultiplication int     `json:"timeMultiplication"`
	TimeDivision       int     `json:"timeDivision"`
	TimeExponentiation int     `json:"timeExponentiation"`
	Index              int     `json:"index"`
}

type Server struct {
	Port     int
	LastTime time.Time
	Free     int
	Total    int
}

type Record struct {
	Expression string
	ID         int
	TimeIn     string
	TimeOut    string
	Answer     string
	Error      string
}

func isValidInfixExpression(expression string) bool {
	// Проверка на пустую строку
	if len(expression) == 0 {
		return false
	}

	// Проверка на допустимые символы
	validChars := "+-()*/^ 1234567890"
	for _, char := range expression {
		if !strings.ContainsRune(validChars, char) {

			return false
		}
	}

	// Проверка на пробелы между операциями и операндами
	for _, char := range "+-*/^" {
		if strings.ContainsRune(expression, char) {
			index := strings.IndexRune(expression, char)
			if expression[index+1] != ' ' || expression[index-1] != ' ' {
				return false
			}
		}
	}

	if strings.Contains(expression, "  ") {
		return false
	}

	// Проверка на правильность вложенности скобок
	openBrackets := strings.Count(expression, "(")
	closeBrackets := strings.Count(expression, ")")
	if openBrackets != closeBrackets {
		return false
	}

	// Проверяем, что между цифрами и скобками есть пробелы
	for _, char := range "()" {
		if strings.ContainsRune(expression, char) {
			index := strings.IndexRune(expression, char)
			if expression[index+1] == ' ' && index == 0 {
			} else if expression[index-1] == ' ' && index == len(expression)-1 {
			} else if expression[index+1] == ' ' && expression[index-1] == ' ' {
			} else {
				return false
			}
		}
	}

	if openBrackets != 0 {
		pattern := `(\()\s(\d+)`
		match, _ := regexp.MatchString(pattern, expression)

		pattern2 := `(\d+)\s(\))`
		match2, _ := regexp.MatchString(pattern2, expression)
		return match && match2
	}
	return true
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/" || token == "^"
}

func isLeftAssociative(token string) bool {
	leftAssociative := map[string]bool{
		"+": true,
		"-": true,
		"*": true,
		"/": true,
		"^": false,
	}
	return leftAssociative[token]
}

func precedence(token string) int {
	operators := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
	}
	return operators[token]
}

func writeToCSV(index_loc int, num1, num2 float64, operation string, times [5]int) error {
	data := Data{
		Num1:               num1,
		Num2:               num2,
		Operation:          operation,
		TimeAddition:       times[0],
		TimeSubtraction:    times[1],
		TimeMultiplication: times[2],
		TimeDivision:       times[3],
		TimeExponentiation: times[4],
		Index:              index_loc,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	file, err := os.OpenFile(global_path_linglist_csv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{fmt.Sprint(index_loc), string(jsonData)})
	if err != nil {
		return fmt.Errorf("failed to write to CSV: %w", err)
	}

	return nil
}

func infixToPostfix(infix string) []string {
	tokens := strings.Split(infix, " ")
	var result []string
	var stack []string

	for _, token := range tokens {
		if token == "" {
			continue
		}
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			result = append(result, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				op := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				result = append(result, op)
			}
			stack = stack[:len(stack)-1] // pop "("
		} else if isOperator(token) {
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) && ((isLeftAssociative(token) && precedence(token) <= precedence(stack[len(stack)-1])) || (!isLeftAssociative(token) && precedence(token) < precedence(stack[len(stack)-1]))) {
				op := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				result = append(result, op)
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		result = append(result, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return result
}

type Stack []float64

func (s *Stack) Push(value float64) {
	(*s) = append((*s), value)
}

func (s *Stack) Pop() float64 {
	if len(*s) == 0 {
		return -1
	}
	top := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return top
}

func (s *Stack) Top() float64 {
	if len(*s) == 0 {
		return -1
	}
	return (*s)[len(*s)-1]
}

func (c *Data) SetIndexToMinusOne() {
	c.Index = -1
}

func ReadCalculations(filename string) ([]Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Пропускаем первую строку
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}

	var calculations []Data
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Разделяем строку на индекс и JSON-строку
		jsonStr := strings.Trim(record[1], `"`)

		// Разбираем JSON-строку в структуру
		var calc Data
		err = json.Unmarshal([]byte(jsonStr), &calc)
		if err != nil {
			return nil, err
		}

		calculations = append(calculations, calc)
	}

	return calculations, nil
}

// ReadServerAddresses считывает адреса серверов из указанного CSV файла.
func ReadServerAddresses(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Пропускаем первую строку
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}

	var serverAddresses []string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Разделяем строку на части
		port, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		free, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
		total, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		// Проверяем условие и добавляем адрес сервера в список
		if free < total {
			serverAddresses = append(serverAddresses, fmt.Sprintf("http://localhost:%d/", port))
		}
	}

	return serverAddresses, nil
}

func SendJSONToServers(serverAddresses []string, calculations []Data) error {
	for _, calc := range calculations {
		if calc.Index == -1 {
			continue
		}
		for _, serverAddress := range serverAddresses {
			// Пропускаем расчеты с индексом -1

			jsonData, err := json.Marshal(calc)
			if err != nil {
				return err
			}

			resp, err := http.Post(serverAddress, "application/json", strings.NewReader(string(jsonData)))
			if err != nil {
				return err
			}
			resp.Body.Close()

			fmt.Printf("Отправлено JSON на %s: %s\n", serverAddress, string(jsonData))
			break
		}
	}

	return nil
}

// WriteCalculations записывает расчеты в указанный CSV файл.
func WriteCalculations(filename string, calculations []Data) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	err = writer.Write([]string{"index", "json"})
	if err != nil {
		return err
	}

	// Записываем расчеты
	for _, calc := range calculations {
		jsonData, err := json.Marshal(calc)
		if err != nil {
			return err
		}
		err = writer.Write([]string{strconv.Itoa(calc.Index), string(jsonData)})
		if err != nil {
			return err
		}
	}

	return nil
}

func checkAndReplaceIndex(index int, filename string) (float64, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV: %w", err)
	}

	for i, record := range records {
		if len(record) < 2 {
			continue
		}

		if record[0] == strconv.Itoa(index) {
			answer, err := strconv.ParseFloat(record[1], 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse float: %w", err)
			}

			// Заменяем индекс на -1
			records[i][0] = "-1"

			// Перезаписываем файл с обновленными данными
			file.Seek(0, 0)  // Перемещаемся в начало файла
			file.Truncate(0) // Очищаем файл
			writer := csv.NewWriter(file)
			if err := writer.WriteAll(records); err != nil {
				return 0, fmt.Errorf("failed to write CSV: %w", err)
			}
			writer.Flush()

			return answer, nil
		}
	}

	return 0, fmt.Errorf("index not found")
}

func evaluationPostfix(postfix []string) float64 {
	var intStack Stack
	for _, char := range postfix {
		opchar := char

		if opchar >= "0" && opchar <= "9" {
			i1, _ := strconv.Atoi(opchar)
			intStack.Push(float64(i1))
		} else {
			opr1 := intStack.Top()
			intStack.Pop()
			opr2 := intStack.Top()
			intStack.Pop()
			var answer float64
			switch char {
			case "^":
				index += 1
				idx_now := index
				err := writeToCSV(idx_now, float64(opr2), float64(opr1), "^", [5]int{global_TimeAddition, global_TimeSubtraction, global_TimeMultiplication, global_TimeDivision, global_TimeExponentiation})
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Data written successfully")
				}

				for {
					ansMutex.Lock()
					answer, err = checkAndReplaceIndex(idx_now, global_path_answer_csv)
					if err == nil {
						ansMutex.Unlock()
						break
					}
					ansMutex.Unlock()
					time.Sleep(time.Second)
				}
				// x := math.Pow(float64(opr2), float64(opr1))
				intStack.Push(float64(int(answer)))
			case "+":
				index += 1
				idx_now := index
				err := writeToCSV(index, float64(opr2), float64(opr1), "+", [5]int{global_TimeAddition, global_TimeSubtraction, global_TimeMultiplication, global_TimeDivision, global_TimeExponentiation})
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Data written successfully")
				}

				for true {
					ansMutex.Lock()
					answer, err = checkAndReplaceIndex(idx_now, global_path_answer_csv)
					if err == nil {
						ansMutex.Unlock()
						break
					}
					ansMutex.Unlock()
					time.Sleep(time.Second)
				}
				fmt.Println(answer)
				intStack.Push(answer)
			case "-":
				index += 1
				idx_now := index
				err := writeToCSV(index, float64(opr2), float64(opr1), "-", [5]int{global_TimeAddition, global_TimeSubtraction, global_TimeMultiplication, global_TimeDivision, global_TimeExponentiation})
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Data written successfully")
				}

				for {
					ansMutex.Lock()
					answer, err = checkAndReplaceIndex(idx_now, global_path_answer_csv)
					if err == nil {
						ansMutex.Unlock()
						break
					}
					ansMutex.Unlock()
					time.Sleep(time.Second)
				}

				intStack.Push(answer)
			case "*":
				index += 1
				idx_now := index
				err := writeToCSV(index, float64(opr2), float64(opr1), "*", [5]int{global_TimeAddition, global_TimeSubtraction, global_TimeMultiplication, global_TimeDivision, global_TimeExponentiation})
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Data written successfully")
				}

				for {
					ansMutex.Lock()
					answer, err = checkAndReplaceIndex(idx_now, global_path_answer_csv)
					if err == nil {
						ansMutex.Unlock()
						break
					}
					ansMutex.Unlock()
					time.Sleep(time.Second)
				}

				intStack.Push(answer)
			case "/":
				index += 1
				idx_now := index
				err := writeToCSV(index, float64(opr2), float64(opr1), "/", [5]int{global_TimeAddition, global_TimeSubtraction, global_TimeMultiplication, global_TimeDivision, global_TimeExponentiation})
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Data written successfully")
				}

				for {
					ansMutex.Lock()
					answer, err = checkAndReplaceIndex(idx_now, global_path_answer_csv)
					if err == nil {
						ansMutex.Unlock()
						break
					}
					ansMutex.Unlock()
					time.Sleep(time.Second)
				}

				intStack.Push(answer)
			}
		}
	}
	return intStack.Top()
}

func main_work(infixExpression string, id_loc int) {
	record := Record{
		Expression: infixExpression,
		ID:         id_loc,
		TimeIn:     time.Now().Format("2006-01-02  15:04:05"),
		TimeOut:    "",
		Answer:     "",
		Error:      "",
	}
	fileMutex.Lock()
	// defer
	if err := writeOrUpdateRecord(record); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Record written or updated successfully")
	}

	lol, err := getDataByID(id_loc, global_path_expression_csv)
	if err != nil {
		record := Record{
			Expression: infixExpression,
			ID:         id_loc,
			TimeIn:     time.Now().Format("2006-01-02  15:04:05"),
			TimeOut:    "",
			Answer:     "",
			Error:      "500",
		}
		if err := writeOrUpdateRecord(record); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Record written or updated successfully")
		}

		return
	}

	if isValidInfixExpression(infixExpression) { // проверка корректности
		fmt.Println("OK")
		lol.Error = fmt.Sprint(200)
		if err := writeOrUpdateRecord(*lol); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Record written or updated successfully")
		}

		fileMutex.Unlock()

		postfixExpression := infixToPostfix(infixExpression) // постфиксная форма

		result := evaluationPostfix(postfixExpression)

		// defer
		fmt.Println("Result:", result) // Expected output: 10
		fileMutex.Lock()
		lol, _ := getDataByID(id_loc, global_path_expression_csv)
		lol.Answer = fmt.Sprint(result)
		lol.TimeOut = time.Now().Format("2006-01-02  15:04:05")
		if err := writeOrUpdateRecord(*lol); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Record written or updated successfully")
		}
		fileMutex.Unlock()
	} else {
		fmt.Println("NOT OK")
		lol.Error = fmt.Sprint(400)
		if err := writeOrUpdateRecord(*lol); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Record written or updated successfully")
		}
		fileMutex.Unlock()
	}
}

func CheckAgent() {
	// Читаем файл agent.csv
	file, err := os.Open(global_path_agent_csv)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Пропускаем первую строку с заголовками
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		fmt.Println("Error reading header:", err)
		return
	}

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Парсим данные и создаем список серверов
	servers := make([]Server, 0)
	for _, record := range records {
		port, _ := strconv.Atoi(record[0])
		lastTime, _ := time.Parse("2006-01-02 15:04:05", record[1])
		free, _ := strconv.Atoi(record[2])
		total, _ := strconv.Atoi(record[3])
		servers = append(servers, Server{Port: port, LastTime: lastTime, Free: free, Total: total})
	}

	// Проверяем серверы и удаляем неактивные
	for i := len(servers) - 1; i >= 0; i-- {
		server := servers[i]
		// Проверяем, что время последней активности сервера было более  5 минут назад
		if time.Since(server.LastTime).Minutes()+float64(global_min_to_add) > 5 {
			// Выполняем GET-запрос
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d", server.Port))
			if err != nil || resp.StatusCode != http.StatusOK {
				// Если запрос не удался, удаляем сервер
				servers = append(servers[:i], servers[i+1:]...)
			}
			if resp != nil {
				resp.Body.Close()
			}
		} else if time.Since(server.LastTime).Minutes()+float64(global_min_to_add) > 1 {
			fmt.Println(server.Port)
			_, _ = http.Get(fmt.Sprintf("http://localhost:%d", server.Port))
			return
		}
	}

	// Сохраняем обновленный список серверов обратно в файл
	file, err = os.Create(global_path_agent_csv)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	// Записываем заголовки в файл
	header := []string{"port", "last_time", "free", "total"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing header to file:", err)
		return
	}

	for _, server := range servers {
		record := []string{
			strconv.Itoa(server.Port),
			server.LastTime.Format("2006-01-02 15:04:05"),
			strconv.Itoa(server.Free),
			strconv.Itoa(server.Total),
		}
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record to file:", err)
			return
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
}

func SendToAgent() {
	calculations, err := ReadCalculations(global_path_linglist_csv)
	if err != nil {
		fmt.Println("Ошибка чтения расчетов:", err)
		return
	}

	serverAddresses, err := ReadServerAddresses(global_path_agent_csv)
	if err != nil {
		fmt.Println("Ошибка чтения адресов серверов:", err)
		return
	}

	if len(serverAddresses) == 0 {
		if len(calculations) != 0 {
			fmt.Println("Файл agent.csv пустой, нет адресов серверов для отправки")
		}
		return
	}

	err = SendJSONToServers(serverAddresses, calculations)
	if err != nil {
		fmt.Println("Ошибка отправки JSON на серверы:", err)
		return
	}

	// Обновляем индекс каждого расчета до -1
	for i := range calculations {
		// calculations[i].SetIndexToMinusOne()
		calculations = RemoveCalculationAtIndex(calculations, i)
	}
	// fmt.Println(calculations)
	if calculations != nil {
		err = WriteCalculations(global_path_linglist_csv, calculations)
		if err != nil {
			fmt.Println("Ошибка записи расчетов обратно в CSV:", err)
		}
	}
}

func RemoveCalculationAtIndex(calculations []Data, index int) []Data {
	if index < 0 || index >= len(calculations) {
		// Если индекс невалиден, возвращаем исходный срез
		return calculations
	}
	// Если индекс валиден, удаляем элемент и возвращаем обновленный срез
	return append(calculations[:index], calculations[index+1:]...)
}

func writeOrUpdateRecord(record Record) error {
	// Открываем файл для чтения и записи
	file, err := os.OpenFile(global_path_expression_csv, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Читаем существующие записи
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read records: %w", err)
	}

	// Проверяем, есть ли уже запись с таким же ID
	found := false
	for i, row := range records {
		if len(row) > 1 && row[1] == fmt.Sprint(record.ID) {
			// Заменяем существующую запись
			records[i] = []string{
				record.Expression,
				fmt.Sprint(record.ID),
				record.TimeIn,
				record.TimeOut,
				record.Answer,
				record.Error,
			}
			found = true
			break
		}
	}

	if !found {
		// Добавляем новую запись
		records = append(records, []string{
			record.Expression,
			fmt.Sprint(record.ID),
			record.TimeIn,
			record.TimeOut,
			record.Answer,
			record.Error,
		})
	}

	// Очищаем файл перед записью обновленных записей
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	// Создаем новый writer для записи обновленных записей
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write records: %w", err)
	}

	// Очищаем буфер writer'а и закрываем файл
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func getDataByID(id int, filePath string) (*Record, error) {
	// Открываем файл на чтение
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Создаем новый Reader для чтения CSV-файла
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Ищем строку с указанным ID
	for _, record := range records {
		if len(record) > 1 && record[1] == strconv.Itoa(id) {
			// Если строка с таким же ID найдена, преобразуем данные в структуру Data
			timeIn, _ := time.Parse("2006-01-02  15:04:05", record[2])
			timeOut, _ := time.Parse("2006-01-02  15:04:05", record[3])
			answer, _ := strconv.ParseFloat(record[4], 64)
			data := &Record{
				Expression: record[0],
				ID:         id,
				TimeIn:     timeIn.Format("2006-01-02  15:04:05"),
				TimeOut:    timeOut.Format("2006-01-02  15:04:05"),
				Answer:     fmt.Sprint(answer),
				Error:      record[5],
			}
			return data, nil
		}
	}

	// Если строка с указанным ID не найдена, возвращаем ошибку
	return nil, fmt.Errorf("no data found for ID %d", id)
}

func WriteSingleRowToCSV(filename, row string) error {
	// Открываем файл с очисткой содержимого
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Записываем строку в файл
	if _, err := file.WriteString(row + "\n"); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

var index int

var fileMutex sync.Mutex
var ansMutex sync.Mutex

var global_TimeAddition int
var global_TimeSubtraction int
var global_TimeMultiplication int
var global_TimeDivision int
var global_TimeExponentiation int

var global_id int

var global_path_agent_csv string
var global_path_answer_csv string
var global_path_expression_csv string
var global_path_linglist_csv string
var global_min_to_add int

func main() {
	global_path_agent_csv = "agent.csv"
	global_path_answer_csv = "answer.csv"
	global_path_expression_csv = "expression.csv"
	global_path_linglist_csv = "long_list.csv"

	global_min_to_add = 180

	index = 0
	global_id = 0
	global_TimeAddition = 1
	global_TimeSubtraction = 1
	global_TimeMultiplication = 1
	global_TimeDivision = 1
	global_TimeExponentiation = 1

	WriteSingleRowToCSV(global_path_answer_csv, "index,answer")
	WriteSingleRowToCSV(global_path_expression_csv, "expression,id,time in,time out,answer,error")
	WriteSingleRowToCSV(global_path_linglist_csv, "index,json")

	go func() {
		for {
			time.Sleep(time.Millisecond * 250)
			SendToAgent()
			CheckAgent()
		}
	}()

	http.HandleFunc("/time", postHandler1)
	http.HandleFunc("/expression", postHandler2)
	http.HandleFunc("/agents", getHandler1)
	http.HandleFunc("/expressions", getHandler2)

	http.ListenAndServe(":8080", nil)
}

func postHandler1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var timeOp TimeOperation
	err := json.NewDecoder(r.Body).Decode(&timeOp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Пример обработки данных из JSON
	fmt.Printf("OK")
	global_TimeAddition = timeOp.TimeAddition
	global_TimeDivision = timeOp.TimeDivision
	global_TimeExponentiation = timeOp.TimeExponentiation
	global_TimeMultiplication = timeOp.TimeMultiplication
	global_TimeSubtraction = timeOp.TimeSubtraction
	// Здесь можно добавить  логику обработки данных, например, выполнение математических операций

	w.WriteHeader(http.StatusOK)
}

func postHandler2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var expr Expression
	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Пример обработки данных из JSON
	fmt.Printf("Received expression: %s\n", expr.Expression)
	global_id += 1
	go main_work(expr.Expression, global_id)

	w.WriteHeader(http.StatusOK)
}

func getHandler1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из CSV файла agent.csv
	file, err := os.Open(global_path_agent_csv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка данных в ответе
	for _, record := range records {
		fmt.Fprintln(w, strings.Join(record, ","))
	}
}

func getHandler2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из CSV файла expression.csv
	file, err := os.Open(global_path_expression_csv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка данных в ответе
	for _, record := range records {
		fmt.Fprintln(w, strings.Join(record, ","))
	}
}
