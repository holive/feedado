package s3

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTempFileBuffer(t *testing.T) {
	b, err := NewTempFileBuffer()
	if err != nil {
		t.Fatal(err)
	}

	letters := "abcdefghijklm"

	t.Run("WriteAt", func(t *testing.T) {

		for i, letter := range letters {
			if _, err = b.WriteAt([]byte{byte(letter)}, int64(i)); err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("Read", func(t *testing.T) {
		var bt []byte
		bt, err = ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		if string(bt) != letters {
			t.Errorf("ReadAll = %v, expected %v", string(bt), letters)
		}
	})

	t.Run("Close and remove file", func(t *testing.T) {
		if err = b.Close(); err != nil {
			t.Fatal(err)
		}

		if _, err = os.Stat(b.tmp.Name()); !os.IsNotExist(err) {
			t.Fatal(err)
		}
	})
}
