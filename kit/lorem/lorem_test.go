package lorem_test

import (
	"testing"

	"github.com/sempernow/uqc/kit/lorem"
)

func TestLorem(t *testing.T) {

	t.Skip()

	t.Log("Host():", lorem.Host())
	t.Log("URL():", lorem.URL())
	t.Log("Email():", lorem.Email())
	t.Log("Word(5,8): ", lorem.Word(5, 8))
	t.Log("Sentence(5,8):", lorem.Sentence(5, 8))
	t.Log("Paragraph(5,8):", lorem.Paragraph(5, 8))
}
