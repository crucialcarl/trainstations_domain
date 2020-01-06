package storage

import "trainstations_domain/stations"
import "testing"

//import "fmt"

func TestGetAllStations(t *testing.T) {
	var stationRepository stations.Repository
	stationRepository = NewMemoryStationStorage()
	have, _ := stationRepository.GetAllStations()
	want := []stations.Station{{Abbr: "MONT", Name: "Montgomery"}}
	if (have[0].Name != want[0].Name) || (have[0].Abbr != want[0].Abbr) {
		t.Error(have[0], want[0])
	}
}