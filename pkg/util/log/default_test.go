package log

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T) {
	logger := SetDefault(nil)
	assert.NotNil(t, logger)

	logger = logrus.New()
	assert.Equal(t, logger, SetDefault(logger))
}
