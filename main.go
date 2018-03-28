package main

import (
	"golang.org/x/exp/io/spi"
	"time"
	"github.com/btittelbach/go-bbhw"
	"bufio"
	"io/ioutil"
	"os"
)

var pins[]bbhw.FakeGPIO
//var readings[]float64
var activation bbhw.FakeGPIO
var control bbhw.PWMPin
var activePos, lastPos float64
var recordingComplete bool
var input spi.Device

func main () {

	input, err := spi.Open(&spi.Devfs{
		Dev:		"dev/spidev0.0",
		Mode:		spi.Mode3,		//SPI devices designated as spidev[port #][device #]
		MaxSpeed:	22100,
	})
	if err != nil {
		panic(err)}

	input.SetBitsPerWord(10)//10 bits per ADC packet
	defer input.Close()			//Close accelerometer SPI comms after program is over

	pins := make([]bbhw.FakeGPIO, 15) //6 readings from passive materials, 1 reading from active, 1 for powering active
	pins[0], pins[1] = null                     //Passive #1 - Control
	pins[2], pins[3] = null                     //Passive #2 - MaPS
	pins[4], pins[5] = null                     //Passive #3 - Cable Mount
	pins[6], pins[7] = null                     //Passive #4 - Neoprene
	pins[8], pins[9] = null                     //Passive #5 - Sorbothane (POSSIBLY REMOVING)
	pins[10], pins[11] = null                   //Passive #6 - Elastomeric Resin
	pins[12], pins[13] = null                   //Active Read
	pins[14], pins[15] = null                   //Active Control

	activePos = 0	//Initialize readings to 0

	countdown := time.NewTimer(30)	//Starting countdown (can be interrupted by activation line)
	go func() {					//Function called when countdown finishes
		<-countdown.C			//wait for the countdown to finish
		recordReact()}()		//begin recordReact after the countdown finishes

	start := false
	for start == false {					//Loops while waiting for an activation
		start, err = activation.GetState()	//Checks activation state, true ends loop
		if err != nil {
			panic(err)
		}
	}

	if start == true{			//if the activation line reads HIGH
		countdown.Stop()		//stop the timer
		countdown = time.NewTimer(15 * time.Minute)	//Begins timer for 15 minute recording
		recordingComplete = false					//Flag for finishing recording
		go func(){
			<- countdown.C
			recordingComplete = true				//Trigger recording complete
		}()
		recordReact()			//Begin the recording and feedback loop
	}
}

func recordReact(){
	lastPos = activePos 	//Set last value for active system
	  /*1. Cycle through Chip Select
		2. New accelerometer feeds data to SPI
	  	3. Record incoming data to file (maybe do several files, maybe do single file that will then require splitting into six parts
	  	4. Use data from active system's accelerometer to drive voice coil*/
	log, err := os.Create("CompiledAccelerometerData.txt")
	if err != nil {panic(err)}
	var reading[] byte					//Input byte array
	for !recordingComplete{
		for _, pin := range pins{
			pin.SetState(true)		//Assuming true means HIGH
			input.Tx(nil, reading)	//Attempt SPI transaction, place results into reading
			pin.SetState(false)



			//driveVoiceCoil()
		}
	}
}

func driveVoiceCoil(difference int){	//Drives the actuator based on difference in current value and last recorded value
	bbhw.SetPWMFreq(control, float64(difference)) //Set PWM frequencies related to some function done to the differences

}


