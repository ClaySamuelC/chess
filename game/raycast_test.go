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
			name: "should return false when raycast hits unit of wrong team color first",
			args: args{11, 0, 1, "White", &map[string]bool{"Rook": true}},
			want: false,
		},
		{
			fen:  "8/3p4/8/8/8/3B4/3Q4/8 w - - 0 1",
			name: "should return false when raycast hits unit of wrong rank first",
			args: args{11, 0, 1, "White", &map[string]bool{"Queen": true, "Rook": true}},
			want: false,
		},
		{
			fen:  "p7/8/8/8/8/8/8/7B w - - 0 1",
			name: "should return true when raycast has correct target in path",
			args: args{0, 1, 1, "White", &map[string]bool{"Queen": true, "Bishop": true}},
			want: true,
		},
		{
			fen:  "Pp6/8/8/8/8/8/8/8 w - - 0 1",
			name: "should return true",
			args: args{0, 1, 0, "Black", &map[string]bool{"Pawn": true}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := CreateGame(tt.fen)
			got := c.RayCast(tt.args.pos, tt.args.xDir, tt.args.yDir, tt.args.targetTeam, tt.args.targetRanks)
			assert.Equal(t, tt.want, got)
		})
	}
}
