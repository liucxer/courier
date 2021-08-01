package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVolumeString(t *testing.T) {
	tt := require.New(t)

	type Case struct {
		v      string
		volume Volume
	}

	cases := []Case{
		{
			v: "content",
			volume: Volume{
				Name:       "content",
				AccessMode: VolumeAccessModeReadWrite,
			},
		},
		{
			v: "content:ro",
			volume: Volume{
				Name:       "content",
				AccessMode: VolumeAccessModeReadOnly,
			},
		},
		{
			v: "content:/tmp",
			volume: Volume{
				Name:       "content",
				MountPath:  "/tmp",
				AccessMode: VolumeAccessModeReadWrite,
			},
		},
		{
			v: "/tmp",
			volume: Volume{
				MountPath: "/tmp",
			},
		},
		{
			v: "container:content",
			volume: Volume{
				Name:       "content",
				AccessMode: VolumeAccessModeReadWrite,
			},
		},
		{
			v: "content:/tmp:ro",
			volume: Volume{
				Name:       "content",
				MountPath:  "/tmp",
				AccessMode: VolumeAccessModeReadOnly,
			},
		},
		{
			v: "/tmp:/tmp:ro",
			volume: Volume{
				LocalPath:  "/tmp",
				MountPath:  "/tmp",
				AccessMode: VolumeAccessModeReadOnly,
			},
		},
	}

	for _, caseItem := range cases {
		volume, err := ParseVolumeString(caseItem.v)
		tt.Nil(err)
		tt.Equal(&caseItem.volume, volume)
	}
}
