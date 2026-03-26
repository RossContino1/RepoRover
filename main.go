package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const appVersion = "1.2.1"

var busy atomic.Bool

//go:embed assets/icon.png
var iconBytes []byte

type ManagerInfo struct {
	Name     string
	IsSystem bool
}

type OSInfo struct {
	ID         string
	PrettyName string
}

type AURAction int

const (
	AURActionSkip AURAction = iota
	AURActionOpenTerminal
)

type AURDecision struct {
	Proceed bool
	Action  AURAction
}

type TerminalLaunch struct {
	Name string
	Args []string
}

func main() {
	iconRes := fyne.NewStaticResource("icon.png", iconBytes)

	a := app.NewWithID("com.bytesbreadbbq.reporover")
	a.Settings().SetTheme(theme.LightTheme())
	a.SetIcon(iconRes)

	w := a.NewWindow("RepoRover")
	w.SetIcon(iconRes)
	w.Resize(fyne.NewSize(920, 700))

	osLabel := widget.NewLabel("Detected OS: not checked yet")
	pmLabel := widget.NewLabel("Primary package manager: not checked yet")
	rebootLabel := widget.NewLabel("Reboot status: not checked yet")

	includeFlatpak := widget.NewCheck("Include Flatpak", func(bool) {})
	includeFlatpak.SetChecked(true)

	includeSnap := widget.NewCheck("Include Snap", func(bool) {})
	includeSnap.SetChecked(true)

	dryRun := widget.NewCheck("Dry Run", func(bool) {})
	systemOnly := widget.NewCheck("System Only", func(bool) {})

	statusLabel := widget.NewLabel("Ready")
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)

	output := widget.NewRichText(
		&widget.TextSegment{Text: "RepoRover GUI ready.\n"},
	)
	output.Wrapping = fyne.TextWrapWord
	outputScroll := container.NewScroll(output)

	var plainLog strings.Builder
	plainLog.WriteString("RepoRover GUI ready.\n")

	appendOutput := func(text string) {
		fyne.Do(func() {
			line := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), text)
			output.Segments = append(output.Segments, &widget.TextSegment{Text: line + "\n"})
			output.Refresh()
			outputScroll.ScrollToBottom()
			plainLog.WriteString(line + "\n")
		})
	}

	clearOutput := func() {
		fyne.Do(func() {
			output.Segments = []widget.RichTextSegment{
				&widget.TextSegment{Text: "RepoRover GUI ready.\n"},
			}
			output.Refresh()
			plainLog.Reset()
			plainLog.WriteString("RepoRover GUI ready.\n")
		})
	}

	copyLog := func() {
		w.Clipboard().SetContent(plainLog.String())
		dialog.ShowInformation("Copy Log", "Log copied to clipboard.", w)
	}

	setHeader := func(osText, pmText, rebootText string) {
		fyne.Do(func() {
			osLabel.SetText(osText)
			pmLabel.SetText(pmText)
			rebootLabel.SetText(rebootText)
		})
	}

	setStatus := func(text string, value float64) {
		fyne.Do(func() {
			statusLabel.SetText(text)
			progressBar.SetValue(value)
		})
	}

	var checkButton *widget.Button
	var runButton *widget.Button

	setButtonsEnabled := func(enabled bool) {
		fyne.Do(func() {
			if enabled {
				checkButton.Enable()
				runButton.Enable()
			} else {
				checkButton.Disable()
				runButton.Disable()
			}
		})
	}

	startOperation := func(name string) bool {
		if !busy.CompareAndSwap(false, true) {
			appendOutput("Another operation is already running.")
			setStatus("Busy: "+name+" blocked", progressBar.Value)
			return false
		}
		setButtonsEnabled(false)
		return true
	}

	endOperation := func() {
		busy.Store(false)
		setButtonsEnabled(true)
	}

	openLink := func(raw string) {
		u, err := url.Parse(raw)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err := a.OpenURL(u); err != nil {
			dialog.ShowError(err, w)
		}
	}

	showAbout := func() {
		aboutWin := a.NewWindow("About RepoRover")
		aboutWin.Resize(fyne.NewSize(560, 430))
		aboutWin.SetIcon(iconRes)

		icon := canvas.NewImageFromResource(iconRes)
		icon.SetMinSize(fyne.NewSize(64, 64))
		icon.FillMode = canvas.ImageFillContain

		title := widget.NewLabel("RepoRover")
		title.Alignment = fyne.TextAlignCenter
		title.TextStyle = fyne.TextStyle{Bold: true}

		subtitle := widget.NewLabel("Universal Linux package manager updater.")
		subtitle.Alignment = fyne.TextAlignCenter
		subtitle.Wrapping = fyne.TextWrapWord

		version := widget.NewLabel("Version " + appVersion)
		version.Alignment = fyne.TextAlignCenter

		body := widget.NewLabel(
			"Features:\n" +
				"• Detects supported package managers\n" +
				"• Runs updates with pkexec where needed\n" +
				"• Optional Flatpak and Snap updates\n" +
				"• Safer Arch / AUR flow when an AUR helper is present\n" +
				"• Supports paru or yay for interactive AUR updates\n" +
				"• Health-checks AUR helpers before using them\n" +
				"• Launches AUR updates in terminal for interactive prompts\n" +
				"• Dry Run and System Only modes\n" +
				"• Copyable log output for testing\n\n" +
				"Notes:\n" +
				"• Some AppImage environments may require FUSE\n" +
				"• Reboot detection depends on distro support\n" +
				"• System package tools are provided by the host distro\n" +
				"• paru or yay is only used if already installed",
		)
		body.Wrapping = fyne.TextWrapWord

		buttons := container.NewCenter(
			container.NewHBox(
				widget.NewButton("GitHub", func() {
					openLink("https://github.com/RossContino1/RepoRover")
				}),
				widget.NewButton("Website", func() {
					openLink("https://bytesbreadbbq.com/")
				}),
				widget.NewButton("Close", func() {
					aboutWin.Close()
				}),
			),
		)

		content := container.NewBorder(
			nil,
			buttons,
			nil,
			nil,
			container.NewPadded(
				container.NewVBox(
					container.NewCenter(icon),
					title,
					subtitle,
					version,
					widget.NewSeparator(),
					body,
				),
			),
		)

		aboutWin.SetContent(content)
		aboutWin.Show()
	}

	showHelpWindow := func() {
		helpWin := a.NewWindow("RepoRover Help")
		helpWin.Resize(fyne.NewSize(780, 620))
		helpWin.SetIcon(iconRes)

		helpText := widget.NewLabel(
			`RepoRover Help

Version ` + appVersion + `

Overview
RepoRover checks for supported Linux update systems and can run updates from a GUI using pkexec for privileged commands.

Basic Operation
1. Click "Check Systems" to detect your distro and available package managers.
2. Review the detected OS and primary package manager.
3. Choose your options:
   - Include Flatpak
   - Include Snap
   - Dry Run
   - System Only
4. Click "Run Updates".

Options
Include Flatpak
Runs flatpak update if Flatpak is installed.

Include Snap
Runs snap refresh if Snap is installed.

Dry Run
Shows what would be executed without running the commands.

System Only
Runs only the primary system package manager and skips Flatpak/Snap.

Arch / AUR Support
On Arch-based systems, RepoRover:
- runs pacman -Syu for official repositories with pkexec
- checks for AUR-only updates with paru -Qua or yay -Qua if an AUR helper is installed
- health-checks AUR helpers and will fall back if one is broken
- lets you choose whether to skip the AUR step or open it in a terminal
- does not try to force AUR helper prompts through the GUI, because paru/yay may need an interactive terminal for sudo and package prompts

Progress and Logs
The status bar shows the current task.
The progress bar shows overall session progress.
The output area shows detailed command activity.
Use "Copy Log" if you want to save test results.

Reboot Checks
DNF:
Uses "dnf needs-restarting -r"
APT:
Checks /var/run/reboot-required
Other package managers may report "No reboot check available".

AppImage / FUSE Note
If you later distribute RepoRover as an AppImage, some systems may require FUSE support for the AppImage to run.

Privileges
RepoRover uses pkexec for system package manager commands that need elevation.
That means the desktop environment may prompt for an administrator password.

AUR Terminal Note
When AUR updates are detected, RepoRover tries to launch a terminal emulator and run paru -Sua or yay -Sua there.
That allows sudo prompts, package prompts, and other AUR helper interactions to behave normally.

Testing Tips
- Use Dry Run first on a new distro
- Use System Only if you want to skip Flatpak and Snap
- Launch the installed desktop entry if you want proper dock icon behavior on GNOME
- Test detect, dry run, real run, Flatpak, Snap, AUR, and reboot status separately`,
		)
		helpText.Wrapping = fyne.TextWrapWord

		helpScroll := container.NewScroll(container.NewPadded(helpText))

		helpWin.SetContent(container.NewBorder(
			nil,
			container.NewCenter(
				container.NewHBox(
					widget.NewButton("GitHub", func() {
						openLink("https://github.com/RossContino1/RepoRover")
					}),
					widget.NewButton("Website", func() {
						openLink("https://bytesbreadbbq.com/")
					}),
					widget.NewButton("Close", func() {
						helpWin.Close()
					}),
				),
			),
			nil,
			nil,
			helpScroll,
		))

		helpWin.Show()
	}

	checkSystemsAction := func() {
		if !startOperation("check systems") {
			return
		}
		setStatus("Checking systems...", 0)

		go func() {
			defer endOperation()

			info := getOSInfo()
			managers := detectManagers()
			pm := detectPrimaryManager(managers)
			if pm == "" {
				pm = "none detected"
			}

			setHeader(
				"Detected OS: "+info.PrettyName,
				"Primary package manager: "+pm,
				"Reboot status: not checked yet",
			)

			appendOutput("")
			appendOutput("Checking available update systems...")
			appendOutput("OS: " + info.PrettyName)
			appendOutput("pkexec: " + commandStatus("pkexec"))
			appendOutput("flatpak: " + commandStatus("flatpak"))
			appendOutput("snap: " + commandStatus("snap"))
			appendOutput("paru: " + commandStatus("paru"))
			appendOutput("yay: " + commandStatus("yay"))

			if hasCommand("pacman") {
				aurStatus := detectWorkingAURHelper(appendOutput)
				if aurStatus.Working {
					appendOutput("AUR helper: working (" + aurStatus.Name + ")")
					if term, ok := detectTerminalCommand(); ok {
						appendOutput("terminal emulator: detected (" + term.Name + ")")
					} else {
						appendOutput("terminal emulator: not detected")
					}
				} else {
					appendOutput("AUR helper: no working helper detected")
				}
			}

			if len(managers) == 0 {
				appendOutput("No supported update systems detected.")
				setStatus("No supported update systems detected", 1)
				return
			}

			for _, m := range managers {
				kind := "secondary"
				if m.IsSystem {
					kind = "system"
				}
				appendOutput(fmt.Sprintf("✔ %s (%s)", m.Name, kind))
			}

			if hasCommand("pacman") {
				aurStatus := detectWorkingAURHelper(nil)
				if aurStatus.Working {
					appendOutput("✔ working AUR helper detected: " + aurStatus.Name + " (AUR updates can be reviewed before running)")
				} else {
					appendOutput("➜ no working AUR helper detected (AUR updates will be skipped)")
				}
			}

			setStatus("System check complete", 1)
		}()
	}

	runUpdatesAction := func() {
		if !startOperation("run updates") {
			return
		}
		setStatus("Preparing update session...", 0)

		go func() {
			defer endOperation()

			info := getOSInfo()
			managers := detectManagers()
			pm := detectPrimaryManager(managers)
			if pm == "" {
				pm = "none detected"
			}

			setHeader(
				"Detected OS: "+info.PrettyName,
				"Primary package manager: "+pm,
				"Reboot status: running...",
			)

			appendOutput("")
			appendOutput("Starting update session...")

			if len(managers) == 0 {
				appendOutput("No supported update systems detected.")
				setHeader(
					"Detected OS: "+info.PrettyName,
					"Primary package manager: "+pm,
					"Reboot status: not available",
				)
				setStatus("No supported update systems detected", 1)
				return
			}

			selected := buildSelectedManagers(
				managers,
				includeFlatpak.Checked,
				includeSnap.Checked,
				systemOnly.Checked,
				appendOutput,
			)

			totalSteps := len(selected)
			if !dryRun.Checked {
				totalSteps++
			}

			if totalSteps == 0 {
				appendOutput("Nothing selected to run.")
				setStatus("Nothing selected to run", 1)
				return
			}

			done := 0
			for _, m := range selected {
				setStatus("Running "+m.Name+"...", float64(done)/float64(totalSteps))

				ok := runManager(a, w, m, dryRun.Checked, appendOutput)
				if !ok {
					appendOutput("Update session stopped because " + m.Name + " failed.")
					setStatus(m.Name+" failed", 1)
					return
				}

				done++
				setStatus(m.Name+" complete", float64(done)/float64(totalSteps))
			}

			if dryRun.Checked {
				setHeader(
					"Detected OS: "+info.PrettyName,
					"Primary package manager: "+pm,
					"Reboot status: skipped during dry run",
				)
				appendOutput("Dry run complete.")
				setStatus("Dry run complete", 1)
				return
			}

			setStatus("Checking reboot status...", float64(done)/float64(totalSteps))
			rebootMessage := checkRebootStatus()
			done++

			setHeader(
				"Detected OS: "+info.PrettyName,
				"Primary package manager: "+pm,
				"Reboot status: "+rebootMessage,
			)

			appendOutput("All selected update commands finished.")
			appendOutput("Update session complete.")
			appendOutput("Reboot status: " + rebootMessage)
			setStatus("Update session complete", float64(done)/float64(totalSteps))
		}()
	}

	clearOutputAction := func() {
		if busy.Load() {
			appendOutput("Cannot clear output while another operation is running.")
			return
		}
		clearOutput()
		setStatus("Output cleared", 0)
	}

	copyLogAction := func() {
		copyLog()
	}

	checkButton = widget.NewButton("Check Systems", checkSystemsAction)
	runButton = widget.NewButton("Run Updates", runUpdatesAction)
	clearButton := widget.NewButton("Clear Output", clearOutputAction)
	copyButton := widget.NewButton("Copy Log", copyLogAction)

	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Check Systems", checkSystemsAction),
		fyne.NewMenuItem("Run Updates", runUpdatesAction),
		fyne.NewMenuItem("Copy Log", copyLogAction),
		fyne.NewMenuItem("Clear Output", clearOutputAction),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			if busy.Load() {
				appendOutput("Cannot quit while another operation is running.")
				return
			}
			w.Close()
		}),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Program Help", showHelpWindow),
		fyne.NewMenuItem("About RepoRover", showAbout),
	)

	w.SetMainMenu(fyne.NewMainMenu(fileMenu, helpMenu))

	titleIcon := canvas.NewImageFromResource(iconRes)
	titleIcon.SetMinSize(fyne.NewSize(28, 28))
	titleIcon.FillMode = canvas.ImageFillContain

	titleText := widget.NewLabel("  RepoRover " + appVersion)
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	title := container.NewCenter(
		container.NewHBox(
			titleIcon,
			titleText,
		),
	)

	leftPanel := container.NewVBox(
		widget.NewLabel("Options"),
		includeFlatpak,
		includeSnap,
		dryRun,
		systemOnly,
	)

	topPanel := container.NewVBox(
		title,
		widget.NewSeparator(),
		osLabel,
		pmLabel,
		rebootLabel,
		widget.NewSeparator(),
		container.NewCenter(
			container.NewHBox(
				checkButton,
				runButton,
				clearButton,
				copyButton,
			),
		),
	)

	firstRunNote := widget.NewLabel("Note: First run may take several minutes while RepoRover gathers updates.")
	firstRunNote.Wrapping = fyne.TextWrapWord

	bottomPanel := container.NewVBox(
		statusLabel,
		progressBar,
		firstRunNote,
	)

	content := container.NewBorder(
		topPanel,
		bottomPanel,
		leftPanel,
		nil,
		outputScroll,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func commandStatus(name string) string {
	if hasCommand(name) {
		return "detected"
	}
	return "not detected"
}

func buildSelectedManagers(managers []ManagerInfo, includeFlatpak, includeSnap, systemOnly bool, log func(string)) []ManagerInfo {
	var selected []ManagerInfo

	for _, m := range managers {
		if m.IsSystem {
			selected = append(selected, m)
		}
	}

	if systemOnly {
		log("➜ secondary package managers skipped by System Only")
		return selected
	}

	for _, m := range managers {
		if m.IsSystem {
			continue
		}
		if m.Name == "flatpak" && !includeFlatpak {
			log("➜ flatpak skipped")
			continue
		}
		if m.Name == "snap" && !includeSnap {
			log("➜ snap skipped")
			continue
		}
		selected = append(selected, m)
	}

	return selected
}

func getOSInfo() OSInfo {
	info := OSInfo{
		ID:         "unknown",
		PrettyName: "Unknown Linux",
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return info
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := strings.Trim(parts[1], `"`)
		values[key] = value
	}

	if v, ok := values["ID"]; ok && v != "" {
		info.ID = v
	}
	if v, ok := values["PRETTY_NAME"]; ok && v != "" {
		info.PrettyName = v
	}

	return info
}

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

type AURHelperStatus struct {
	Name     string
	Tested   bool
	Working  bool
	Details  string
	ExitCode int
}

func detectAURHelper() string {
	return detectWorkingAURHelper(nil).Name
}

func detectWorkingAURHelper(log func(string)) AURHelperStatus {
	candidates := []string{"paru", "yay"}
	var firstFailure AURHelperStatus

	for _, helper := range candidates {
		status := probeAURHelper(helper)

		if log != nil {
			switch {
			case !status.Tested:
				log("AUR helper check: " + helper + " not installed")
			case status.Working:
				log("AUR helper check: " + helper + " OK")
			default:
				log("AUR helper check: " + helper + " failed: " + status.Details)
			}
		}

		if status.Working {
			return status
		}
		if status.Tested && firstFailure.Name == "" {
			firstFailure = status
		}
	}

	return firstFailure
}

func probeAURHelper(helper string) AURHelperStatus {
	status := AURHelperStatus{Name: helper}

	if !hasCommand(helper) {
		status.Details = "not installed"
		return status
	}

	status.Tested = true

	out, err := runCommand(helper, "--version")
	trimmed := strings.TrimSpace(out)

	if err == nil {
		status.Working = true
		if trimmed == "" {
			status.Details = "version check succeeded"
		} else {
			status.Details = "version check succeeded: " + oneLine(trimmed)
		}
		return status
	}

	status.Details = explainCommandError(err)
	if trimmed != "" {
		status.Details += " | " + oneLine(trimmed)
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		status.ExitCode = exitErr.ExitCode()
	}
	return status
}

func oneLine(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " | ")
	return s
}

