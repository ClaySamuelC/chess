package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRaycast(t *testing.T) {
	type args struct {
		pos         int
		xDir        int
		yDir        int
		targetTeam  string
		targetRanks *map[string]bool
	}
	tests := []struct {
		fen  string
		name string
		args args
		want bool
	}{
		{
			fen:  "8/3p4/8/8/8/8/3P4/8 w HAha - 0 1",
			name: "should return false",
			args: args{11, 0, 1, "White", &map[string]bool{"Rook": true}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := CreateGame(tt.fen)
			got := c.rayCast(tt.args.pos, tt.args.xDir, tt.args.yDir, tt.args.targetTeam, tt.args.targetRanks)
			assert.Equal(t, tt.want, got)
		})
	}
}
