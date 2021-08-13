package updater

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

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

type Updater struct {
	UpdaterConfig
	Logger Logger
}

func (u *Updater) isLoginPage(wd selenium.WebDriver) (bool, error) {
	elems, err := wd.FindElements(selenium.ByCSSSelector, "input#uid")
	if err != nil {
		return false, err
	}

	return len(elems) != 0, nil
}

func (u *Updater) login(wd selenium.WebDriver) error {
	elem, err := wd.FindElement(selenium.ByCSSSelector, "input#uid")
	if err != nil {
		return err
	}

	login := os.Getenv("HUAWEI_LOGIN")
	err = elem.SendKeys(login)
	if err != nil {
		return err
	}

	elem, err = wd.FindElement(selenium.ByCSSSelector, "input#password")
	if err != nil {
		return err
	}
	pass := os.Getenv("HUAWEI_PASS")
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

	return nil
}

func (u *Updater) getTaskShowPage(wd selenium.WebDriver) error {
	url := "https://app-ru.huawei.com/sdcp/apt/apt#!epasbrgtask/index/myTaskShow.html?currentStatus=DB&indexCode=undefined"
	err := wd.Get(url)
	if err != nil {
		return err
	}
	if isLogin, err := u.isLoginPage(wd); err != nil {
		return err
	} else if isLogin {
		err = u.login(wd)
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
		err = wd.Get(url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Updater) clickExport(wd selenium.WebDriver) error {
	wd.SwitchFrame("microServicesIframeItemTodo")
	defer wd.SwitchFrame(nil)

	elem, err := wd.FindElement(selenium.ByCSSSelector, "div.button-export button.hae-btn")
	if err != nil {
		return err
	}
	return elem.Click()
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

	return elem.Click()
}

func (u *Updater) UpdateReview() error {
	ops := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(u.SeleniumPath, u.Port, ops...)
	if err != nil {
		return err
	}
	defer service.Stop()

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	args := []string{"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7"}
	if u.Mode == "headless" {
		args = append(args, []string{
			"--headless",
			"--no-sandbox",
		}...)
	}

	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: args,
	}
	caps.AddChrome(chromeCaps)

	addr := fmt.Sprintf("http://localhost:%d/wd/hub", u.Port)
	wd, err := selenium.NewRemote(caps, addr)
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	err = u.getTaskShowPage(wd)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)

	err = u.clickExport(wd)
	if err != nil {
		panic(err)
	}
	time.Sleep(20 * time.Second)

	err = u.downloadReview(wd)
	if err != nil {
		panic(err)
	}

	paths, err := filepath.Glob(u.DownloadPath + "*.xlsx")

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
	for _, path := range paths {
		os.Remove(path)
	}

	return nil
}

func (u *Updater) UpdateEvery(dur time.Duration) {
	go func() {
		for {
			err := u.UpdateReview()
			if err != nil {
				u.Logger.Print(err)
			}
			time.Sleep(dur)
		}
	}()
}

func NewUpdater(conf UpdaterConfig) *Updater {
	return &Updater{
		UpdaterConfig: conf,
	}
}
