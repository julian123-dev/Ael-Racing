package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// countCapitalIterative counts uppercase letters A-Z iteratively
func countCapitalIterative(s string) int {
	count := 0
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			count++
		}
	}
	return count
}

// countCapitalRecursive counts uppercase letters A-Z recursively
func countCapitalRecursive(s string, idx int) int {
	if idx >= len(s) {
		return 0
	}
	c := 0
	if s[idx] >= 'A' && s[idx] <= 'Z' {
		c = 1
	}
	return c + countCapitalRecursive(s, idx+1)
}

// measureTime runs f and returns elapsed time in milliseconds.
func measureTime(f func(), isRecursive bool) float64 {
	start := time.Now()
	if isRecursive {
		ch := make(chan float64)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					ch <- -1
				}
			}()
			f()
			ch <- float64(time.Since(start).Nanoseconds()) / 1e6
		}()
		return <-ch
	}
	f()
	return float64(time.Since(start).Nanoseconds()) / 1e6
}

// generateRandomText returns a random string
func generateRandomText(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	sizes := []int{3000, 8000, 31000, 49000, 56000, 131000}

	iterTimes := make([]float64, len(sizes))
	recTimes := make([]float64, len(sizes))

	// print table header
	fmt.Println("=== HASIL PENGUJIAN JUMLAH HURUF KAPITAL ===")
	fmt.Println("+--------+-----------------------+-----------------------+")
	fmt.Printf("| %-6s | %-21s | %-21s |\n", "n", "Waktu Rekursif (ms)", "Waktu Iteratif (ms)")
	fmt.Println("+--------+-----------------------+-----------------------+")

	for i, n := range sizes {
		text := generateRandomText(n)

		rec := measureTime(func() {
			_ = countCapitalRecursive(text, 0)
		}, true)

		it := measureTime(func() {
			_ = countCapitalIterative(text)
		}, false)

		recTimes[i] = rec
		iterTimes[i] = it

		recStr := fmt.Sprintf("%.4f", rec)
		if rec < 0 {
			recStr = "ERR"
		}

		fmt.Printf("| %6d | %21s | %21.4f |\n", n, recStr, it)
		fmt.Println("+--------+-----------------------+-----------------------+")
	}

	// save table
	tf, err := os.Create("TIMERRUNNING.txt")
	if err == nil {
		defer tf.Close()
		fmt.Fprintln(tf, "=== HASIL PENGUJIAN JUMLAH HURUF KAPITAL ===")
		fmt.Fprintln(tf, "+--------+-----------------------+-----------------------+")
		fmt.Fprintf(tf, "| %-6s | %-21s | %-21s |\n", "n", "Waktu Rekursif (ms)", "Waktu Iteratif (ms)")
		fmt.Fprintln(tf, "+--------+-----------------------+-----------------------+")
		for i, n := range sizes {
			rec := recTimes[i]
			recStr := fmt.Sprintf("%.4f", rec)
			if rec < 0 {
				recStr = "ERR"
			}
			fmt.Fprintf(tf, "| %6d | %21s | %21.4f |\n", n, recStr, iterTimes[i])
			fmt.Fprintln(tf, "+--------+-----------------------+-----------------------+")
		}
		fmt.Fprintln(tf, "\nCatatan: 'ERR' berarti fungsi rekursif stack overflow.")
	}

	// ---------- PLOT ----------
	p := plot.New() // ← BENAR, 1 RETURN VALUE SAJA

	p.Title.Text = "Perbandingan Waktu Eksekusi: Rekursif vs Iteratif"
	p.X.Label.Text = "Panjang input (n)"
	p.Y.Label.Text = "Waktu (ms)"

	ptsRec := plotter.XYs{}
	ptsIt := plotter.XYs{}
	for i, n := range sizes {
		if recTimes[i] >= 0 {
			ptsRec = append(ptsRec, plotter.XY{X: float64(n), Y: recTimes[i]})
		}
		ptsIt = append(ptsIt, plotter.XY{X: float64(n), Y: iterTimes[i]})
	}

	err = plotutil.AddLinePoints(p, "Rekursif", ptsRec, "Iteratif", ptsIt)
	if err != nil {
		lr, _ := plotter.NewLine(ptsRec)
		li, _ := plotter.NewLine(ptsIt)
		p.Add(lr, li)
		p.Legend.Add("Rekursif", lr)
		p.Legend.Add("Iteratif", li)
	}

	p.Y.Min = 0
	maxY := 0.0
	for _, v := range iterTimes {
		if v > maxY {
			maxY = v
		}
	}
	for _, v := range recTimes {
		if v > maxY {
			maxY = v
		}
	}
	p.Y.Max = math.Ceil(maxY*1.2 + 1)

	if err := p.Save(10*vg.Inch, 5*vg.Inch, "GRAFIK.png"); err != nil {
		fmt.Println("Gagal menyimpan grafik:", err)
		return
	}

	fmt.Println("SUKSES! File 'GRAFIK.png' dan 'TIMERRUNNING.txt' dibuat.")
}
