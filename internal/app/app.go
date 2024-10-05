package app

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/brys0/goTAB/internal/fio"
	"github.com/brys0/goTAB/internal/worker"

	"github.com/brys0/goTAB/internal/api"
	"github.com/brys0/goTAB/internal/cpu"
	"github.com/brys0/goTAB/internal/gpu"
	"github.com/brys0/goTAB/internal/hardware"
	"github.com/brys0/goTAB/internal/memory"
	tabOS "github.com/brys0/goTAB/internal/os"
	"github.com/brys0/goTAB/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type App struct {
	SelectedGraphics []int          `json:"-"`
	CPUTest          bool           `json:"-"`
	Errors           []error        `json:"-"`
	API              *api.APIClient `json:"-"`
	Platform         *string        `json:"-"`
	Theme            *huh.Theme
	Token            *string                `json:"token"`
	HwInfo           *hardware.HardwareInfo `json:"hwinfo"`
	Tests            []api.Test             `json:"-"`
	TotalTests       int
}

func CreateNewTAB(testCPU bool, server string) *App {
	uri, err := url.Parse(server)

	if err != nil {
		log.Fatalf("The server url provided is invalid", "err", err)
	}

	return &App{
		SelectedGraphics: []int{},
		CPUTest:          testCPU,
		Errors:           []error{},
		API:              api.CreateNewAPI(uri),
		Platform:         nil,

		Theme: CreateJellyfinStyle(),
		// json data
		Token:      nil,
		HwInfo:     nil,
		Tests:      []api.Test{},
		TotalTests: 0,
	}
}

func CreateJellyfinStyle() *huh.Theme {
	base := huh.ThemeCharm()

	base.Blurred.BlurredButton.Background(lipgloss.Color("#303030"))
	// Jellyfin styled button
	base.Focused.Title = base.Focused.Title.Foreground(lipgloss.Color("#fff"))

	base.Focused.FocusedButton = base.Focused.FocusedButton.
		Foreground(lipgloss.Color("rgba(255, 255, 255, 0.8)")).
		Background(lipgloss.Color("#0cb0e8"))

	return base
}

func (app *App) FetchPlatformInfo() []api.Platform {
	platforms, err := app.API.FetchSupportedPlatforms()

	if err != nil {
		log.Fatal("An error occured when fetching platform info from the server", "err", err)
	}

	return platforms
}

func (app *App) GetHardware() {
	memoryInfo, err := memory.GetMemoryInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system memoryInfo info", err)
	}

	cpuInfo, err := cpu.GetCPUInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system CPU info", err)
	}

	gpuInfo, err := gpu.GetGPUInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system GPU info", err)
	}

	osInfo, err := tabOS.GetOSInfo()

	if err != nil {
		log.Fatal("Error occurred trying to get system info", err)
	}

	app.HwInfo = &hardware.HardwareInfo{
		FFmpeg: nil,
		OS:     osInfo,
		CPU:    cpuInfo,
		Memory: memoryInfo,
		GPU:    gpuInfo,
	}
}

func (app *App) ValidatePlatform() {
	if app.HwInfo == nil {
		log.Fatal("must run GetHardware() before ValidatePlatform()")
	}

	log.Info("Fetching valid platforms...")
	platforms := app.FetchPlatformInfo()

	var platform_id *string

	for i := range platforms {
		platform := platforms[i]

		if strings.Contains(strings.ToLower(app.HwInfo.OS.Name), strings.ToLower(platform.Type)) {
			platform_id = &platform.Id
		}
	}

	if platform_id == nil {
		log.Fatal("your os doesn't appear to match any on the server.", "os", app.HwInfo.OS.Name)
	}

	app.Platform = platform_id
	log.Info("Validated platform", "platform", *platform_id)
}

func (app *App) GetPlatformTests() {
	if app.Platform == nil {
		log.Fatal("must run ValidatePlatform() before GetPlatformTests()")
	}

}

func (app *App) FetchPlatformTests() {
	testRoot, err := app.API.FetchPlatformTests(*app.Platform)

	if err != nil {
		log.Fatal("An error occured when fetching platform info from the server", "err", err)
	}

	app.Token = &testRoot.Token
	app.Tests = testRoot.Tests
	app.HwInfo.FFmpeg = &testRoot.FFmpeg

	for testIndex := range app.Tests {
		test := app.Tests[testIndex]

		app.TotalTests += len(test.Data)
	}
}

func (app *App) PromptGPU() error {
	if app.HwInfo == nil {
		log.Fatal("No hardware info")
	}

	gpuSelectOptions := make([]huh.Option[int], len(app.HwInfo.GPU))

	for id := range app.HwInfo.GPU {
		card := app.HwInfo.GPU[id]
		gpuSelectOptions[id] = huh.Option[int]{Key: card.Description, Value: id}
	}

	gpuMultiSelect := huh.NewMultiSelect[int]().
		Title("Select a GPU").
		Options(gpuSelectOptions...).
		Description("Press [Space] to select and [Enter] to submit").
		Value(&app.SelectedGraphics)

	gpuSelectForm := huh.NewForm(huh.NewGroup(gpuMultiSelect))

	err := gpuSelectForm.Run()

	if err != nil {
		log.Fatal("Error occurred trying to select a GPU", err)
	}

	return nil
}

