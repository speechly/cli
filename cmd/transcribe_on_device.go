//go:build on_device
// +build on_device

package cmd

/*
 #cgo CFLAGS: -I${SRCDIR}/../decoder/include
 #cgo darwin LDFLAGS: -L${SRCDIR}/../decoder/lib -Wl,-rpath,decoder/lib -lspeechlyDecoder -lz -framework Foundation -lc++ -framework Security
 #cgo linux LDFLAGS: -L${SRCDIR}/../decoder/lib -Wl,-rpath,$ORIGIN/../decoder/lib -Wl,--start-group -lstdc++ -lpthread -ldl -lm -lspeechly -lz
 #cgo tflite LDFLAGS: -ltensorflowlite_c
 #cgo coreml LDFLAGS: -framework coreml
 #include <Decoder.h>
 #include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"os"
	"io"
	"path"
	"strings"
	"unsafe"

	"github.com/go-audio/audio"
)

func transcribeOnDevice(model string, corpusPath string) ([]AudioCorpusItem, error) {
	df, err := NewDecoderFactory(model)
	if err != nil {
		return nil, err
	}

	if corpusPath == "STDIN" {
		d, _ := df.NewStream("")
		decodeStdin(d)
	}

	ac, err := readAudioCorpus(corpusPath)
	if err != nil {
		return nil, err
	}

	bar := getBar("Transcribing", "utt", len(ac))
	var results []AudioCorpusItem
	for _, aci := range ac {
		d, err := df.NewStream("")
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}

		audioFilePath := path.Join(path.Dir(corpusPath), aci.Audio)
		if corpusPath == aci.Audio {
			audioFilePath = corpusPath
		}
		transcript, err := decodeAudioCorpusItem(audioFilePath, aci, d)
		if err != nil {
			barClearOnError(bar)
			return results, err
		}

		err = bar.Add(1)
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}

		results = append(results, AudioCorpusItem{Audio: aci.Audio, Hypothesis: transcript})
	}

	return results, nil
}

func decodeStdin(d *cDecoder) (error) {
	cErr := C.DecoderError{}
	C.Decoder_EnableVAD(d.decoder, 1, &cErr);

	go func () {
		cErr := C.DecoderError{}
		buffer := make([]byte, 4096)
		sampleBuffer := make([]float32, 2048)
		for {
			if _, err := io.ReadFull(os.Stdin, buffer); err != nil {
				fmt.Println("error:", err)
				break;
			}
			bufferPos := 0
			for i := 0; i < 4096; i += 2 {
				s := int16((uint16(buffer[i]) | (uint16(buffer[i + 1]) << 8)))
				// fmt.Println(s)
				sampleBuffer[bufferPos] = float32(s) / 32768.0
				bufferPos++
			}
			C.Decoder_WriteSamples(d.decoder, (*C.float)(unsafe.Pointer(&sampleBuffer[0])), C.size_t(2048), C.int(0), &cErr)
		}
	}()

	for {
		res := C.Decoder_WaitResults(d.decoder, &cErr)
		if cErr.error_code != C.uint(0) {
			return fmt.Errorf("failed reading transcript from decoder, error code %d", cErr.error_code)
		}
		word := C.GoString(res.word)
		fmt.Println(word)
		C.CResultWord_Destroy(res)
	}

	return nil
}

func decodeAudioCorpusItem(audioFilePath string, aci AudioCorpusItem, d *cDecoder) (string, error) {
	cErr := C.DecoderError{}

	err := readAudio(audioFilePath, aci, func(buffer audio.IntBuffer, n int) error {
		samples := buffer.AsFloat32Buffer().Data
		C.Decoder_WriteSamples(d.decoder, (*C.float)(unsafe.Pointer(&samples[0])), C.size_t(n), C.int(0), &cErr)
		if cErr.error_code != C.uint(0) {
			return fmt.Errorf("failed writing samples to decoder, error code %d", cErr.error_code)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	C.Decoder_WriteSamples(d.decoder, nil, C.size_t(0), C.int(1), &cErr)
	if cErr.error_code != C.uint(0) {
		return "", fmt.Errorf("failed writing samples to decoder, error code %d", cErr.error_code)
	}

	var words []string
	for {
		res := C.Decoder_WaitResults(d.decoder, &cErr)
		if cErr.error_code != C.uint(0) {
			return "", fmt.Errorf("failed reading transcript from decoder, error code %d", cErr.error_code)
		}
		word := C.GoString(res.word)
		C.CResultWord_Destroy(res)

		if word == "" {
			break
		}
		word = strings.ToLower(word);
		words = append(words, word)
	}
	return strings.Join(words, " "), nil
}

type decoderFactory struct {
	factory *C.DecoderFactoryHandle
	bfr     []byte
}

func NewDecoderFactory(bundlePath string) (*decoderFactory, error) {
	model0 := C.CString(mustGetModelFile(bundlePath))
	defer func() {
		C.free(unsafe.Pointer(model0))
	}()

	bfr, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, err
	}

	cErr := C.DecoderError{}
	f := C.DecoderFactory_CreateFromModelArchive(unsafe.Pointer(&bfr[0]), C.ulong(len(bfr)), &cErr)
	if cErr.error_code != C.uint(0) {
		return nil, fmt.Errorf("failed to load on-device model bundle, error code %d", cErr.error_code)
	}
	return &decoderFactory{
		factory: f,
		bfr:     bfr,
	}, nil
}

func mustGetModelFile(file string) string {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic(fmt.Sprintf("required model file does not exist: %s", file))
	}
	return file
}

type cDecoder struct {
	decoder *C.DecoderHandle
	index   int
}

func (d *decoderFactory) NewStream(deviceID string) (*cDecoder, error) {
	cDeviceID := C.CString(deviceID)
	cErr := C.DecoderError{}
	decoder := C.DecoderFactory_GetDecoder(d.factory, cDeviceID, &cErr)
	if cErr.error_code != C.uint(0) {
		return nil, fmt.Errorf("failed creating decoder instance, error code %d", cErr.error_code)
	}
	defer C.free(unsafe.Pointer(cDeviceID))
	C.Decoder_SetParamI(decoder, C.SPEECHLY_DECODER_BLOCK_MULTIPLIER_I, 6, &cErr);
	return &cDecoder{
		decoder: decoder,
	}, nil
}
