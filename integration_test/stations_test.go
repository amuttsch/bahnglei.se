package integrationtest

import (
	"testing"
)

func TestKarlsruhe(t *testing.T) {
	repo, ctx := RepoTestCase(t)

	stationData := getStationData(t, ctx, repo, 2574283615)

	if stationData.Station.Name != "Karlsruhe Hauptbahnhof" {
		t.Errorf("Expected name to be 'Karlsruhe Hauptbahnhof', got %s", stationData.Station.Name)
	}

	if stationData.Station.Tracks != 16 {
		t.Errorf("Expected 'Karlsruhe Hauptbahnhof' to have 16 tracks, got %d", stationData.Station.Tracks)
	}

	stopPositionPlatformOne := GetStopPosition("1", stationData.StopPositions)
	if stopPositionPlatformOne == nil {
		t.Fatalf("Could not find stop position 1 in Karlsruhe Hauptbahnhof")
	}

	if stopPositionPlatformOne.Neighbors != "1;2" {
		t.Fatalf("Karlsruhe Hauptbahnhof platform 1 should have neighbor platform 2, got %s", stopPositionPlatformOne.Neighbors)
	}
}

func TestKarlsruheDurlach(t *testing.T) {
	repo, ctx := RepoTestCase(t)

	stationData := getStationData(t, ctx, repo, 2958859136)

	if stationData.Station.Tracks != 5 {
		t.Errorf("Expected 'Karlsruhe Durlach' to have 5 tracks, got %d", stationData.Station.Tracks)
	}

	stopPositionPlatformOne := GetStopPosition("1", stationData.StopPositions)
	if stopPositionPlatformOne == nil {
		t.Fatalf("Could not find stop position 1 in Karlsruhe Durlach")
	}
	if stopPositionPlatformOne.Neighbors != "1" {
		t.Fatalf("Karlsruhe Durlach platform 1 does not have neighbors, got %s", stopPositionPlatformOne.Neighbors)
	}

	stopPositionPlatformTwo := GetStopPosition("2", stationData.StopPositions)
	if stopPositionPlatformTwo == nil {
		t.Fatalf("Could not find stop position 2 in Karlsruhe Durlach")
	}
	if stopPositionPlatformTwo.Neighbors != "2;5" {
		t.Fatalf("Karlsruhe Hauptbahnhof platform 2 should have neighbor platform 5, got %s", stopPositionPlatformOne.Neighbors)
	}

	stopPositionPlatformSix := GetStopPosition("6", stationData.StopPositions)
	if stopPositionPlatformSix == nil {
		t.Fatalf("Could not find stop position 6 in Karlsruhe Durlach")
	}
	if stopPositionPlatformSix.Neighbors != "6;9" {
		t.Fatalf("Karlsruhe Hauptbahnhof platform 6 should have neighbor platform 9, got %s", stopPositionPlatformOne.Neighbors)
	}
}

func TestDarmstadtSued(t *testing.T) {
	repo, ctx := RepoTestCase(t)

	stationData := getStationData(t, ctx, repo, 6673320792)

	if stationData.Station.Tracks != 2 {
		t.Errorf("Expected 'Darmstadt Süd' to have 2 tracks, got %d", stationData.Station.Tracks)
	}

	stopPositionPlatformOne := GetStopPosition("1", stationData.StopPositions)
	if stopPositionPlatformOne == nil {
		t.Fatalf("Could not find stop position 1 in Darmstadt Süd")
	}
	if stopPositionPlatformOne.Neighbors != "1;2" {
		t.Fatalf("Darmstadt Süd platform 1 does have neighbor platform 2, got %s", stopPositionPlatformOne.Neighbors)
	}
}
