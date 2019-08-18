package main

import (
	"fmt"
	"flag"
	"io"
	"os"
	"time"

	"github.com/google/gopacket/pcap"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	/*
		var values plotter.Values
		for i := 0; i < 1000; i++ {
			values = append(values, rand.NormFloat64())
		}
	*/

	var bpf string

	flag.Usage = func(){
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [-bpf <bpf filter>] /path/to/packet/capture/file1 [.../file2]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&bpf, "bpf", "", "Process only packets which match this BPF, all packets if left empty.")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}
	plt.Title.Text = "jitter (nanoseconds)"

	index := uint(0)
	for _, pcapPath := range flag.Args() {
		h, err := pcap.OpenOffline(pcapPath)
		if err != nil {
			panic(fmt.Sprintf("failed to open packet capture file: %v", err))
		}
		if bpf != "" {
			if err = h.SetBPFFilter(bpf); err != nil {
				panic(fmt.Sprintf("failed to set BPF filter: %v", err))
			}
		}

		var values plotter.Values

		var t [2]time.Time
		for i := uint64(0); ; i++ {
			_, ci, err := h.ReadPacketData()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(fmt.Sprintf("failed to read packet: %v", err))
			}
			t[i%2] = ci.Timestamp
			if i >= 1 {
				delta := t[i%2].Sub(t[(i-1)%2]) // subtract previous timestamp from current
				values = append(values, float64(delta))
			}
		}

		boxPlot, err := plotter.NewBoxPlot(1*vg.Centimeter, float64(index), values)
		if err != nil {
			panic(err)
		}
		plt.Add(boxPlot)

		h.Close()

		index++
	}

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "jitter.png"); err != nil {
		panic(err)
	}
}