func detectManagers() []ManagerInfo {
	all := []ManagerInfo{
		{Name: "dnf", IsSystem: true},
		{Name: "apt", IsSystem: true},
		{Name: "zypper", IsSystem: true},
		{Name: "pacman", IsSystem: true},
		{Name: "flatpak", IsSystem: false},
		{Name: "snap", IsSystem: false},
	}

	var found []ManagerInfo
	for _, m := range all {
		if hasCommand(m.Name) {
			found = append(found, m)
		}
	}
	return found
}

func detectPrimaryManager(managers []ManagerInfo) string {
	for _, m := range managers {
		if m.IsSystem {
			return m.Name
		}
	}
	return ""
}

func runManager(a fyne.App, w fyne.Window, m ManagerInfo, dryRun bool, log func(string)) bool {
	log("")
	log("Running " + m.Name + "...")

	switch m.Name {
	case "dnf":
		return runSteps(m.Name, dryRun, log,
			[]string{"pkexec", "dnf", "upgrade", "--refresh", "-y"},
		)
	case "apt":
		return runSteps(m.Name, dryRun, log,
			[]string{"pkexec", "apt", "update"},
			[]string{"pkexec", "apt", "upgrade", "-y"},
		)
	case "zypper":
		return runSteps(m.Name, dryRun, log,
			[]string{"pkexec", "zypper", "refresh"},
			[]string{"pkexec", "zypper", "update", "-y"},
		)
	case "pacman":
		return runPacmanFlow(a, w, dryRun, log)
	case "flatpak":
		return runSteps(m.Name, dryRun, log,
			[]string{"flatpak", "update", "-y"},
		)
	case "snap":
		return runSteps(m.Name, dryRun, log,
			[]string{"pkexec", "snap", "refresh"},
		)
	default:
		log("✖ unsupported manager: " + m.Name)
		return false
	}
}

