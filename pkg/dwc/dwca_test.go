package dwc

import (
	"testing"
)

func TestOpen(t *testing.T) {
	dwca, err := Open("testdata")
	if err != nil {
		t.Fatal(err)
	}

	if dwca.Core.RowType != "http://rs.tdwg.org/dwc/terms/Occurrence" {
		t.Fatal("RowType Occurrence Expected")
	}
}
