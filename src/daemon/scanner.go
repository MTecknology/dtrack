package daemon

import (
	// DTrack
	. "dtrack/common"
	"dtrack/state"
	// Standard
	"context"
	"strconv"
)

// Segment of WAV data
type audio_segment struct {
	count	uint64
	data	[]byte
}

// State machine for scanners that review audio chuunks
type segment_scanner struct {
	name	string
	ctx	context.Context
	cancel	context.CancelFunc
	segment	chan *audio_segment
}

// Start segment scanner thread for each trained model
func initialize_scanners(scanners []*segment_scanner, parent context.Context) {
	for i, model := range state.Runtime.Record_Inspect_Models {
		p, c := context.WithCancel(parent)
		new_scanner := &segment_scanner{
			name:		model,
			ctx:		p,
			cancel:		c,
			segment:	make(
					  chan *audio_segment,
					  state.Runtime.Record_Inspect_Backlog),
		}
		scanners[i] = new_scanner
		go new_scanner.scan_segments()
	}
}

// Primary loop that tests each audio segment against a trained model
func (stream *segment_scanner) scan_segments() {
	//frame := make([]audio_segment, 2)
	for {
		select {
		case incoming_segment, ok := <-stream.segment:
			if !ok {
				Warn("Worker unexpectedly closed: %s", stream.name)
				return
			}
			Debug("%s is processing segment %s", stream.name,
				strconv.FormatUint(incoming_segment.count, 10))
			// TODO: frame[1]->frame[0]; segment->frame[1]; inference()
		}
	}
}