func runSteps(managerName string, dryRun bool, log func(string), steps ...[]string) bool {
	for _, step := range steps {
		log("> " + strings.Join(step, " "))

		if dryRun {
			log("  [dry-run] not executed")
			continue
		}

		out, err := runCommand(step[0], step[1:]...)
		logCommandOutput(out, log, "  command completed with no terminal output")

		if err != nil {
			log("✖ " + managerName + " failed: " + explainCommandError(err))
			return false
		}

		log("  step completed successfully")
	}

	log("✔ " + managerName + " complete")
	return true
}

func runPacmanFlow(a fyne.App, w fyne.Window, dryRun bool, log func(string)) bool {
	log("")
	log("Pacman official repository update phase...")
	pacmanStep := []string{"pkexec", "pacman", "-Syu", "--noconfirm", "--color=never"}
	log("> " + strings.Join(pacmanStep, " "))

	if dryRun {
		log("  [dry-run] official repo update would run with --noconfirm")
	} else {
		out, err := runCommand(pacmanStep[0], pacmanStep[1:]...)
		logCommandOutput(out, log, "  pacman completed with no terminal output")
		if err != nil {
			log("✖ pacman failed: " + explainCommandError(err))
			return false
		}
		log("  pacman official repo step completed successfully")
	}

	aurStatus := detectWorkingAURHelper(log)
	aurHelper := aurStatus.Name
	if aurHelper == "" || !aurStatus.Working {
		log("➜ no working AUR helper detected; skipping AUR updates")
		log("✔ pacman complete")
		return true
	}

	aurUpdates, aurErr := getAURUpdates(aurHelper)
	if aurErr != nil {
		log("➜ selected AUR helper failed update query: " + aurHelper + " -Qua: " + aurErr.Error())

		fallback := alternateAURHelper(aurHelper)
		if fallback != "" {
			fallbackStatus := probeAURHelper(fallback)
			if fallbackStatus.Working {
				log("➜ falling back to alternate AUR helper: " + fallback)
				aurHelper = fallback
				aurUpdates, aurErr = getAURUpdates(aurHelper)
			}
		}

		if aurErr != nil {
			log("➜ RepoRover will skip AUR for this session")
			log("✔ pacman complete")
			return true
		}
	}

	if len(aurUpdates) == 0 {
		log("➜ no AUR updates found with " + aurHelper)
		log("✔ pacman complete")
		return true
	}

	log("")
	log("===== AUR UPDATES DETECTED =====")
	log("AUR helper: " + aurHelper)
	log(fmt.Sprintf("AUR package count: %d", len(aurUpdates)))
	for _, pkg := range aurUpdates {
		log("  • " + pkg)
	}
	log("================================")

	if dryRun {
		log("  [dry-run] AUR decision dialog would appear here")
		log("  [dry-run] if selected, RepoRover would try to open a terminal for: " + aurHelper + " -Sua")
		log("✔ pacman complete")
		return true
	}

	decision, ok := promptAURDecision(a, w, aurHelper, aurUpdates)
	if !ok {
		log("✖ AUR choice dialog could not be displayed")
		return false
	}
	if !decision.Proceed {
		log("Update session cancelled at AUR decision dialog")
		return false
	}

	switch decision.Action {
	case AURActionSkip:
		log("➜ user chose to skip the AUR step")
	case AURActionOpenTerminal:
		log("➜ user chose to open AUR updates in terminal")
		usedHelper, err := launchAURInTerminalWithFallback(aurHelper, log)
		if err != nil {
			log("✖ could not launch terminal for AUR updates: " + err.Error())
			return false
		}
		log("  terminal launched for: " + usedHelper + " -Sua")
		log("  complete the AUR prompts in the terminal window")
		log("  when finished, review the result and close that terminal window")
	default:
		log("✖ unknown AUR action")
		return false
	}

	log("✔ pacman complete")
	return true
}

