package daemon

import (
	// DTrack
	"dtrack/log"
	//"dtrack/state"
)

// Segment of WAV data
type audio_segment struct {
	count uint
	data  []byte
}

// Primary loop that tests each audio segment against a trained model
func scan_segments(name string, segment chan audio_segment) {
	var last_segment audio_segment
	for {
		// Wait for incoming segments
		incoming_segment, ok := <-segment
		if !ok {
			log.Die("Scanner unexpectedly closed: %s", name)
			return
		}

		// Save first segment seen, but delay processing until next segment
		if last_segment.data == nil {
			last_segment = incoming_segment
			continue
		}

		// Forward to ML process and test for match
		log.Trace("Model %s is scanning %d", name, last_segment.count)
		//check_window := append(last_segment.data, incoming_segment.data...)
		//TODO

		// Swap last-seen segment and resume waiting for next segment
		last_segment = incoming_segment
	}
}
