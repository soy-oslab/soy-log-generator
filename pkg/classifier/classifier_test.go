package classifier

import (
	"log"
	"os"
	"testing"
)

func TestClassifier(t *testing.T) {
	file, err := os.CreateTemp("", "model.sav.")
	if err != nil {
		t.Errorf("create temp file failed")
	}
	defer file.Close()
	defer os.Remove(file.Name())
	_, err = InitClassfier("")
	if err != nil {
		t.Errorf("expected no error but error occurred")
	}
	c1, err := InitClassfier(file.Name())
	if err != nil {
		t.Errorf("init classification failed")
	}
	c1.Learn("Hello World This is the Test", Hot)
	c1.Learn("Hate this isn't test", Cold)
	test, likely := c1.Classify("Hat this test")
	log.Println(test, likely)
	err = c1.Backup()
	if err != nil {
		t.Errorf("backup failed")
	}
	c2, err := InitClassfier(file.Name())
	if err != nil {
		t.Errorf("load exist model failed")
	}
	test, likely = c2.Classify("Hello this test")
	log.Println(test, likely)
}
