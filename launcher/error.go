package launcher

import (
	"fmt"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

func (l *Launcher) onError(prefix, description string, e error) {
	message := fmt.Sprintf("[%s] %s:", prefix, description)
	errs.Log(message, e)
	tasks.Shutdown()
}
