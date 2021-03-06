package entity

import (
	"errors"

	pb "github.com/domino14/liwords/rpc/api/proto/realtime"
)

// Variants, time controls, etc.

type Variant string
type TimeControl string

const (
	VarClassic   Variant = "classic"
	VarAWorth100         = "a-is-worth-100"
	VarDogworms          = "dogworms" // OMGWords scrambled = dogworms?
	VarSuper             = "superomg"
)

const (
	TCRegular    TimeControl = "regular"    // > 14/0
	TCRapid                  = "rapid"      // 6/0 to <= 14/0
	TCBlitz                  = "blitz"      // > 2/0 to < 6/0
	TCUltraBlitz             = "ultrablitz" // 2/0 and under
	TCCorres                 = "corres"
)

const (
	// Cutoffs in seconds for different time controls.
	CutoffUltraBlitz = 2 * 60
	CutoffBlitz      = 6 * 60
	CutoffRapid      = 14 * 60
)

// Calculate "total" time assuming there are 16 turns in a game per player.
const turnsPerGame = 16 // just an estimate.

// TotalTimeEstimate estimates the amount of time this game will take, per side.
func TotalTimeEstimate(gamereq *pb.GameRequest) int32 {
	return gamereq.InitialTimeSeconds +
		(gamereq.MaxOvertimeMinutes * 60) +
		(gamereq.IncrementSeconds * turnsPerGame)
}

func VariantFromGameReq(gamereq *pb.GameRequest) (TimeControl, Variant, error) {
	// hardcoded values here; fix sometime
	var timefmt TimeControl

	totalTime := TotalTimeEstimate(gamereq)

	if totalTime <= CutoffUltraBlitz {
		timefmt = TCUltraBlitz
	} else if totalTime <= CutoffBlitz {
		timefmt = TCBlitz
	} else if totalTime <= CutoffRapid {
		timefmt = TCRapid
	} else {
		timefmt = TCRegular
	}
	var variant Variant
	if gamereq.Rules.BoardLayoutName == CrosswordGame {
		variant = VarClassic
	} else {
		return "", "", errors.New("unsupported game type")
	}

	return timefmt, variant, nil
}
