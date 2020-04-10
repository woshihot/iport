package iport

import (
	"isesol.com/iport/options"
	"testing"
)

func TestCreate(t *testing.T) {
	o := options.NewOption()
	o.BoxName = "BOX0318060105"
	o.BoxLicense = "1"
	o.BoxMac = "1"

	i := Create(o)
	i.Start()
	select {}

}
