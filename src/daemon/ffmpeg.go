package daemon

import (
	// DTrack
	. "dtrack/common"
	"dtrack/state"
	// Standard
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

// MKV Filename:  YYYY-MM-DD_HHmmss
const mkv_name = "2006-01-02_150405.mkv"

// Maintain copy of accumulated arguments
var session_arguments []string
var build_ffargs sync.Once

// Run ffmpeg command, saving A/V to MKV and Audio-only to Stream
func run_ffmpeg(parent context.Context, stdout *io.PipeWriter) {
	// defer stdout.Close()  (Do not close stream)
	// Build ffmpeg command
	build_ffargs.Do(func() {
		session_arguments = ffmpeg_arguments()
	})
	args := session_arguments
	mkv := state.Runtime.Workspace + "/recordings/" + time.Now().Format(mkv_name)
	args = append(args, mkv)
	ffmpeg := exec.CommandContext(parent, "ffmpeg", args...)
	ffmpeg.Stderr = os.Stderr
	ffmpeg.Stdout = stdout

	// Start ffmpeg process
	Debug("New ffmpeg process, saving to: %s", mkv)
	if ffmpeg.Start() != nil {
		Die("Failed to intialize ffmpeg")
	}
	if ffmpeg.Wait() != nil {
		Warn("ffmpeg finished with errors")
		// Extra pause for potential device thrashing
		time.Sleep(1 * time.Second)
	}
}

// Return the string for an ffmpeg command with the pattern:
// ffmpeg [basic-options] \
//   [audio-options] [audio-device] \
//   [video-options] [video-device] \
//   [output-wav] [to-stdout] \
//   [output-wav&vid] [to-mkv] [MISSING:filename]
func ffmpeg_arguments() []string {
	// 5+2+_+2+2+_+2+7+_+10 = 30 (+filename +vars)
	arg_count := 30 + 1 +
		len(state.Runtime.Record_Audio_Options) +
		len(state.Runtime.Record_Video_Options) +
		len(state.Runtime.Record_Video_Advanced)
	if !state.Runtime.Has_Models {
		arg_count -= 7
	}
	cmd := make([]string, 0, arg_count)

	// basic-options  +5
	cmd = append(cmd, "-y", "-loglevel", "fatal", "-nostdin", "-nostats")
	
	// audio-options  +2 +X
	cmd = append(cmd, "-t", state.Runtime.Record_Duration)
	cmd = append(cmd, state.Runtime.Record_Audio_Options...)
	// audio-device   +2
	cmd = append(cmd, "-i", state.Runtime.Record_Audio_Device)

	// video-options  +2 +X
	cmd = append(cmd, "-t", state.Runtime.Record_Duration)
	cmd = append(cmd, state.Runtime.Record_Video_Options...)
	// video-device   +2
	cmd = append(cmd, "-i", state.Runtime.Record_Video_Device)

	// wav-to-stdout  +7
	if state.Runtime.Has_Models {
		cmd = append(cmd, "-map", "0:a", "-c:a", "pcm_s16le", "-f", "wav", "-")
	}
	// wav&vid-to-mkv +X +10
	cmd = append(cmd, state.Runtime.Record_Video_Advanced...)
	cmd = append(cmd, "-map", "0:a", "-map", "[dtstamp]")
	cmd = append(cmd, "-c:a", "pcm_s16le", "-c:v", "libx264")
	cmd = append(cmd, "-preset", state.Runtime.Record_Compression)

	return cmd
}
