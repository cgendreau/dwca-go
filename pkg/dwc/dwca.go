package dwc

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"
)

const (
	MetaXMLFileName = "meta.xml"
)

type Field struct {
	Index byte   `xml:"index,attr"`
	Term  string `xml:"term,attr"`
}

// DwcFile represents the definition of a single file within an archive
type DwcFile struct {
	Encoding           string   `xml:"encoding,attr"`
	FieldsTerminatedBy string   `xml:"fieldsTerminatedBy,attr"`
	LinesTerminatedBy  string   `xml:"linesTerminatedBy,attr"`
	FieldsEnclosedBy   string   `xml:"fieldsEnclosedBy,attr"`
	IgnoreHeaderLines  string   `xml:"ignoreHeaderLines,attr"`
	RowType            string   `xml:"rowType,attr"`
	Files              []string `xml:"files>location"`
	Fields             []Field  `xml:"field"`
	FieldsMap          map[string]byte
	path               string
	file               *os.File
}

//DwcArchive represents the metadata of a DarwinCore Archive (how the data is structured inside the archive)
type DwcArchive struct {
	Core struct {
		DwcFile
		ID Field `xml:"id"`
	} `xml:"core"`
	Extension []struct {
		DwcFile
		CoreID Field `xml:"coreid"`
	} `xml:"extension"`
	Metadata string `xml:"metadata,attr"`
	path     string
}

func (file *DwcArchive) postDecode(rootPath string) error {
	file.path = rootPath

	if err := file.Core.postDecode(rootPath); err != nil {
		return err
	}

	if file.Extension != nil {
		for _, ext := range file.Extension {
			if err := ext.postDecode(rootPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (dwcFile *DwcFile) postDecode(rootPath string) error {
	if dwcFile.Fields == nil {
		return nil
	}
	dwcFile.FieldsMap = make(map[string]byte)
	for _, f := range dwcFile.Fields {
		dwcFile.FieldsMap[f.Term] = f.Index
	}

	if len(dwcFile.Files) > 1 {
		return errors.New("dwc package only supports 1 location per file")
	}

	dwcFile.path = filepath.Join(rootPath, dwcFile.Files[0])

	return nil
}

// NewDwcArchive opens a folder, reads the meta xml file and returns an initialized *DwcArchive
func NewDwcArchive(name string) (*DwcArchive, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New("This function only accepts folder for now")
	}

	//open meta.xml file to read the description of the archive
	file, err := os.Open(filepath.Join(name, MetaXMLFileName))

	if err != nil {
		return nil, err
	}
	defer file.Close()

	dwca := &DwcArchive{}
	if err := xml.NewDecoder(file).Decode(dwca); err != nil {
		return nil, err
	}

	if err := dwca.postDecode(name); err != nil {
		return nil, err
	}

	return dwca, nil
}

// IndexOf returns the index of term within the archive
func (dwcFile DwcFile) IndexOf(term string) byte {
	return dwcFile.FieldsMap[term]
}

//Open returns a new csv *Reader for the current DwcFile
func (dwcFile *DwcFile) Open() (*csv.Reader, error) {

	file, err := os.Open(dwcFile.path)
	if err != nil {
		return nil, err
	}
	// keep the pointer, we will close it later
	dwcFile.file = file

	csvReader := csv.NewReader(file)

	//Handle the field delimiter
	newstr, err := strconv.Unquote(`"` + dwcFile.FieldsTerminatedBy + `"`)
	fieldDelimiter, _ := utf8.DecodeRuneInString(newstr)
	csvReader.Comma = fieldDelimiter

	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1

	return csvReader, nil
}

//Close closes the underlying resource
func (dwcFile *DwcFile) Close() {
	if dwcFile.file != nil {
		dwcFile.file.Close()
	}
}