func (app *App) PromptConfirmTests() {
	if len(app.Tests) == 0 {
		log.Fatal("No tests found")
	}

	confirm := false
	huh.NewConfirm().
		Title("Test Confirmation").
		Description(fmt.Sprintf("Will run %d tests", app.TotalTests)).
		Affirmative("Yes, Run tests").
		Negative("No, exit").
		Value(&confirm).WithTheme(app.Theme).Run()

	if !confirm {
		return
	}

	log.Info("Starting TUI...")

	app.StartTUIApp()
}

// TODO: Create TUI
func (app *App) StartTUIApp() {
	if app.HwInfo.FFmpeg == nil {
		log.Fatal("FFmpeg source must be defined")
	}

	log.Info("Source url", "ffmpeg", app.HwInfo.FFmpeg)
	model := tui.CreateTUIDisplay()

	p := tea.NewProgram(model)

	go app.MainTUIApp(p)

	p.Run()
}

func (app *App) MainTUIApp(p *tea.Program) {
	app.DownloadFFmpeg(p) // Download or check ffmpeg hash
	app.DownloadVideos(p) // Download the videos to transcode from the test

	// p.ReleaseTerminal()
	selectedGraphics := strings.ToLower(app.HwInfo.GPU[app.SelectedGraphics[0]].Vendor)
	completed_tests := 0

	for t := range app.Tests {
		test := app.Tests[t]
		finished_subtests := 0

		for d := range test.Data {
			data := test.Data[d]

			for a := range data.Arguments {
				arg := data.Arguments[a]

				if arg.Type == selectedGraphics {
					p.Send(tui.TaskInfo(fmt.Sprintf("%s (%s -> %s) (Subtests %d/%d)", test.Name, data.FromResolution, data.ToResolution, finished_subtests, len(test.Data))))

					worker.CreateWorkManager(p, strings.Replace(strings.ReplaceAll(arg.Arguments, "{gpu}", "0"), "-hwaccel cuda", "-hwaccel nvdec", 1), fmt.Sprintf("videos/%v.mkv", test.Name))

					finished_subtests++
					p.Send(tui.ProgressInfo{Type: "task", Progress: float64(finished_subtests) / float64(len(test.Data))})
				}
			}
		}

		completed_tests++
		p.Send(tui.ProgressInfo{Type: "overall", Progress: float64(completed_tests) / float64(len(app.Tests))})
	}

	// p.Quit()

	// log.Fatal("Coming soon.")
}

func (app *App) DownloadFFmpeg(p *tea.Program) {
	location := "ffmpeg/" + path.Base(app.HwInfo.FFmpeg.FFmpegSourceURL)

	if !app.ConfirmHashForFile(location, app.HwInfo.FFmpeg.FFmpegHashes[0].Type, app.HwInfo.FFmpeg.FFmpegHashes[0].Hash) {
		p.Send(tui.TaskInfo(fmt.Sprintf("Downloading FFmpeg (%s)", app.HwInfo.FFmpeg.FFmpegVersion)))
		api.DownloadFile(app.HwInfo.FFmpeg.FFmpegSourceURL, location, func(ratio float64) {
			p.Send(tui.ProgressInfo{Type: "task", Progress: ratio})
			p.Send(tui.ProgressInfo{Type: "overall", Progress: ratio})
		})
	}

	// Unzip files
	p.Send(tui.TaskInfo("Extracting FFmpeg..."))
	fio.Unzip(location, "ffmpeg/bin/")
	p.Send(tui.TaskInfo("Extracted FFmpeg"))
}

func (app *App) DownloadVideos(p *tea.Program) {
	total_files := len(app.Tests)
	completed_files := 0
	for id := range app.Tests {
		test := app.Tests[id]
		location := "videos/" + path.Base(test.SourceURL)

		// Don't check hash, since it doesn't exist for
		p.Send(tui.TaskInfo(fmt.Sprintf("Downloading %s", test.Name)))

		if app.CheckFileExists(location) {
			if test.SourceHashes != nil && app.ConfirmHashForFile(location, test.SourceHashes[0].Type, test.SourceHashes[0].Hash) {
				completed_files++
				continue
			}

			completed_files++
			continue
		}

		api.DownloadFile(test.SourceURL, location, func(ratio float64) {
			p.Send(tui.ProgressInfo{Type: "task", Progress: ratio})

			// Holy math
			// Calculate overall based on num of total files
			p.Send(tui.ProgressInfo{Type: "overall", Progress: float64(completed_files)/float64(total_files) + ratio/float64(total_files)})
		})

		completed_files++
	}
}

func (app *App) ConfirmHashForFile(file string, hash_type string, hashed string) bool {
	var hasher hash.Hash
	fsFile, err := os.OpenFile(file, os.O_RDONLY, 044)

	if err != nil {
		log.Error("Could not open file to check hash sum", "file", file)

		return false
	}

	switch hash_type {
	case "md5":
		hasher = md5.New()
	case "sha512":
		hasher = sha512.New()
	case "sha256":
		hasher = sha512.New()
	default:
		log.Fatal("Did not find a hashing algorithim", "hash_type", hash_type)
	}

	_, err = io.Copy(hasher, fsFile)

	if err != nil {
		log.Fatal("Could not copy file to compare hashes", "file", file)
	}

	hashed_string := hex.EncodeToString(hasher.Sum(nil))
	return hashed_string == hashed
}

func (app *App) CheckFileExists(file string) bool {
	_, err := os.OpenFile(file, os.O_RDONLY, 044)

	return err == nil
}
