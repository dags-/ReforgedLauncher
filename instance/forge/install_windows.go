package forge

func buildCommand(classpath, gameDir string) *exec.Cmd {
	cmd := exec.Command("java", "-classpath", classpath, "Main", i.GameDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
