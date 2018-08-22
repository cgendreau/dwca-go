package dwc

import (
	"io"
	"log"
	"testing"
)

func TestNewDwcArchive(t *testing.T) {
	dwca, err := NewDwcArchive("testdata")
	if err != nil {
		t.Fatal(err)
	}

	if dwca.Core.RowType != "http://rs.tdwg.org/dwc/terms/Occurrence" {
		t.Fatal("RowType Occurrence Expected")
	}

	if dwca.Core.Files[0] != "occurrence.txt" {
		t.Fatal("File ocurrence.txt expected, got " + dwca.Core.Files[0])
	}

	if len(dwca.Extension) != 2 {
		t.Fatal("Expected 2 extensions")
	}

	if dwca.Extension[0].Files[0] != "identification.txt" {
		t.Fatal("First extension file identification.txt expected, got " + dwca.Extension[0].Files[0])
	}

	if dwca.Core.IndexOf("http://rs.tdwg.org/dwc/terms/datasetID") != 6 {
		t.Fatal("datasetID Index should be 6 ")
	}
}

func TestDwcFileReader(t *testing.T) {
	testFolder := "testdata"

	dwca, err := NewDwcArchive(testFolder)
	if err != nil {
		t.Fatal(err)
	}

	csvReader, err := dwca.Core.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer dwca.Core.Close()

	numberOfLine := 0
	for {
		_, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		numberOfLine++
	}

	if numberOfLine != 21 {
		t.Fatal("File should should have 21 lines")
	}
}
