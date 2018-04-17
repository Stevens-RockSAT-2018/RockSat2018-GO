package main

import (
	"fmt"
	"os"
	"time"

	"github.com/btittelbach/go-bbhw"
	"golang.org/x/exp/io/spi"
)

var log os.File

var pins []bbhw.FakeGPIO

//var readings[]float64
var activation bbhw.GPIOControllablePin
var control bbhw.PWMPin
var activePos, lastPos float64
var recordingComplete bool
var input *spi.Device

// #define SETUP_FLAG     0b10000000
// #define SCAN_MODE_NONE 0b00000110
// #define SCAN_MODE_0_N  0b00000000
// #define SCAN_MODE_N_4 0b00000100
const (
	SetupFlag    = 0x80
	ScanModeNone = 0x06
	ScanMode0N   = 0x00
	ScanModeN4   = 0x04
)

func readWriteByte(b byte) byte {
	output := []byte{0}
	input.Tx([]byte{b}, output)
	return output[0]
}

func readAccel() int {
	var data int
	regData := SetupFlag | 1<<3 | ScanMode0N
	msb := readWriteByte(byte(regData))
	// time.Sleep(time.Microsecond * 20)
	lsb := readWriteByte(0x00)
	data = (int(msb)<<6 | int(lsb)>>2) << 16

	msb = readWriteByte(0x00)
	// time.Sleep(time.Microsecond * 20)
	lsb = readWriteByte(0x00)
	data |= (int(msb)<<6 | int(lsb)>>2)

	return data
}

func main() {
	var err error

	input, err = spi.Open(&spi.Devfs{
		Dev:  "/dev/spidev1.0",
		Mode: spi.Mode0,
		// Mode:     spi.Mode3, //SPI devices designated as spidev[port #][device #]
		MaxSpeed: 22100,
	})
	if err != nil {
		panic(err)
	}

	// write := [...]byte{1, 2, 3, 4}
	// read := [...]byte{}

	//
	// if err := input.Tx([]byte{
	// 	1, 2, 3, 4,
	// }, nil); err != nil {
	// 	panic(err)
	// }
	// input.SetBitsPerWord(10) //10 bits per ADC packet

	for true {
		fmt.Println(readAccel())
		time.Sleep(time.Second)
	}

	defer input.Close() //Close accelerometer SPI comms after program is over

}
