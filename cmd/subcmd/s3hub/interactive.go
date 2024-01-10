package s3hub

import tui "github.com/nao1215/rainbow/ui/s3hub"

// interactive starts s3hub command interactive UI.
func interactive() error {
	return tui.RunS3hubUI()
}
