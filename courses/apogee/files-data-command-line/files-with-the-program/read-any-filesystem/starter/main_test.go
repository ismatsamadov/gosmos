package main

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestReadStationsFS(t *testing.T) {
	sqt := Station{Code: "SQT", Name: "Salar Quieto", MaskDeg: 3, Bands: []Band{"S", "X"}}
	small := fstest.MapFS{"stations.csv": {Data: []byte("code,name,lat_deg,lon_deg,mask_deg,bands\nSQT,Salar Quieto,-23.4,-104.6,3,S X\n")}}
	got, err := readStationsFS(small)
	if err != nil || !reflect.DeepEqual(got, []Station{sqt}) {
		t.Errorf("readStationsFS(a table of one station) = %v, %v, want %v, nil", got, err, []Station{sqt})
	}

	built, err := readStationsFS(resources)
	if err != nil || len(built) != 5 || !reflect.DeepEqual(built[2], sqt) {
		t.Fatalf("readStationsFS(resources) = %v, %v, want 5 stations with SQT third, nil", built, err)
	}
	data, err := resources.ReadFile("stations.csv")
	if err != nil {
		t.Fatal(err)
	}
	copied, err := readStationsFS(fstest.MapFS{"stations.csv": {Data: data}})
	if err != nil || !reflect.DeepEqual(copied, built) {
		t.Errorf("readStationsFS(the same table in memory) = %v, %v, want %v, nil", copied, err, built)
	}

	got, err = readStationsFS(fstest.MapFS{"notes.txt": {Data: []byte("no table here\n")}})
	if got != nil || !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("readStationsFS(no stations.csv) = %v, %v, want nil and a not-exist error", got, err)
	}
}
