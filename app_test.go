package alkaid

import (
	"fmt"
	"testing"
)

func TestNewApp(t *testing.T) {
	_, err := NewApp()
	fmt.Println(err)
}
