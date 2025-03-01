package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/env"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubitect_Empty(t *testing.T) {
	k := Kubitect{
		Url:     URL(env.ConstProjectUrl),
		Version: MasterVersion("master"),
	}

	assert.NoError(t, k.Validate())
	assert.NoError(t, Kubitect{}.Validate())
}
