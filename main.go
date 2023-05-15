package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const version = "2023.2.5.15"
const serviceName = "json api service"
const serviceDescription = "json api service"

const config = "user=postgres password=pj79.. dbname=system host=localhost port=5432 sslmode=disable application_name=xstock_webservice"

const (
	administrator = iota + 1
	poweruser
	user
)

type SalaryPerUser struct {
	FirstName     string
	LastName      string
	Email         string
	Age           int
	MonthlySalary []MonthlySalary
}

type MonthlySalary struct {
	Basic int
	HRA   int
	TA    int
}

type program struct{}

func main() {
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] starting..."))
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] © "+strconv.Itoa(time.Now().Year())+" Jachym Jahodaa"))
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "ERR [SYSTEM] Cannot start: "+err.Error()))
	}

	err = s.Run()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "ERR [SYSTEM] Cannot start: "+err.Error()))
	}
}

func (p *program) Start(service.Service) error {
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] started"))
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] stopped"))
	return nil
}

func (p *program) run() {
	router := httprouter.New()
	router.GET("/allusers", loadAllUsersFromJson)
	router.GET("/user", loadSpecificUserFromJson)

	err := http.ListenAndServe(":100", router)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "ERR [SYSTEM] Problem starting service: "+err.Error()))
		os.Exit(-1)
	}
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] running"))

}

func loadSpecificUserFromJson(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var url = request.URL.String()
	if !strings.Contains(url, "?") {
		fmt.Fprintf(writer, "Bad data request")
		return
	}
	parsedUrl := strings.Split(url, "?")
	if !strings.Contains(parsedUrl[1], "=") {
		fmt.Fprintf(writer, "Bad data request")
		return
	}
	letter := strings.Split(parsedUrl[1], "=")[1]

	jsonfile := readJson()

	for i, data := range jsonfile {
		data.LastName = jsonfile[i].LastName
		if data.LastName[0:1] == letter {
			dataForPage, err := json.Marshal(jsonfile[i])
			if err != nil {
				fmt.Fprintf(writer, err.Error())
				return
			}
			fmt.Fprintf(writer, string(dataForPage))
			return
		}
	}
}

func loadAllUsersFromJson(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	jsonfile := readJson()
	dataForPage, err := json.Marshal(jsonfile)
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}

	fmt.Fprintf(writer, string(dataForPage))
}

func readJson() []SalaryPerUser {
	file, err := os.Open("test.json")
	if err != nil {
		log.Println("Error opening json file:", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading json data:", err)
	}
	var jsonfile []SalaryPerUser
	err = json.Unmarshal(data, &jsonfile)
	if err != nil {
		log.Println("Error unmarshalling json data:", err)
	}
	return jsonfile
}