func getAURUpdates(aurHelper string) ([]string, error) {
	out, err := runCommand(aurHelper, "-Qua")
	trimmed := strings.TrimSpace(out)

	if err != nil && trimmed == "" {
		return nil, err
	}

	if trimmed == "" {
		return []string{}, nil
	}

	var updates []string
	for _, line := range strings.Split(trimmed, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		updates = append(updates, line)
	}

	return updates, nil
}

func promptAURDecision(a fyne.App, w fyne.Window, aurHelper string, aurUpdates []string) (AURDecision, bool) {
	result := make(chan AURDecision, 1)
	shown := make(chan bool, 1)

	fyne.Do(func() {
		updatesBox := widget.NewTextGrid()
		updatesBox.SetText(strings.Join(aurUpdates, "\n"))
		updatesScroll := container.NewScroll(updatesBox)
		updatesScroll.SetMinSize(fyne.NewSize(680, 220))

		radio := widget.NewRadioGroup([]string{
			"Skip AUR step for now",
			"Open terminal and run " + aurHelper + " -Sua",
		}, nil)
		radio.SetSelected(radio.Options[1])

		warning := widget.NewLabel(
			"AUR updates were detected. RepoRover can open a terminal window and run " + aurHelper + " -Sua there so sudo and package prompts work normally. When the update finishes, review the output and close the terminal window.",
		)
		warning.Wrapping = fyne.TextWrapWord

		label := widget.NewLabel(
			fmt.Sprintf("AUR packages reported by %s -Qua (%d):", aurHelper, len(aurUpdates)),
		)

		content := container.NewBorder(
			container.NewVBox(
				warning,
				widget.NewSeparator(),
				label,
			),
			nil,
			nil,
			nil,
			container.NewVBox(
				updatesScroll,
				widget.NewSeparator(),
				radio,
			),
		)

		dlg := dialog.NewCustomConfirm(
			fmt.Sprintf("AUR Updates Detected (%d)", len(aurUpdates)),
			"OK",
			"Cancel",
			content,
			func(ok bool) {
				if !ok {
					result <- AURDecision{Proceed: false}
					return
				}

				decision := AURDecision{Proceed: true, Action: AURActionSkip}
				switch radio.Selected {
				case radio.Options[1]:
					decision.Action = AURActionOpenTerminal
				default:
					decision.Action = AURActionSkip
				}
				result <- decision
			}, w)
		dlg.Resize(fyne.NewSize(760, 560))
		dlg.Show()
		shown <- true
	})

	select {
	case <-shown:
		decision := <-result
		return decision, true
	case <-time.After(2 * time.Second):
		return AURDecision{}, false
	}
}

