package updater

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Stezok/excel-repot-service/internal/service"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Logger interface {
	Print(...interface{})
}

type UpdaterConfig struct {
	SeleniumPath string
	Mode         string
	Port         int

	ReviewPath   string
	DownloadPath string
}

type UpdateJob struct {
	HuaweiLogin    string
	HuaweiPassword string
	ProjectID      string
}

type Updater struct {
	UpdaterConfig
	JobChannel        chan UpdateJob
	Logger            Logger
	ReportService     service.ReportService
	UpdateTimeService service.UpdateTimeService
}

func (u *Updater) isLoginPage(wd selenium.WebDriver) (bool, error) {
	elems, err := wd.FindElements(selenium.ByCSSSelector, "input#uid")
	if err != nil {
		return false, err
	}

	return len(elems) != 0, nil
}

func (u *Updater) login(wd selenium.WebDriver, job UpdateJob) error {
	err := wd.Get("https://uniportal.huawei.com/uniportal/")
	if err != nil {
		return err
	}

	elem, err := wd.FindElement(selenium.ByCSSSelector, "input#uid")
	if err != nil {
		return err
	}

	login := job.HuaweiLogin
	err = elem.SendKeys(login)
	if err != nil {
		return err
	}

	elem, err = wd.FindElement(selenium.ByCSSSelector, "input#password")
	if err != nil {
		return err
	}
	pass := job.HuaweiPassword
	err = elem.SendKeys(pass)
	if err != nil {
		return err
	}

	elem, err = wd.FindElement(selenium.ByCSSSelector, "input.login_submit_pwd_v2")
	if err != nil {
		return err
	}
	err = elem.Click()
	if err != nil {
		return err
	}

	time.Sleep(3)
	return nil
}

func (u *Updater) getTaskShowPage(wd selenium.WebDriver, job UpdateJob) error {
	url := "https://isdp-ru.huawei.com/sdcloud/apt/aui/index.html#/itemToDoIsdp?currentStatus=DB&projectNumber=%s"
	url = fmt.Sprintf(url, job.ProjectID)

	err := wd.Get(url)
	if err != nil {
		return err
	}
	if isLogin, err := u.isLoginPage(wd); err != nil {
		return err
	} else if isLogin {
		err = u.login(wd, job)
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
		err = wd.Get(url)
		if err != nil {
			return err
		}
	}

	time.Sleep(5 * time.Second)
	return nil
}

func (u *Updater) clickExport(wd selenium.WebDriver) error {
	wd.SwitchFrame("microServicesIframeItemTodo")
	defer wd.SwitchFrame(nil)

	// _, err := wd.ExecuteScript(`document.querySelector("div.button-export button.hae-btn").click()`, nil)
	// return err
	elem, err := wd.FindElement(selenium.ByCSSSelector, "div.button-export button.hae-btn")
	if err != nil {
		return err
	}
	err = elem.Click()
	if err != nil {
		return err
	}

	time.Sleep(15 * time.Second)
	return nil
}

func (u *Updater) downloadReview(wd selenium.WebDriver) error {
	err := wd.Get("https://app-ru.huawei.com/sdcp/apt/#!jalor/async/task/listExport.html")
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	elem, err := wd.FindElement(selenium.ByCSSSelector, `tr.grid-row[_row="0"] td[_col="2"] a`)
	if err != nil {
		return err
	}

	err = elem.Click()
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	return nil
}

func (u *Updater) UpdateReview(job UpdateJob) error {
	defer recover()
	ops := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(u.SeleniumPath, u.Port, ops...)
	if err != nil {
		return err
	}
	defer service.Stop()

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	args := []string{
		"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7",
		"--no-sandbox",
	}

	if u.Mode == "headless" {
		args = append(args, "--headless")
	}

	prefs := make(map[string]interface{})
	prefs["download.default_directory"] = u.DownloadPath

	chromeCaps := chrome.Capabilities{
		Prefs: prefs,
		Args:  args,
	}
	caps.AddChrome(chromeCaps)

	addr := fmt.Sprintf("http://localhost:%d/wd/hub", u.Port)
	wd, err := selenium.NewRemote(caps, addr)
	if err != nil {
		return err
	}
	defer wd.Quit()

	err = u.login(wd, job)
	if err != nil {
		return err
	}

	err = u.getTaskShowPage(wd, job)
	if err != nil {
		return err
	}

	err = u.clickExport(wd)
	if err != nil {
		return err
	}

	err = u.downloadReview(wd)
	if err != nil {
		return err
	}

	paths, err := filepath.Glob(u.DownloadPath + "*.xlsx")
	if err != nil {
		return err
	}

	newFile, err := os.Open(paths[0])
	if err != nil {
		return err
	}
	defer newFile.Close()

	oldFile, err := os.Create(u.ReviewPath)
	if err != nil {
		return err
	}
	defer oldFile.Close()

	_, err = io.Copy(oldFile, newFile)
	if err != nil {
		return err
	}

	newFile.Close()
	oldFile.Close()

	for _, path := range paths {
		os.Remove(path)
	}

	_, err = u.ReportService.UpdateReports(job.ProjectID)
	if err != nil {
		return err
	}

	err = u.UpdateTimeService.SetLastUpdateTime(job.ProjectID, time.Now().Unix())
	return err
}

func (u *Updater) handle() {
	for {
		job := <-u.JobChannel
		u.UpdateReview(job)
	}
}

func (u *Updater) Run() {
	u.handle()
}

func (u *Updater) PushJobEvery(job UpdateJob, dur time.Duration) {
	go func() {
		for {
			err := u.UpdateReview(job)
			if err != nil {
				u.Logger.Print(err)
			}
			time.Sleep(dur)
		}
	}()
}

// func (u *Updater) UpdateEvery(dur time.Duration) {
// 	go func() {
// 		for {
// 			err := u.UpdateReview()
// 			if err != nil {
// 				u.Logger.Print(err)
// 			}
// 			time.Sleep(dur)
// 		}
// 	}()
// }

func NewUpdater(conf UpdaterConfig, reportService service.ReportService, updateTimeService service.UpdateTimeService) *Updater {
	return &Updater{
		UpdaterConfig:     conf,
		JobChannel:        make(chan UpdateJob, 1),
		ReportService:     reportService,
		UpdateTimeService: updateTimeService,
		Logger:            log.Default(),
	}
}
