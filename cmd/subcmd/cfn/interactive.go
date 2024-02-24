package cfn

import tui "github.com/nao1215/rainbow/ui/cfn"

// interactive starts cfn command interactive UI.
func interactive() error {
	return tui.RunCfnUI()
}
