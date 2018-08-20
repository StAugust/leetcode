// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"math"
	"net/http"
	"log"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, _, _ := corner(i+1, j)
			bx, by, _, _ := corner(i, j)
			cx, cy, _, _ := corner(i, j+1)
			dx, dy, _, _ := corner(i+1, j+1)
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")

	http.HandleFunc("/", handler)
	//http.ListenAndServe("localhost:8000", nil)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, _, _ := corner(i+1, j)
			bx, by, _, _ := corner(i, j)
			cx, cy, _, _ := corner(i, j+1)
			dx, dy, cr, cb := corner(i+1, j+1)
			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='#%02x00%02x' />\n",
				ax, ay, bx, by, cx, cy, dx, dy, cr, cb)
		}
	}
	fmt.Fprintf(w, "</svg>")
}

func corner(i, j int) (float64, float64, uint8, uint8) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f2(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	cr , cb := mapColor2(x, y)
	return sx, sy, cr, cb
}

func mapColor(x, y float64) (uint8, uint8) {
	height := math.Sin(math.Hypot(x, y))
	
	return uint8( height * 128 + 128 ),uint8( -height * 128 + 128)
}


func mapColor2(x, y float64) (uint8, uint8) {
	height := math.Sin(math.Hypot(x - float64(int(x) / 10 * 10) ,y - float64(int(y) / 10 * 10)))
	
	return uint8( height * 128 + 128 ),uint8( -height * 128 + 128)
}


func f2(x, y float64) float64 {
	r := math.Hypot(x - float64(int(x) / 10 * 10) ,y - float64(int(y) / 10 * 10))
	return math.Sin(r) / r
}


func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	res := math.Sin(r) / r
	if math.IsInf(res, 1) {
		return math.MaxFloat64
	}
	if math.IsInf(res, -1) {
		return -math.MaxFloat64
	}
	return math.Sin(r) / r
}

//!-
