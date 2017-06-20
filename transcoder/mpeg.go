package transcoder

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type MpegTs struct {
	cmd     *exec.Cmd
	tmpFile string
}

func NewMpegTs(input io.Reader, output io.Writer) (*MpegTs, error) {
	var tc MpegTs
	tmpFile, err := ioutil.TempFile("", "StreamService")
	if err != nil {
		return nil, err
	}
	tc.tmpFile = tmpFile.Name()
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, input)
	if err != nil {
		return nil, err
	}
	tc.init(tc.tmpFile, output)
	return &tc, nil
}

func NewMpegTsFile(filename string, output io.Writer) *MpegTs {
	var tc MpegTs
	tc.init(filename, output)
	return &tc
}

func (tc *MpegTs) init(filename string, output io.Writer) {
	tc.cmd = exec.Command("ffmpeg", "-i", filename, "-map_metadata", "-1", "-c", "copy", "-bsf", "h264_mp4toannexb", "-f", "mpegts", "-")
	//tc.cmd = exec.Command("ffmpeg", "-i", filename, "-c", "copy", "-f", "mpegts", "-")
	//tc.cmd = exec.Command("ffmpeg", "-i", filename, "-vcodec", "libx264", "-acodec", "aac", "-profile:v", "baseline", "-f", "mpegts", "-")

	// ffmpeg -i demo.mp4 -codec copy -map 0 -f mpegts -vbsf h264_mp4toannexb -flags -global_header -segment_format mpegts -segment_list_flags +live out.ts
	//tc.cmd = exec.Command("ffmpeg", "-i", filename, "-codec", "copy", "-map", "0", "-f", "mpegts", "-vbsf", "h264_mp4toannexb", "-flags", "-global_header", "-segment_format", "mpegts", "-segment_list_flags", "-se", "+live", "-")
	tc.cmd.Stdout = output
}

func (tc *MpegTs) TempFileName() string {
	return tc.tmpFile
}

func (tc *MpegTs) Transcode() error {
	return tc.cmd.Run()
}

func (tc *MpegTs) Dispose() error {
	if tc.tmpFile != "" {
		return os.Remove(tc.tmpFile)
	}
	return nil
}