func detectTerminalCommand() (TerminalLaunch, bool) {
	if termEnv := strings.TrimSpace(os.Getenv("TERMINAL")); termEnv != "" {
		if _, err := exec.LookPath(termEnv); err == nil {
			return TerminalLaunch{Name: termEnv}, true
		}
	}

	candidates := []string{
		"x-terminal-emulator",
		"kgx",
		"gnome-terminal",
		"konsole",
		"xfce4-terminal",
		"mate-terminal",
		"tilix",
		"lxterminal",
		"alacritty",
		"kitty",
		"wezterm",
		"xterm",
	}

	for _, name := range candidates {
		if _, err := exec.LookPath(name); err == nil {
			return TerminalLaunch{Name: name}, true
		}
	}

	return TerminalLaunch{}, false
}

func buildTerminalArgs(termName, shellCommand string) []string {
	switch termName {
	case "x-terminal-emulator":
		return []string{"-e", "sh", "-lc", shellCommand}
	case "kgx":
		return []string{"--", "sh", "-lc", shellCommand}
	case "gnome-terminal":
		return []string{"--", "sh", "-lc", shellCommand}
	case "konsole":
		return []string{"-e", "sh", "-lc", shellCommand}
	case "xfce4-terminal":
		return []string{"--command=sh -lc '" + escapeSingleQuotes(shellCommand) + "'"}
	case "mate-terminal":
		return []string{"--", "sh", "-lc", shellCommand}
	case "tilix":
		return []string{"-e", "sh", "-lc", shellCommand}
	case "lxterminal":
		return []string{"-e", "sh -lc '" + escapeSingleQuotes(shellCommand) + "'"}
	case "alacritty":
		return []string{"-e", "sh", "-lc", shellCommand}
	case "kitty":
		return []string{"sh", "-lc", shellCommand}
	case "wezterm":
		return []string{"start", "--always-new-process", "sh", "-lc", shellCommand}
	case "xterm":
		return []string{"-e", "sh", "-lc", shellCommand}
	default:
		return []string{"-e", "sh", "-lc", shellCommand}
	}
}

