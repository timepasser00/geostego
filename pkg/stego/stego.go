package stego

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/draw"
)

var ErrMessageTooLarge = errors.New("message is too large for this image")

var errHeaderNotFound = errors.New("header not found in the image")
var errInsfficientMessage = errors.New("message not complete")
var errBadImage = errors.New("image too small")

func Encode(img image.Image, msg []byte) (image.Image, error) {

	msgLen := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLen, uint32(len(msg)))

	data := append(msgLen, msg...)

	bounds := img.Bounds()

	height, width := bounds.Max.Y, bounds.Max.X
	pixelCnt := height * width

	bitsAvailable := pixelCnt * 3

	bitsRequired := len(data) * 8

	if bitsRequired > bitsAvailable {
		return nil, ErrMessageTooLarge
	}

	modifiedImg := image.NewRGBA(bounds)

	draw.Draw(modifiedImg, bounds, img, image.Point{}, draw.Src)

	dataBitIndex := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			if dataBitIndex >= bitsRequired {
				return modifiedImg, nil
			}

			r, g, b, a := modifiedImg.At(x, y).RGBA()

			c := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}

			setLSB := func(val *uint8) {

				if dataBitIndex >= bitsRequired {
					return
				}

				byteIndex := dataBitIndex / 8
				bitPos := 7 - (dataBitIndex % 8)

				bit := (data[byteIndex] >> byte(bitPos)) & 1

				*val = (*val & 0xFE) | bit

				dataBitIndex++
			}

			setLSB(&c.R)
			setLSB(&c.G)
			setLSB(&c.B)

			modifiedImg.SetRGBA(x, y, c)
		}
	}

	return modifiedImg, nil

}

func Decode(img image.Image) ([]byte, error) {

	bound := img.Bounds()

	height, width := bound.Max.Y, bound.Max.X

	maxEncodedBits := height * width * 3

	if maxEncodedBits < 32 {
		return nil, errHeaderNotFound
	}

	var (
		lengthHeader   uint32 = 0
		headerBitsRead        = 0
		readingBody           = false

		messageBuffer  []byte
		msgBitsRead         = 0
		currentMsgByte byte = 0
	)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			vals := []uint8{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
			}

			for _, val := range vals {
				bit := val & 1

				if !readingBody {
					lengthHeader = (lengthHeader << 1) | uint32(bit)
					headerBitsRead++

					if headerBitsRead == 32 {
						msgLen := int(lengthHeader)
						
						if (msgLen*8 + 32) >= maxEncodedBits {
							return nil, errInsfficientMessage
						}

						messageBuffer = make([]byte, msgLen)
						readingBody = true
					}
				} else {
					shift := 7 - (msgBitsRead % 8)
					currentMsgByte |= (bit << uint8(shift))
					msgBitsRead++

					if msgBitsRead%8 == 0 {
						byteIndex := (msgBitsRead) / 8 - 1
						messageBuffer[byteIndex] = currentMsgByte
						currentMsgByte = 0
					}

					if msgBitsRead >= int(lengthHeader)*8 {
						return messageBuffer, nil
					}
				}
			}
		}
	}
	return nil, errBadImage
}
