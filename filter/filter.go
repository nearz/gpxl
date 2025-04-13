package filter

import (
	"github.com/nearz/gpxl/pxl"
)

type Filter interface {
	Render(*pxl.Pxl)
}