func launchAURInTerminal(aurHelper string) error {
	_, err := launchAURInTerminalWithFallback(aurHelper, nil)
	return err
}

func launchAURInTerminalWithFallback(aurHelper string, log func(string)) (string, error) {
	attemptOrder := []string{aurHelper}
	if fallback := alternateAURHelper(aurHelper); fallback != "" {
		attemptOrder = append(attemptOrder, fallback)
	}

	var launchErrs []string

	for i, helper := range attemptOrder {
		status := probeAURHelper(helper)
		if !status.Working {
			msg := helper + " failed health check"
			if status.Details != "" {
				msg += ": " + status.Details
			}
			launchErrs = append(launchErrs, msg)
			if log != nil {
				log("➜ " + msg)
			}
			continue
		}

		if err := launchSpecificAURInTerminal(helper); err != nil {
			launchErrs = append(launchErrs, helper+": "+err.Error())
			if log != nil {
				log("➜ terminal launch failed for " + helper + ": " + err.Error())
			}
			continue
		}

		if i > 0 && log != nil {
			log("➜ falling back to alternate AUR helper: " + helper)
		}
		return helper, nil
	}

	if len(launchErrs) == 0 {
		return "", fmt.Errorf("no usable AUR helper found")
	}
	return "", fmt.Errorf(strings.Join(launchErrs, " | "))
}

