package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundImages(t *testing.T) {
	// prepare
	var actual, expected []uint
	ids := []uint{9, 1, 2, 3, 4, 5, 6}
	images := []ImageFile{
		{ImageID: 1},
		{ImageID: 2},
		{ImageID: 4},
		{ImageID: 6},
	}

	// scenario: success
	actual = notFoundImages(ids, images)
	expected = []uint{9, 3, 5}
	assert.EqualValues(t, expected, actual)

	// scenario: fail
	actual = notFoundImages(ids, images)
	expected = []uint{}
	assert.NotEqual(t, expected, actual)
	actual = notFoundImages(ids, images)
	expected = []uint{9, 1, 2, 3, 4, 5, 6}
	assert.NotEqual(t, expected, actual)
}
