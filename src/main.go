//##
// Audio Disturbance Tracker (DTrack)
//
// License:   AGPL-3
// Copyright: 2024-2026, Michael Lustfield (MTecknology)
// Authors:   See history with "git log" or "git blame"
//##
package main

import (
	// Bootstrap
	"dtrack/common"
	"dtrack/state"
	// Actions
	"dtrack/ai"
	"dtrack/daemon"
	"dtrack/review"
)

func main() {
	// Bootstrap
	parse_flags()
	common.Debug_Enabled = *app_verbose
	state.Load_Configuration(*app_config_path)

	// Kickoff
	action_map := map[string]func() {
		"train":  ai.Train,
		"monitor": daemon.Start,
		"review":  review.Start,
	}
	action_map[*app_action]()

	// Post-processing
	common.Clean_Workspace()
}
