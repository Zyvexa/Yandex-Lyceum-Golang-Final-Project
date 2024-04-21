package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/Zyvexa/grps_test/proto"
	"google.golang.org/grpc"
)

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

// Функция для поиска свободного порта начиная с указанного порта
func findFreePort(startPort int) (int, error) {
	for port := startPort; port <= 65535; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return -1, fmt.Errorf("unable to find a free port starting from %d", startPort)
}

func updateDataFile(portt int) {
	filepath := global_path_agent_csv
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	currentTime := time.Now().Local().Format("2006-01-02 15:04:05 -0700 MST")
	found := false
	for idx, record := range records {
		if len(record) > 0 && idx != 0 {
			recordPort, _ := strconv.Atoi(strings.TrimSpace(record[0]))
			if recordPort == portt {
				strPort := fmt.Sprintf("%d", portt)
				records[idx] = []string{strPort, currentTime, strconv.Itoa(free), strconv.Itoa(total)}
				found = true
				break
			}
		}
	}

	if !found {
		strPort := fmt.Sprintf("%d", portt)
		newRecord := []string{strPort, currentTime, strconv.Itoa(free), strconv.Itoa(total)}
		records = append(records, newRecord)
	}

	file.Seek(0, 0)
	w := csv.NewWriter(file)
	err = w.WriteAll(records)
	if err != nil {
		log.Fatalf("error writing to file: %v", err)
	}
	w.Flush()
}

func writeToCSV(index int, answer float64) error {
	ansMutex.Lock()
	defer ansMutex.Unlock()
	// Открываем файл для записи, если он не существует, он будет создан
	file, err := os.OpenFile(global_path_answer_csv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Форматируем строку с индексом и ответом
	line := fmt.Sprintf("%d,%.2f\n", index, answer)

	// Записываем строку в файл
	if _, err := file.WriteString(line); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func calculate(num1, num2 float64, operation string, time_addition, time_subtraction, time_multiplication, time_division, time_exponentiation, index int) {
	var result float64

	switch operation {
	case "^":
		time.Sleep(time.Duration(time_exponentiation) * time.Second)
		result = math.Pow(num1, num2)
	case "+":
		time.Sleep(time.Duration(time_addition) * time.Second)
		result = num1 + num2
	case "-":
		time.Sleep(time.Duration(time_subtraction) * time.Second)
		result = num1 - num2
	case "*":
		time.Sleep(time.Duration(time_multiplication) * time.Second)
		result = num1 * num2
	case "/":
		time.Sleep(time.Duration(time_division) * time.Second)
		if num2 != 0 {
			result = num1 / num2
		} else {
			fmt.Println("Division by zero is not allowed")
			free -= 1
			result = 0
		}
	default:
		fmt.Println("Invalid operation!")
		free -= 1
		result = 0
	}

	writeToCSV(index, float64(result))

	fmt.Printf("Result of %.2f %s %.2f = %.2f\n", num1, operation, num2, result)
	free -= 1

	updateDataFile(port)
}

// //////////////////grps
type Server struct {
	pb.GeometryServiceServer // сервис из сгенерированного пакета
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SendDataToAgent(
	ctx context.Context,
	in *pb.Data,
) (*pb.Empty, error) {

	var data Data
	err := json.Unmarshal([]byte(in.Json), &data)

	if err != nil {
		fmt.Println("ERROR")
		return &pb.Empty{}, err
	}
	fmt.Println("Accepted for processing")
	free += 1
	updateDataFile(port)
	go calculate(data.Num1, data.Num2, data.Operation, data.TimeAddition, data.TimeSubtraction, data.TimeMultiplication, data.TimeDivision, data.TimeExponentiation, data.Index)

	return &pb.Empty{}, nil
}

func (s *Server) UpdateAgent(
	ctx context.Context,
	in *pb.Empty,
) (*pb.Empty, error) {
	updateDataFile(port)
	fmt.Println("Update")
	return &pb.Empty{}, nil
}

var ansMutex sync.Mutex

var port int

var free int
var total int

var global_path_agent_csv string
var global_path_answer_csv string

func main() {

	global_path_agent_csv = "..\\data\\agent.csv"
	global_path_answer_csv = "..\\data\\answer.csv"

	free = 0
	fmt.Print("Введите максимальное количество операций: ")
	_, err := fmt.Scan(&total)
	if err != nil || total <= 0 {
		fmt.Println("Некоректный ввод", err)
		return
	}

	fmt.Println("Максимальное количество - ", total)

	host := "localhost"

	port, err = findFreePort(9000)
	if err != nil {
		log.Fatalf("Error finding free port: %v", err)
	}
	fmt.Printf("Found free port: %d\n", port)

	updateDataFile(port)

	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части GeometryService
	geomServiceServer := NewServer()
	// зарегистрируем нашу реализацию сервера
	pb.RegisterGeometryServiceServer(grpcServer, geomServiceServer)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