func launchSpecificAURInTerminal(aurHelper string) error {
	term, ok := detectTerminalCommand()
	if !ok {
		return fmt.Errorf("no supported terminal emulator was detected")
	}

	shellCommand := aurHelper + ` -Sua; status=$?; echo; if [ $status -eq 0 ]; then echo "AUR update finished successfully."; else echo "` + aurHelper + ` exited with status $status."; fi; sleep 2; exit $status`
	args := buildTerminalArgs(term.Name, shellCommand)

	cmd := exec.Command(term.Name, args...)
	return cmd.Start()
}

func alternateAURHelper(current string) string {
	switch current {
	case "paru":
		if hasCommand("yay") {
			return "yay"
		}
	case "yay":
		if hasCommand("paru") {
			return "paru"
		}
	}
	return ""
}

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, `'`, `'\''`)
}

func logCommandOutput(out string, log func(string), emptyMessage string) {
	if strings.TrimSpace(out) == "" {
		log(emptyMessage)
		return
	}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		log(line)
	}
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func explainCommandError(err error) string {
	if err == nil {
		return ""
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Sprintf("command exited with status %d", exitErr.ExitCode())
	}
	return err.Error()
}

func checkRebootStatus() string {
	if hasCommand("dnf") {
		err := exec.Command("pkexec", "dnf", "needs-restarting", "-r").Run()
		if err == nil {
			return "No reboot required"
		}
		if _, ok := err.(*exec.ExitError); ok {
			return "Reboot required"
		}
		return "Reboot check error"
	}

	if hasCommand("apt") {
		_, err := os.Stat("/var/run/reboot-required")
		if err == nil {
			return "Reboot required"
		}
		if os.IsNotExist(err) {
			return "No reboot required"
		}
		return "Reboot check error"
	}

	return "Reboot check not supported on this distro"
}
