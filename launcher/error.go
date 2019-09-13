package launcher

import (
	"fmt"
	"os"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

func (l *Launcher) onError(prefix, description string, e error) {
	message := fmt.Sprintf("[%s] %s:", prefix, description)
	errs.Log(message, e)
	os.Exit(1)
}
