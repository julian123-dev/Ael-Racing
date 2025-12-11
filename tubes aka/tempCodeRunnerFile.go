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

// ITERATIF
func countCapitalIterative(s string) int {
	count := 0
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			count++
		}
	}
	return count
}

// REKURSIF
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

// MENGUKUR WAKTU
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

// Membuat teks random
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

	fmt.Println("=== HASIL PENGUJIAN JUMLAH HURUF KAPITAL ===")
	fmt.Println("+--------+-----------------------+-----------------------+")
	fmt.Printf("| %-6s | %-21s | %-21s |\n", "n", "Rekursif (ms)", "Iteratif (ms)")
	fmt.Println("+--------+-----------------------+-----------------------+")

	for i, n := range sizes {
		text := generateRandomText(n)

		// Rekursif
		rec := measureTime(func() {
			_ = countCapitalRecursive(text, 0)
		}, true)

		// Iteratif
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

	// Simpan ke file txt
	file, err := os.Create("TIMERRUNNING.txt")
	if err == nil {
		defer file.Close()

		fmt.Fprintln(file, "=== HASIL PENGUJIAN JUMLAH HURUF KAPITAL ===")
		fmt.Fprintln(file, "+--------+-----------------------+-----------------------+")
		fmt.Fprintf(file, "| %-6s | %-21s | %-21s |\n", "n", "Rekursif (ms)", "Iteratif (ms)")
		fmt.Fprintln(file, "+--------+-----------------------+-----------------------+")

		for i, n := range sizes {
			rec := recTimes[i]
			recStr := fmt.Sprintf("%.4f", rec)
			if rec < 0 {
				recStr = "ERR"
			}
			fmt.Fprintf(file, "| %6d | %21s | %21.4f |\n", n, recStr, iterTimes[i])
			fmt.Fprintln(file, "+--------+-----------------------+-----------------------+")
		}

		fmt.Fprintln(file, "\nCatatan: 'ERR' = rekursif gagal karena stack overflow.")
	}

	// =============================
	// BAGIAN PLOTTING (TIDAK MERAH)
	// =============================

	p, err := plot.New()
	if err != nil {
		fmt.Println("Gagal membuat plot:", err)
		return
	}

	p.Title.Text = "Perbandingan Waktu Eksekusi Rekursif vs Iteratif"
	p.X.Label.Text = "Panjang Input (n)"
	p.Y.Label.Text = "Waktu (ms)"

	// Membuat titik data
	ptsRec := make(plotter.XYs, 0)
	ptsIt := make(plotter.XYs, 0)

	for i, n := range sizes {
		if recTimes[i] >= 0 {
			ptsRec = append(ptsRec, plotter.XY{X: float64(n), Y: recTimes[i]})
		}
		ptsIt = append(ptsIt, plotter.XY{X: float64(n), Y: iterTimes[i]})
	}

	// Tambahkan grafik
	err = plotutil.AddLinePoints(p,
		"Rekursif", ptsRec,
		"Iteratif", ptsIt,
	)
	if err != nil {
		lr, _ := plotter.NewLine(ptsRec)
		li, _ := plotter.NewLine(ptsIt)
		p.Add(lr, li)
		p.Legend.Add("Rekursif", lr)
		p.Legend.Add("Iteratif", li)
	}

	// Atur rentang Y
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

	// Simpan gambar
	err = p.Save(10*vg.Inch, 5*vg.Inch, "GRAFIK.png")
	if err != nil {
		fmt.Println("Gagal menyimpan grafik:", err)
		return
	}

	fmt.Println("SUKSES! File GRAFIK.png dan TIMERRUNNING.txt berhasil dibuat.")
}
