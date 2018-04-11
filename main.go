package main

import (
	"github.com/btittelbach/go-bbhw"
	"golang.org/x/exp/io/spi"
	// "bufio"
)

var pins []bbhw.FakeGPIO

//var readings[]float64
var activation bbhw.FakeGPIO
var control bbhw.PWMPin
var activePos, lastPos float64
var recordingComplete bool
var input spi.Device

func main() {

	input, err := spi.Open(&spi.Devfs{
		Dev:      "/dev/spidev1.0",
		Mode:     spi.Mode3, //SPI devices designated as spidev[port #][device #]
		MaxSpeed: 22100,
	})
	if err != nil {
		panic(err)
	}

	// write := [...]byte{1, 2, 3, 4}
	// read := [...]byte{}

	if err := input.Tx([]byte{
		1, 2, 3, 4,
	}, nil); err != nil {
		panic(err)
	}
	// input.SetBitsPerWord(10) //10 bits per ADC packet
	defer input.Close() //Close accelerometer SPI comms after program is over

}
