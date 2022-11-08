package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

// https://docs.fileformat.com/audio/wav/
func main() {

	outputFile, err := os.Create("out.wav")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() { outputFile.Close() }()

	err = WriteChunk(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}

const (
	RiffHeaderSize         = 8
	SoundDuration  float64 = 10.0 // seconds

	// Frequency with 440 pitch standard https://en.wikipedia.org/wiki/A440_(pitch_standard)
	Frequency float64 = 440.0 // Hertz/Hz | 440 Hz, or 440 cycles per second
	Volume    float64 = 0.5   // 0.0 silent to 1.0 max
	Channels          = 2

	// SamplesPerSecond
	// Sample (convert continuous data to discrete data) https://en.wikipedia.org/wiki/Sampling_(signal_processing)
	// PCM https://en.wikipedia.org/wiki/Pulse-code_modulation
	SamplesPerSecond = 44100

	// BitPerSample / Bit Depth https://en.wikipedia.org/wiki/Audio_bit_depth
	BitsPerSample = 16
)

func WriteChunk(out io.Writer) error {

	format, formatSize, err := writeFormatChunk()
	if err != nil {
		return err
	}

	data, dataSize, err := writeDataChunk()
	if err != nil {
		return err
	}

	header, err := writeHeaderChunk(formatSize, dataSize)
	if err != nil {
		return err
	}

	allChunk := append(header, format...)
	allChunk = append(allChunk, data...)

	_, err = out.Write(allChunk)
	if err != nil {
		return err
	}

	return nil
}

func writeOffset(offset int, data, input []byte) {
	for i := 0; i < len(data); i++ {
		input[offset:][i] = data[i]
	}
}

func writeHeaderChunk(formatSize, dataSize uint32) ([]byte, error) {
	packet := make([]byte, RiffHeaderSize+4)

	// RIFF ascii
	writeOffset(0, []byte{0x52, 0x49, 0x46, 0x46}, packet)

	// total size
	totalSize := (4 + RiffHeaderSize) + (formatSize + RiffHeaderSize) + (dataSize + RiffHeaderSize)
	fmt.Println("totalSize: ", totalSize)
	writeOffset(4, EncodeUint32LE(totalSize), packet)

	// WAVE ascii
	writeOffset(8, []byte{0x57, 0x41, 0x56, 0x45}, packet)

	return packet, nil
}

func writeFormatChunk() ([]byte, uint32, error) {
	var size uint32 = 16
	packet := make([]byte, RiffHeaderSize+size)

	var (
		dataRate       uint32 = (Channels * BitsPerSample * SamplesPerSecond) / 8
		bytesPerSample uint32 = (BitsPerSample-1)/8 + 1
		blockAlignment uint16 = Channels * uint16(bytesPerSample)
	)

	fmt.Println("dataRate: ", dataRate)
	fmt.Println("bytesPerSample: ", bytesPerSample)
	fmt.Println("blockAlignment: ", blockAlignment)

	// fmt ascii
	writeOffset(0, []byte{0x66, 0x6d, 0x74, 0x20}, packet)

	writeOffset(4, EncodeUint32LE(16), packet)
	writeOffset(8, EncodeUint16LE(1), packet)
	writeOffset(10, EncodeUint16LE(Channels), packet)
	writeOffset(12, EncodeUint32LE(SamplesPerSecond), packet)
	writeOffset(16, EncodeUint32LE(dataRate), packet)
	writeOffset(20, EncodeUint16LE(blockAlignment), packet)
	writeOffset(22, EncodeUint16LE(BitsPerSample), packet)

	return packet, size, nil
}

func writeDataChunk() ([]byte, uint32, error) {
	var (
		sampleCount    uint32 = uint32(SoundDuration) * SamplesPerSecond
		bytesPerSample uint32 = (BitsPerSample-1)/8 + 1
		dataSize       uint32 = sampleCount * bytesPerSample * Channels
		i              uint32
		channel        int
		offset         = 0
	)

	fmt.Println("--------")
	fmt.Println("sampleCount: ", sampleCount)
	fmt.Println("bytesPerSample: ", bytesPerSample)
	fmt.Println("dataSize: ", dataSize)

	packet := make([]byte, RiffHeaderSize+dataSize)

	// data ascii
	writeOffset(0, []byte{0x64, 0x61, 0x74, 0x61}, packet)

	// data size
	writeOffset(4, EncodeUint32LE(dataSize), packet)

	// https://en.wikipedia.org/wiki/Sampling_(signal_processing)
	// convert continuous data to discrete data
	for i = 0; i < sampleCount; i++ {
		var (
			sample int
			freq   float64 = Frequency
			w      float64 = 2 * math.Pi * float64(i) / (SamplesPerSecond / freq)
			vol    float64 = math.Sin(w) * Volume
		)
		// fmt.Println("w: ", w)

		// if BitsPerSample less than or equal 8
		// store sample as 8 bit integer sign or unsigned
		// otherwise
		// store sample as 16 bit (that range from -32768 to 32767) integer signed or unsigned Little Endian
		if BitsPerSample <= 8 {
			var rangeTop int = (1 << BitsPerSample) - 1
			sample = int(((vol + 1) / 2) * float64(rangeTop))
		} else {
			var (
				rangeTopFull int = (1 << BitsPerSample) - 1
				rangeTop     int = 1 << (BitsPerSample - 1)
			)

			sample = int((((vol + 1) / 2) * float64(rangeTopFull)) - float64(rangeTop))
		}

		for channel = 0; channel < Channels; channel++ {
			switch bytesPerSample {
			case 1:
				// todo
				fmt.Println("1 ", sample)
				writeOffset(8+offset, []byte{byte(sample)}, packet)
				break
			case 2:
				// todo
				writeOffset(8+offset, EncodeInt16LE(int16(sample)), packet)
				// fmt.Println("2 ", sample, "offset ", 8+offset, "vol ", vol, "bd8: ", ((vol + 1) / 2) * float64((1 << 8) - 1))
				break
			case 3:
				// todo
				writeOffset(8+offset, EncodeUint24LE(uint32(sample)), packet)
				break
			case 4:
				// todo
				writeOffset(8+offset, EncodeUint32LE(uint32(sample)), packet)
				break
			default:
				return nil, dataSize, errors.New("invalid bytesPerSample")
			}

			offset = offset + int(bytesPerSample)
		}
	}

	return packet, dataSize, nil
}

func EncodeUint16LE(x uint16) []byte {
	bytes := make([]byte, 2)
	for i := 0; i < 2; i++ {
		shift := i * 8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeInt16LE(x int16) []byte {
	bytes := make([]byte, 2)
	for i := 0; i < 2; i++ {
		shift := i * 8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeUint16BE(x uint16) []byte {
	bytes := make([]byte, 2)
	for i := 0; i < 2; i++ {
		shift := 8 - i*8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeInt16BE(x int16) []byte {
	bytes := make([]byte, 2)
	for i := 0; i < 2; i++ {
		shift := 8 - i*8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeUint24LE(x uint32) []byte {
	bytes := make([]byte, 3)
	for i := 0; i < 3; i++ {
		shift := i * 8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeInt24LE(x int32) []byte {
	bytes := make([]byte, 3)
	for i := 0; i < 3; i++ {
		shift := i * 8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeUint24BE(x uint32) []byte {
	bytes := make([]byte, 3)
	for i := 0; i < 3; i++ {
		shift := 16 - i*8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeInt24BE(x uint32) []byte {
	bytes := make([]byte, 3)
	for i := 0; i < 3; i++ {
		shift := 16 - i*8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeUint32LE(x uint32) []byte {
	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		shift := i * 8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeInt32LE(x int32) []byte {
	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		shift := i * 8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeUint32BE(x uint32) []byte {
	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		shift := 24 - i*8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func EncodeInt32BE(x int32) []byte {
	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		shift := 24 - i*8
		c := (x >> shift) & 0xff
		bytes[i] = byte(c)
	}
	return bytes
}

func DecodeUint32LE(bytes []byte) uint32 {
	var d uint32 = 0
	for i, b := range bytes {
		shift := i * 8
		c := uint32(b) << shift
		d |= c
	}
	return d
}

func DecodeUint32BE(bytes []byte) uint32 {
	var d uint32 = 0
	for i, b := range bytes {
		shift := 24 - i*8
		c := uint32(b) << shift
		d |= c
	}
	return d
}
