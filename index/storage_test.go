package index

import (
	"bufio"
	"github.com/alldroll/suggest/alphabet"
	"github.com/alldroll/suggest/compression"
	"github.com/alldroll/suggest/dictionary"
	"log"
	"os"
	"testing"
)

func TestOnDiscWriter_Save(t *testing.T) {
	header, err := os.Create("../testdata/db/test.hd")
	defer func() {
		header.Close()
		os.Remove(header.Name())
	}()
	if err != nil {
		log.Fatal(err)
	}

	docList, err := os.Create("../testdata/db/test.dl")
	defer func() {
		docList.Close()
		os.Remove(docList.Name())
	}()

	if err != nil {
		log.Fatal(err)
	}

	dict, err := os.Open("../testdata/cars.dict")
	defer dict.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(dict)
	collection := make([]string, 0)
	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	indices := buildIndex()
	writer := NewOnDiscIndicesWriter(compression.VBEncoder(), header, docList, 0)
	err = writer.Save(indices)
	if err != nil {
		t.Error(err)
	}

	reader := NewOnDiscIndicesReader(compression.VBDecoder(), header, docList, 0)
	_, err = reader.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func buildIndex() Indices {
	dict, err := os.Open("../testdata/cars.dict")
	defer dict.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(dict)
	collection := make([]string, 0)
	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	dictionary := dictionary.NewInMemoryDictionary(collection)

	alphabet := alphabet.NewCompositeAlphabet([]alphabet.Alphabet{
		alphabet.NewRussianAlphabet(),
		alphabet.NewEnglishAlphabet(),
		alphabet.NewNumberAlphabet(),
		alphabet.NewSimpleAlphabet([]rune{'$'}),
	})

	indexer := NewIndexer(3, NewGenerator(3, alphabet), NewCleaner(alphabet.Chars(), "$", [2]string{"$", "$"}))

	return indexer.Index(dictionary)
}