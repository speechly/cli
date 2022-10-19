//go:build on_device
// +build on_device

package cmd

/*
 #cgo CFLAGS: -I${SRCDIR}/../decoder/include
 #cgo darwin LDFLAGS: -L${SRCDIR}/../decoder/macos-x86_64/lib -Wl,-rpath,decoder/macos-x86_64/lib -lspeechly -lz -framework Foundation -lc++ -framework Security
 #cgo linux LDFLAGS: -L${SRCDIR}/../decoder/linux-x86_64/lib -Wl,-rpath,$ORIGIN/../decoder/linux-x86_64/lib -Wl,--start-group -lstdc++ -lpthread -ldl -lm -lspeechly -lz
 #cgo tflite LDFLAGS: -ltensorflowlite_c
 #cgo coreml LDFLAGS: -framework coreml
 #include <Decoder.h>
 #include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"unsafe"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func transcribeOnDevice(model string, corpusPath string) error {
	ac := readAudioCorpus(corpusPath)
	df, err := NewDecoderFactory(model)
	if err != nil {
		return err
	}
	for _, aci := range ac {
		d, err := df.NewStream("")
		if err != nil {
			return err
		}

		audioFilePath := path.Join(path.Dir(corpusPath), aci.Audio)
		transcript, err := decodeAudioCorpusItem(audioFilePath, aci, d)
		if err != nil {
			return err
		}

		res := &AudioCorpusItem{Audio: aci.Audio, Hypothesis: transcript}
		b, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}

	return nil
}

func decodeAudioCorpusItem(audioFilePath string, aci AudioCorpusItem, d *cDecoder) (string, error) {
	cErr := C.DecoderError{}

	readAudio(audioFilePath, aci, func(buffer audio.IntBuffer, n int) error {
		samples := buffer.AsFloat32Buffer().Data
		C.Decoder_WriteSamples(d.decoder, (*C.float)(unsafe.Pointer(&samples[0])), C.size_t(n), C.int(0), &cErr)
		if cErr.error_code != C.uint(0) {
			return fmt.Errorf("failed writing samples to decoder, error code %d", cErr.error_code)
		}
		return nil
	})

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
	return &cDecoder{
		decoder: decoder,
	}, nil
}

func readAudioCorpus(filename string) []AudioCorpusItem {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("An error occured: %s", err)
	}
	ac := make([]AudioCorpusItem, 0)

	jd := json.NewDecoder(f)
	for {
		var aci AudioCorpusItem
		err := jd.Decode(&aci)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Unmarshaling JSON failed: %s\n", err)
		}
		ac = append(ac, aci)
	}
	return ac
}

func readAudio(audioFilePath string, acItem AudioCorpusItem, callback func(buffer audio.IntBuffer, n int) error) {
	file, err := os.Open(audioFilePath)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalf("Closing file failed: %s\n", err)
		}
	}()
	if err != nil {
		log.Fatalf("Reading audio file failed: %s\n", err)
	}

	ad := wav.NewDecoder(file)
	ad.ReadInfo()
	if !ad.IsValidFile() {
		log.Fatalf("The audio file is not valid.\n")
	}

	afmt := ad.Format()

	if afmt.NumChannels != 1 || afmt.SampleRate != 16000 || ad.BitDepth != 16 {
		log.Fatalf("Only audio with 1ch 16kHz 16bit PCM wav files are supported. The audio file is %dch %dHz %dbit.\n",
			afmt.NumChannels, afmt.SampleRate, ad.BitDepth)
	}

	for {
		bfr := audio.IntBuffer{
			Format:         afmt,
			Data:           make([]int, 2048),
			SourceBitDepth: int(ad.BitDepth),
		}
		n, err := ad.PCMBuffer(&bfr)
		if err != nil {
			log.Fatalf("Reading audio file failed: %s\n", err)
		}

		if n == 0 {
			break
		}

		err = callback(bfr, n)
		if err != nil {
			log.Fatalf("Processing audio failed: %s\n", err)
		}
	}
}
