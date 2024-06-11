package coordinates

import "math"

const R = 6378137.0

func Radians(deg float64) float64 {
	return deg * math.Pi / 180
}

func Degrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func Y2lat(ty int, zoom int) float64 {
	y := float64(ty) / math.Pow(2, float64(zoom))
	l := (1 - 2*y) * math.Pi
	o := math.Atan(math.Sinh(l))
	return Degrees(o)
}

func Lat2y(lat float64) float64 {
	return R * math.Log(math.Tan(math.Pi/4+Radians(lat)/2))
}

func X2lon(tx int, zoom int) float64 {
	x := float64(tx) / math.Pow(2, float64(zoom))
	l := (2*x - 1) * math.Pi
	return Degrees(l)
}

func Lon2x(lon float64) float64 {
	return R * Radians(lon)
}

func Distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	// See https://www.movable-type.co.uk/scripts/latlong.html
	lat1Rad := Radians(lat1)
	lat2Rad := Radians(lat2)
	deltaLatRad := Radians(lat2 - lat1)
	deltaLonRad := Radians(lon2 - lon1)

	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		lat1Rad*lat2Rad*math.Sin(deltaLonRad/2)*math.Sin(deltaLonRad/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
