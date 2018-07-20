package dwc

import (
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
)

const (
	metaXmlFileName = "meta.xml"
)

type Field struct {
	Index string `xml:"index,attr"`
	Term  string `xml:"term,attr"`
}

type DwcFile struct {
	Encoding           string `xml:"encoding,attr"`
	FieldsTerminatedBy string `xml:"fieldsTerminatedBy,attr"`
	LinesTerminatedBy  string `xml:"linesTerminatedBy,attr"`
	FieldsEnclosedBy   string `xml:"fieldsEnclosedBy,attr"`
	IgnoreHeaderLines  string `xml:"ignoreHeaderLines,attr"`
	RowType            string `xml:"rowType,attr"`
}

type DwcArchive struct {
	Core struct {
		DwcFile
		Id     Field   `xml:"id"`
		Fields []Field `xml:"field"`
	} `xml:"core"`
	Extension []struct {
		DwcFile
		CoreId Field   `xml:"coreid"`
		Fields []Field `xml:"field"`
	} `xml:"extension"`
	Metadata string `xml:"metadata,attr"`
}

func Open(name string) (*DwcArchive, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New("This function only accepts folder for now")
	}

	file, err := os.Open(filepath.Join(name, metaXmlFileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metaxml DwcArchive
	err2 := xml.NewDecoder(file).Decode(&metaxml)
	if err2 != nil {
		return nil, err
	}
	return &metaxml, nil
}
