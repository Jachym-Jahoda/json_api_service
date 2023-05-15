package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
	Id              int              `jsonapi:"primary,salary_per_user"`
	FirstName       string           `jsonapi:"attr,first_name"`
	LastName        string           `jsonapi:"attr,last_name"`
	Email           string           `jsonapi:"attr,email"`
	Age             int              `jsonapi:"attr,age"`
	MonthlySalaries []*MonthlySalary `jsonapi:"relation,monthly_salaries"`
}

type MonthlySalary struct {
	Id    int `jsonapi:"primary,monthly_salary"`
	Basic int `jsonapi:"attr,basic"`
	HRA   int `jsonapi:"attr,hra"`
	TA    int `jsonapi:"attr,ta"`
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
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/data/*filepath", http.Dir("json"))

	loadJsonData()
	//router.POST("/load_json_data", loadJsonData)

	err := http.ListenAndServe(":90", router)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "ERR [SYSTEM] Problem starting service: "+err.Error()))
		os.Exit(-1)
	}
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] running"))

}

func loadJsonData() {
	for {
		file, err := os.Open("test.json")
		if err != nil {
			log.Println("Error opening json file:", err)
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Error reading json data:", err)
		}

		var jsonfile []SalaryPerUser
		err = json.Unmarshal(data, &jsonfile)
		if err != nil {
			log.Println("Error unmarshalling json data:", err)
		}
	}
}
