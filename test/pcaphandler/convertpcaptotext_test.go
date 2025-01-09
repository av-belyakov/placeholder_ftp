package pcaphandler_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_ftp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
)

func TestConvertPcapToText(t *testing.T) {
	var (
		filePath     string = "../test_files/"
		readFileName string = "test_pcap_file.pcap"
		//readFileName string = "1616152425_2021_03_19____14_13_45_24636.tdp"
		writeFileName string = "test_pcap_file.pcap.txt"
		//writeFileName string = "1616152425_2021_03_19____14_13_45_24636.tdp.txt"
	)

	logging := logginghandler.New()
	go func() {
		for msgErr := range logging.GetChan() {
			t.Log("ERROR:", msgErr)
		}
	}()

	// для файла по которому выполняется декодирование пакетов
	readFile, err := os.Open(path.Join(filePath, readFileName))
	assert.NoError(t, err)

	// для файла в который выполняется запись информации полученной в результате декодирования
	writeFile, err := os.OpenFile(path.Join(filePath, writeFileName), os.O_RDWR|os.O_CREATE, 0666)
	assert.NoError(t, err)

	err = supportingfunctions.NetworkTrafficDecoder(readFileName, readFile, writeFile, logging)
	assert.NoError(t, err)

	readFile.Close()
	writeFile.Close()
}

/*func TestConvertPcapToTextTwo(t *testing.T) {
	var (
		filePath string = "../test_files/"
		//readFileName string = "test_pcap_file.pcap"
		readFileName string = "1616152425_2021_03_19____14_13_45_24636.tdp"
		//writeFileName string = "___test_pcap_file.pcap.txt"
		writeFileName string = "1616152425_2021_03_19____14_13_45_24636.tdp.txt"
	)

	logging := logginghandler.New()
	go func() {
		for msgErr := range logging.GetChan() {
			fmt.Println("ERROR:", msgErr)
		}
	}()

	// для файла в который выполняется запись информации полученной в результате декодирования
	writeFile, err := os.OpenFile(path.Join(filePath, writeFileName), os.O_RDWR|os.O_CREATE, 0666)
	assert.NoError(t, err)

	err = supportingfunctions.NetTraffPcapDecoder(filePath, readFileName, writeFile, logging)
	assert.NoError(t, err)

	writeFile.Close()
}*/
