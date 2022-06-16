//go:build on_device
// +build on_device

package cmd

/*
 #cgo CFLAGS: -I${SRCDIR}/../decoder/include
 #cgo darwin LDFLAGS: -L${SRCDIR}/../decoder/lib -Wl,-rpath,decoder/lib -lspeechly -framework Foundation -lc++ -O3
 #cgo linux LDFLAGS: -L${SRCDIR}/../decoder/lib -Wl,-rpath,$ORIGIN/../decoder/lib -Wl,--start-group -lstdc++ -lpthread -ldl -lm -lspeechly
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

func transcribeOnDevice(models []string, appID string, corpusPath string) error {
	ac := readAudioCorpus(corpusPath)
	df := NewDecoderFactory(models)
	for _, aci := range ac {
		d, err := df.NewStream(appID, "deviceID")
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
	readAudio(audioFilePath, aci, func(buffer audio.IntBuffer, n int) error {
		samples := buffer.AsFloat32Buffer().Data
		C.Decoder_WriteSamples(d.decoder, (*C.float)(unsafe.Pointer(&samples[0])), C.size_t(n), C.int(0))
		return nil
	})

	C.Decoder_WriteSamples(d.decoder, nil, C.size_t(0), C.int(1))

	var words []string
	for {
		res := C.Decoder_WaitResults(d.decoder)
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
}

func NewDecoderFactory(modelPaths []string) *decoderFactory {
	model0 := C.CString(mustGetModelFile(modelPaths[0]))
	model1 := C.CString(mustGetModelFile(modelPaths[1]))
	model2 := C.CString(mustGetModelFile(modelPaths[2]))
	model3 := C.CString(mustGetModelFile(modelPaths[3]))
	model4 := C.CString(mustGetModelFile(modelPaths[4]))
	model5 := C.CString(mustGetModelFile(modelPaths[5]))
	defer func() {
		C.free(unsafe.Pointer(model0))
		C.free(unsafe.Pointer(model1))
		C.free(unsafe.Pointer(model2))
		C.free(unsafe.Pointer(model3))
		C.free(unsafe.Pointer(model4))
	}()

	return &decoderFactory{
		factory: C.DecoderFactory_Create(model0, model1, model2, model3, model4, model5),
	}
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

func (d *decoderFactory) NewStream(appID string, deviceID string) (*cDecoder, error) {
	cAppID := C.CString(appID)
	cDeviceID := C.CString(deviceID)
	decoder := C.DecoderFactory_GetDecoder(d.factory, cAppID, cDeviceID)
	defer C.free(unsafe.Pointer(cAppID))
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
