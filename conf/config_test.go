package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Initial(t *testing.T) {
	Initial("")
	conf1 := Config

	Initial("C:/Users/dell-20/go/src/github.com/wangff15386/supermarket-go/conf/app.conf")
	conf2 := Config
	assert.Equal(t, conf1, conf2)

	Initial("../conf/app.conf")
	conf3 := Config
	assert.Equal(t, conf2, conf3)

	Initial("/32535y46eutiryukfgjyfhtg")
	conf4 := Config
	assert.Equal(t, conf1, conf4)
}
