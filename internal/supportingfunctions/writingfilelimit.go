package supportingfunctions

import (
	"os"
	"path"
)

// WritingFileLimit выполняет запись, заданной в килобайтах, части исходного файла
// в новый файл находящийся в этой же директории и задает такое же имя файла с добавлением
// в конце файла дополнительного суфикса, например '.limit'
func WritingFileLimit(pathName, fileName, suffix string, maxChunk int) (int64, error) {
	/*f, err := os.Open(path.Join(pathName, fileName))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fdst, err := os.OpenFile(path.Join(pathName, fileName+suffix), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return 0, err
	}
	defer fdst.Close()

	r := io.LimitReader(f, int64(maxChunk*1024*1024))
	num, err := io.Copy(fdst, r)
	if err != nil {
		return 0, err
	}*/

	chunk := int64(maxChunk * 1024 * 1024)
	oldPath := path.Join(pathName, fileName)

	if err := os.Truncate(oldPath, chunk); err != nil {
		return chunk, err
	}

	os.Rename(oldPath, path.Join(pathName, fileName+suffix))

	return chunk, nil
}
