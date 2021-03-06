package webserver

//import "fmt"
import "log"
import "net/http"
import "html/template"
import "trainstations_domain/stations"
import "trainstations_domain/lines"
import "trainstations_domain/trains"
import "trainstations_domain/scoring"
import "trainstations_domain/storage/file"
import "trainstations_domain/storage/bartapi"

type PageData struct {
	SelectedStationsNorth []trains.TrainInfo
	SelectedStationsSouth []trains.TrainInfo
	AllStations           []stations.Station
	AllLines              []lines.Line
}

func StartServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Printf("Starting server on %s", port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}

func home(rw http.ResponseWriter, r *http.Request) {
	stationRepository := file.NewStationStorage("./data/stations.txt")
	allStations, _ := stationRepository.GetAll()

	lineRepository := file.NewLineStorage("./data/lines.txt")
	allLines, _ := lineRepository.GetAll()

	_ = r.ParseForm()
	abbr := r.URL.Query().Get("abbr")
	lines := r.Form["line"]

	stationDataNorth := bartapi.TrainsFromBartAPI(abbr, "n")
	filteredTrainsNorth := filterDestinationByLine(stationDataNorth, lines)
	scoredTrainsNorth := scoring.Score(filteredTrainsNorth)

	stationDataSouth := bartapi.TrainsFromBartAPI(abbr, "s")
	filteredTrainsSouth := filterDestinationByLine(stationDataSouth, lines)
	scoredTrainsSouth := scoring.Score(filteredTrainsSouth)

	page := PageData{
		SelectedStationsNorth: scoredTrainsNorth,
		SelectedStationsSouth: scoredTrainsSouth,
		AllStations:           allStations,
		AllLines:              allLines,
	}

	templates := []string{
		"./ui/html/templates/index.tmpl",
		"./ui/html/templates/header.tmpl",
		"./ui/html/templates/base.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Panic("Error occurred parsing template", err)
	}

	err = tmpl.Execute(rw, page)
	if err != nil {
		log.Panic("Error occurred writing template", err)
	}
}

func filterDestinationByLine(trainsToFilter []trains.TrainInfo, lines []string) []trains.TrainInfo {
	var trainsMatchLine []trains.TrainInfo
	for _, train := range trainsToFilter {
		for _, line := range lines {
			if train.Line == line {
				trainsMatchLine = append(trainsMatchLine, train)
			}
		}
	}
	if len(trainsMatchLine) == 0 {
		return trainsToFilter
	}
	return trainsMatchLine
}
