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

const appVersion = "0.1.0"

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

func main() {
	iconRes := fyne.NewStaticResource("icon.png", iconBytes)

	a := app.NewWithID("com.bytesbreadbbq.sysupdate")
	a.Settings().SetTheme(theme.LightTheme())
	a.SetIcon(iconRes)

	w := a.NewWindow("SysUpdate")
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
		&widget.TextSegment{Text: "SysUpdate GUI ready.\n"},
	)
	output.Wrapping = fyne.TextWrapWord
	outputScroll := container.NewScroll(output)

	var plainLog strings.Builder
	plainLog.WriteString("SysUpdate GUI ready.\n")

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
				&widget.TextSegment{Text: "SysUpdate GUI ready.\n"},
			}
			output.Refresh()
			plainLog.Reset()
			plainLog.WriteString("SysUpdate GUI ready.\n")
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
		aboutWin := a.NewWindow("About SysUpdate")
		aboutWin.Resize(fyne.NewSize(540, 380))
		aboutWin.SetIcon(iconRes)

		icon := canvas.NewImageFromResource(iconRes)
		icon.SetMinSize(fyne.NewSize(64, 64))
		icon.FillMode = canvas.ImageFillContain

		title := widget.NewLabel("SysUpdate")
		title.Alignment = fyne.TextAlignCenter
		title.TextStyle = fyne.TextStyle{Bold: true}

		subtitle := widget.NewLabel("Linux update utility with GUI and CLI goals.")
		subtitle.Alignment = fyne.TextAlignCenter
		subtitle.Wrapping = fyne.TextWrapWord

		version := widget.NewLabel("Version " + appVersion)
		version.Alignment = fyne.TextAlignCenter

		body := widget.NewLabel(
			"Features:\n" +
				"• Detects supported package managers\n" +
				"• Runs updates with pkexec where needed\n" +
				"• Optional Flatpak and Snap updates\n" +
				"• Dry Run and System Only modes\n" +
				"• Copyable log output for testing\n\n" +
				"Notes:\n" +
				"• Some AppImage environments may require FUSE\n" +
				"• Reboot detection depends on distro support\n" +
				"• System package tools are provided by the host distro",
		)
		body.Wrapping = fyne.TextWrapWord

		buttons := container.NewCenter(
			container.NewHBox(
				widget.NewButton("GitHub", func() {
					openLink("https://github.com/RossContino1/SysUpdate")
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
		helpWin := a.NewWindow("SysUpdate Help")
		helpWin.Resize(fyne.NewSize(760, 560))
		helpWin.SetIcon(iconRes)

		helpText := widget.NewLabel(
			`SysUpdate Help

Version ` + appVersion + `

Overview
SysUpdate checks for supported Linux update systems and can run updates from a GUI using pkexec for privileged commands.

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
If you later distribute SysUpdate as an AppImage, some systems may require FUSE support for the AppImage to run.

Privileges
SysUpdate uses pkexec for system package manager commands that need elevation.
That means the desktop environment may prompt for an administrator password.

Testing Tips
- Use Dry Run first on a new distro
- Use System Only if you want to skip Flatpak and Snap
- Launch the installed desktop entry if you want proper dock icon behavior on GNOME
- Test detect, dry run, real run, Flatpak, Snap, and reboot status separately`,
		)
		helpText.Wrapping = fyne.TextWrapWord

		helpScroll := container.NewScroll(container.NewPadded(helpText))

		helpWin.SetContent(container.NewBorder(
			nil,
			container.NewCenter(
				container.NewHBox(
					widget.NewButton("GitHub", func() {
						openLink("https://github.com/RossContino1/SysUpdate")
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
				runManager(m, dryRun.Checked, appendOutput)
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
		fyne.NewMenuItem("About SysUpdate", showAbout),
	)

	w.SetMainMenu(fyne.NewMainMenu(fileMenu, helpMenu))

	titleIcon := canvas.NewImageFromResource(iconRes)
	titleIcon.SetMinSize(fyne.NewSize(28, 28))
	titleIcon.FillMode = canvas.ImageFillContain

	titleText := widget.NewLabel("  SysUpdate " + appVersion)
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

	bottomPanel := container.NewVBox(
		statusLabel,
		progressBar,
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

func runManager(m ManagerInfo, dryRun bool, log func(string)) {
	log("")
	log("Running " + m.Name + "...")

	var steps [][]string

	switch m.Name {
	case "dnf":
		steps = [][]string{
			{"pkexec", "dnf", "upgrade", "--refresh", "-y"},
		}
	case "apt":
		steps = [][]string{
			{"pkexec", "apt", "update"},
			{"pkexec", "apt", "upgrade", "-y"},
		}
	case "zypper":
		steps = [][]string{
			{"pkexec", "zypper", "refresh"},
			{"pkexec", "zypper", "update", "-y"},
		}
	case "pacman":
		steps = [][]string{
			{"pkexec", "pacman", "-Syu", "--noconfirm", "--color=never"},
		}
	case "flatpak":
		steps = [][]string{
			{"flatpak", "update", "-y"},
		}
	case "snap":
		steps = [][]string{
			{"pkexec", "snap", "refresh"},
		}
	default:
		log("✖ unsupported manager: " + m.Name)
		return
	}

	for _, step := range steps {
		log("> " + strings.Join(step, " "))

		if dryRun {
			log("  [dry-run] not executed")
			continue
		}

		out, err := runCommand(step[0], step[1:]...)

		if strings.TrimSpace(out) != "" {
			for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
				log(line)
			}
		} else {
			log("  command completed with no terminal output")
		}

		if err != nil {
			log("✖ " + m.Name + " failed: " + explainCommandError(err))
			return
		}

		log("  step completed successfully")
	}

	log("✔ " + m.Name + " complete")
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
