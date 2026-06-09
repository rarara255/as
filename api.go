// (Router) Multiplexer
// DefaultServeMux(внутри хранится map соответствий между URL <-> Handler)
//Способы регистрации путей:
//1) http.Handle(url, handle)
//2) http.HandleFunc(url,handFunc(w http.ResponseWriter, r *http.Request){})
//Правило строгого совпадения с URL: /api/data != /api/data/
//request -> URi /static/styles.css

// 1.URL "/"
// 2 .URL "/static/"

package main

import (
	"fmt"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	"sync"
	"strconv"
	"log/slog"
)

type Result struct{
	Source string `json:"source"`
	Data string `json:"data"`
}

type CatFact struct {
	Fact string `json:"Fact"`
}

type MainData struct {
	Temp      float64 `json:"temp"`
	FeelsTemp float64 `json:"feels_like"`
}

type Weather struct {
	Name string   `json:"name"`
	Main MainData `json:"main"`
}

func fetch(apiType string, url string, results chan <- Result, wg *sync.WaitGroup){
	defer wg.Done()
	resp, errGet := http.Get(url)
	if errGet != nil{
		results <- Result{Source: apiType, Data: "Не удалось связаться с сервером"}
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var message string

	if apiType == "CAT"{
		var fact CatFact
		json.Unmarshal(body, &fact)
		message = fact.Fact
	}else{
		var weather Weather
		json.Unmarshal(body, &weather)
		message = fmt.Sprintf("Погода: в %s, Температура %.2f, Ощущается как %.2f", weather.Name, weather.Main.Temp, weather.Main.FeelsTemp)
	}

	results <- Result {Source: apiType, Data: message}
}

func apiHandler(w http.ResponseWriter, r *http.Request){
	countSTR :=  r.URL.Query().Get("count")

	count, countErr := strconv.Atoi(countSTR)
	if countErr != nil || count < 1{
		count = 1
	}

	log.Printf("Получен запрос. Запущен сбор данных по %d запросам", count)

	results := make(chan Result, count*2)
	var wg sync.WaitGroup

	for i := 0; i < count; i++{
		wg.Add(2)
		go fetch("CAT", "https://catfact.ninja/fact", results, &wg)
		go fetch("WEATHER", "https://api.openweathermap.org/data/2.5/weather?q=Ulan-Ude&units=metric&appid=fe7a3b2ef70de13b795768ff1f177b24", results, &wg)
	}

	go func() {
		wg.Wait() // находимся в этой точке, пока все горутины не завершат свою работу 
		close(results) // закрытие канала
	}()

	var finalResults []Result

	// данные берутся по мере поступления 
	// for завершится сразу после закрытия канала
	for result := range results {
		finalResults = append(finalResults, result)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(finalResults)
}


func main() {
	mux := http.NewServeMux() // создание локального мультиплексора
	//подключение файлового сервера к роутеру на корневой путь
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fileServer)
	mux.HandleFunc("/api/fetch", apiHandler)

	//используется префиксный путь
	// mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request){
	// 	fmt.Fprintf(w, "Добро пожаловать в API. Вы запросили: %s,  %s", r.URL.Path, r.URL.Hostname())
	// })

	//строгий путь для точного совпадения
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Страница о проекта"))
	})
	//универсальный обработчик(корень)
	//является префиксом для всего
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	// 	if r.URL.Path != "/"{
	// 		http.NotFound(w,r)
	// 		return
	// 	}

	// 	w.Write([]byte("Главная страница"))
	// })

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Scanln()
}
