package s3hub

import "github.com/nao1215/rainbow/ui"

// interactive starts s3hub command interactive UI.
func interactive() error {
	return ui.RunS3hubUI()
}
