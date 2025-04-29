package recorder

import (
	"fmt"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Recorder struct {
}

func New() *Recorder {
	return &Recorder{}
}

// ffmpeg -y -i <input> -t <seconds> -c:v copy -c:a copy -bsf:a aac_adtstoasc <output>
func (r *Recorder) StartRecording(playlistURL string, durationSeconds int, outputPath string) error {
	input := ffmpeg.Input(playlistURL)

	args := map[string]interface{}{
		"t":     strconv.Itoa(durationSeconds),
		"c:v":   "copy",
		"c:a":   "copy",
		"bsf:a": "aac_adtstoasc",
	}

	err := input.Output(outputPath, ffmpeg.KwArgs(args)).Run(ffmpeg.SeparateProcessGroup())
	if err != nil {
		return fmt.Errorf("ffmpeg record failed: %w", err)
	}
	return nil
}
