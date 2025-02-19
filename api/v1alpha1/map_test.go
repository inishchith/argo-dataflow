package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_getContainer(t *testing.T) {
	const x Map = "my-expr"
	c := x.getContainer(getContainerReq{})
	assert.Equal(t, []string{"map", "my-expr"}, c.Args)
}
