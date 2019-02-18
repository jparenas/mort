package engine

import (
	"github.com/aldor007/mort/pkg/config"
	"github.com/aldor007/mort/pkg/object"
	"github.com/aldor007/mort/pkg/response"
	"github.com/aldor007/mort/pkg/transforms"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestImageEngine_Process_Error(t *testing.T) {
	image := response.NewNoContent(500)
	mortConfig := config.Config{}
	mortConfig.Load("testdata/config.yml")
	obj, err := object.NewFileObjectFromPath("/local/small/parent.jpg", &mortConfig)

	assert.Nil(t, err)
	assert.NotNil(t, obj)

	e := NewImageEngine(image)
	res, err := e.Process(obj, []transforms.Transforms{obj.Transforms})

	assert.NotNil(t, err)
	assert.Equal(t, res.StatusCode, 500)
}

func TestImageEngine_Process(t *testing.T) {
	f, err := os.Open("testdata/small.jpg")
	if err != nil {
		panic(err)
	}

	image := response.New(200, f)
	mortConfig := config.Config{}
	mortConfig.Load("testdata/config.yml")
	obj, err := object.NewFileObjectFromPath("/local/small/parent.jpg", &mortConfig)

	assert.Nil(t, err)
	assert.NotNil(t, obj)

	e := NewImageEngine(image)
	res, err := e.Process(obj, []transforms.Transforms{obj.Transforms})

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, res.Headers.Get("content-type"), "image/jpeg")
	// FIXME: falling on travis
	//assert.Equal(t, res.Headers.Get("etag"), "3ddbe830b458045b041d7653fc42970e")
	assert.Equal(t, res.Headers.Get("x-amz-meta-public-width"), "300")
}
