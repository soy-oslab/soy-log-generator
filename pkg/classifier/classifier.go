package classifier

import (
	"log"
	"os"
	"strings"

	s "github.com/bbalet/stopwords"
	b "github.com/jbrukh/bayesian"
)

const (
	// Hot indicates the hot data
	Hot b.Class = "Hot"
	// Cold indicates the cold data
	Cold b.Class = "Cold"
)

// Classifier contains the information for identify the hot or not
type Classifier struct {
	classifier *b.Classifier
	filepath   string
	classes    [2]b.Class
}

// InitClassfier initializes and allocates the Classifier structure
func InitClassfier(savepath string) (*Classifier, error) {
	var err error
	c := new(Classifier)
	stat, err := os.Stat(savepath)
	if stat != nil {
		log.Println(os.IsNotExist(err), stat.Size() > int64(0))
	}
	if !os.IsNotExist(err) && stat.Size() > int64(0) {
		log.Println("this")
		c.classifier, err = b.NewClassifierFromFile(savepath)
	} else {
		log.Println("that")
		c.classifier = b.NewClassifier(Hot, Cold)
	}
	c.filepath = savepath
	c.classes[0] = Hot
	c.classes[1] = Cold
	return c, err
}

// preprocessing preprocesses the string
func preprocessing(str string) []string {
	str = strings.ToLower(str)
	str = strings.Trim(str, " \t\n")
	str = s.CleanString(str, "en", false)
	slice := strings.Split(str, " ")
	return slice
}

// learn learns the string based on bayesian
func (c *Classifier) learn(str string, state b.Class) {
	slice := preprocessing(str)
	c.classifier.Learn(slice, state)
}

// Backup stores the model
func (c *Classifier) Backup() error {
	err := c.classifier.WriteToFile(c.filepath)
	return err
}

// classfy identifies this string is hot or not
func (c *Classifier) classify(str string) (map[b.Class]float64, b.Class) {
	slice := preprocessing(str)
	scores, likely, _ := c.classifier.ProbScores(slice)
	class := c.classes[likely]
	results := make(map[b.Class]float64)
	results[Hot] = scores[0]
	results[Cold] = scores[1]
	return results, class
}
