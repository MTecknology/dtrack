//##
// DTrack Package: Surveilance Monitor
//
// Collects audio+video files and logs any matched audio disturbances.
//##
package daemon

import (
	// DTrack
	. "dtrack/common"
	"dtrack/state"
	// Standard
	"context"
	"time"
	"io"
)

// Primary post-bootstrap entry point
// Initialize audio segment scanners and begin recording process
func Start() {
	process, cancel := context.WithCancel(context.Background())
	defer cancel()
	wav_stream, daemon_stream := io.Pipe()

	// Start scanners if any models are defined
	if state.Runtime.Has_Models {
		Debug("Initializing segment scanners")
		start_scanners(process, wav_stream)
	} else {
		Warn("No inspection models configured; only able to record!")
		go Pipe2DevNull(wav_stream)
	}

	// Start main recording loop that sends data to scanners
	for {
		run_ffmpeg(process, daemon_stream)
		// Pause to prevent thrashing of physical devices
		time.Sleep(50 * time.Millisecond)
	}

}

// Initialize all audio segment scanners and process wav_stream data
func start_scanners(process context.Context, wav_stream *io.PipeReader) {
	// Process manager for segment scanners
	scanners := make(
		[]*segment_scanner,
		len(state.Runtime.Record_Inspect_Models))
	initialize_scanners(scanners, process)
	returned_segments := make(chan *audio_segment)

	// Stream converter
	go stream_to_segment(wav_stream, returned_segments)

	// Distribute new segments to all scanners
	go func() {
	for {
		select {
		// Collect new segment
		case new_segment, ok := <-returned_segments:
			if !ok {
				Die("Stream converter disappeared")
				return
			}
			// Distribute segment to scanners
			for _, w := range scanners {
				select {
				// Send segment to individual scanner
				case w.segment <- new_segment:
				default:
					Warn("Scanner Blocked: %s", w.name)
				}
			}
		}
	}}()
}

// Convert an input wav_stream to 1-second audio clips
func stream_to_segment(stream *io.PipeReader, segments chan<- *audio_segment) {
	defer stream.Close()
	defer close(segments)
	var segment_id uint64 = 0
	if state.Runtime.Record_Inspect_Segment <= 0 {
		Die("MISSING: record_inspect_segment ; See README.rst")
	}

	/* TODO - Replace state.Runtime.Record_Inspect_Segment w/ automatic calculation.
	// Wait until we can identify the byte size of a 1-second audio segment
	Debug("Identifying size of 1-second audio segment ...")
	decoder := wav.NewDecoder(stream)
	if err := decoder.ReadInfo(); err != nil {
		if err != io.EOF {
			Die("Received end of ffmpeg without finding header.")
		}
		return
	}

	// Calculate size of audio segment
	format := decoder.Format()
	if format == nil {
		Die("WAV decoder returned a nil Format")
	}
	bits := decoder.SampleBitDepth()
	// A typical 16-bit 44.1kHz stereo stream is 176,400 bytes per second
	segment_size = format.SampleRate * format.NumChannels * (int(bits) / 8)
	if segment_size <= 0 {
		Die("Unexpected size of audio segment")
	}
	Debug("WAV Header received; segment size set to %d b/s", segment_size)
	*/

	// Start main conversion loop
	for {
		// Allocate a buffer for the audio segment
		segment_data := make([]byte, state.Runtime.Record_Inspect_Segment)
		
		// Read segment_size bytes from the stream
		// io.ReadFull blocks until the buffer is full OR an error occurs
		_, err := io.ReadFull(stream, segment_data)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// WAV stream ended; restart fresh loop
			Debug("Reading segment was reset")
			continue
		}
		if err != nil {
			Die("Unhandled stream read error: %s", err.Error())
		}

		// Add new segment to queue
		segment_id++
		Debug("New segment accumulated: %d", segment_id)
		segments <- &audio_segment {
			count: segment_id,
			data: segment_data,
		}
	}
}
