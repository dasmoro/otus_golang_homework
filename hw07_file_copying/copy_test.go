package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("no to and from", func(t *testing.T) {
		err := Copy("", "", 0, 0, false)
		t.Log(err)
		require.Error(t, err)
	})
	t.Run("bad from good to", func(t *testing.T) {
		err := Copy("%$()#r3", "test.test", 0, 0, false)
		t.Log(err)
		require.Error(t, err)
	})
	t.Run("Regular Copying", func(t *testing.T) {
		f, err := os.Create("test1")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte("test string 29 symbols length"))
		f.Close()

		err = Copy("test1", "copied1", 0, 0, false)
		require.NoError(t, err)
		f, err = os.Open("copied1")
		if err != nil {
			t.Fatal(err)
		}
		fi, err := f.Stat()
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		os.Remove("test1")
		os.Remove("copied1")
		require.Equal(t, int64(29), fi.Size())

	})

	t.Run("Offset Copying", func(t *testing.T) {
		f, err := os.Create("test2")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte("test string 29 symbols length"))
		f.Close()

		err = Copy("test2", "copied2", 5, 0, false)
		require.NoError(t, err)
		f, err = os.Open("copied2")
		if err != nil {
			t.Fatal(err)
		}
		fi, err := f.Stat()
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		os.Remove("test2")
		os.Remove("copied2")
		require.Equal(t, int64(29-5), fi.Size())

	})

	t.Run("Offset Limit Copying", func(t *testing.T) {
		f, err := os.Create("test3")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte("test string 29 symbols length"))
		f.Close()

		err = Copy("test3", "copied3", 5, 29, false)
		require.NoError(t, err)
		f, err = os.Open("copied3")
		if err != nil {
			t.Fatal(err)
		}
		fi, err := f.Stat()
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		os.Remove("test3")
		os.Remove("copied3")
		require.Equal(t, int64(29-5), fi.Size())

	})

	t.Run("Offset out of bounds Copying", func(t *testing.T) {
		f, err := os.Create("test4")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte("test string 29 symbols length"))
		f.Close()
		err = Copy("test4", "copied4", 577, 0, false)
		t.Log(err)
		os.Remove("test4")
		require.Error(t, err)

	})

	t.Run("dir", func(t *testing.T) {
		err := Copy("testdata", "test.test", 0, 0, false)
		t.Log(err)
		require.Error(t, err)
	})

	t.Run("equals", func(t *testing.T) {
		err := Copy("copy.go", "copy.go", 0, 0, false)
		t.Log(err)
		require.Error(t, err)
	})

	t.Run("linux random filesize things", func(t *testing.T) {
		err := Copy("/dev/random", "test.test", 0, 0, false)
		t.Log(err)
		require.Error(t, err)
	})
}
