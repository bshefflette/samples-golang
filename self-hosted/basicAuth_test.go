package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/tebeka/selenium"
	slog "github.com/tebeka/selenium/log"
)

type TestHarness struct {
	wd           selenium.WebDriver
	service      *selenium.Service
	Capabilities selenium.Capabilities
	DefaultURL   string
	URL          *url.URL
}

// GetWebDriver returns the Selenium WebDriver
func (th *TestHarness) GetWebDriver() selenium.WebDriver {
	return th.wd
}

func NewTestHarness(ctx *godog.ScenarioContext, c selenium.Capabilities, defaultURL string) *TestHarness {
	th := &TestHarness{
		Capabilities: c,
		DefaultURL:   defaultURL,
	}

	ctx.BeforeScenario(th.BeforeScenario)
	ctx.AfterScenario(th.AfterScenario)

	return th
}

func (th *TestHarness) BeforeScenario(s *godog.Scenario) {
	var err error

	currentOS := runtime.GOOS
	chromeDriverPath := "selenium/chromedriver-90.0.4430.24-linux64"
	if currentOS == "darwin" {
		chromeDriverPath = "selenium/chromedriver-90.0.4430.24-mac64"
	}
	const (
		// These paths will be different on your system.
		seleniumPath = "selenium/selenium-server-standalone-3.141.59.jar"
		port         = 4444
	)

	sopts := []selenium.ServiceOption{
		// selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.ChromeDriver(chromeDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(nil),                    // Output debug information to STDERR.
	}

	th.service, err = selenium.NewSeleniumService(seleniumPath, port, sopts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}

	th.Capabilities.SetLogLevel(slog.Server, slog.Off)
	th.Capabilities.SetLogLevel(slog.Browser, slog.Off)
	th.Capabilities.SetLogLevel(slog.Client, slog.Off)
	th.Capabilities.SetLogLevel(slog.Driver, slog.Off)
	th.Capabilities.SetLogLevel(slog.Performance, slog.Off)
	th.Capabilities.SetLogLevel(slog.Profiler, slog.Off)
	th.wd, err = selenium.NewRemote(th.Capabilities, th.DefaultURL)
	if err != nil {
		log.Panic(err)
	}
}

func (th *TestHarness) AfterScenario(s *godog.Scenario, err error) {
	th.GetWebDriver().Quit()
	th.service.Stop()
}

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts) // godog v0.10.0 and earlier
	godog.BindCommandLineFlags("godog.", &opts)        // godog v0.11.0 (latest)

}

func TestMain(m *testing.M) {
	status := godog.TestSuite{
		Name:                "BasicAuthTest",
		ScenarioInitializer: BasicAuthInitalizer,
	}.Run()

	os.Exit(status)
}

func BasicAuthInitalizer(ctx *godog.ScenarioContext) {

	th := NewTestHarness(
		ctx,
		selenium.Capabilities{"browserName": "chrome"},
		"",
	)

	ctx.Step(`^I am an annymous user$`, th.iAmAnAnnymousUser)
	ctx.Step(`^I navigate to /([^"]*)$`, th.iNavigateTo)
	ctx.Step(`^I fill in my Password$`, th.iFillInMyPassword)
	ctx.Step(`^I fill in my username$`, th.iFillInMyUsername)
	ctx.Step(`^I should see my profile page details$`, th.iShouldSeeMyProfilePageDetails)
	ctx.Step(`^I submit the login form$`, th.iSubmitTheLoginForm)
}

func (th *TestHarness) iNavigateTo(url string) error {
	err := th.GetWebDriver().Get("http://localhost:8080/" + url)
	if err != nil {
		return err
	}

	err = th.GetWebDriver().WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		_, err := th.GetWebDriver().FindElement(selenium.ByCSSSelector, "[name=\"identifier\"]")
		if err != nil {
			return false, nil
		}
		return true, nil
	}, selenium.DefaultWaitTimeout, 1*time.Second)

	return err
}

func (th *TestHarness) iAmAnAnnymousUser() error {
	return nil
}

func (th *TestHarness) iFillInMyPassword() error {
	elem, err := th.GetWebDriver().FindElement(selenium.ByCSSSelector, "[name=\"credentials.passcode\"]")
	if err != nil {
		return err
	}

	if err := elem.Clear(); err != nil {
		return err
	}

	err = elem.SendKeys("password")

	return nil
}

func (th *TestHarness) iFillInMyUsername() error {
	elem, err := th.GetWebDriver().FindElement(selenium.ByCSSSelector, "[name=\"identifier\"]")
	if err != nil {
		return err
	}

	if err := elem.Clear(); err != nil {
		return err
	}

	err = elem.SendKeys("username")

	return nil
}

func (th *TestHarness) iShouldSeeMyProfilePageDetails() error {

	err := th.GetWebDriver().WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		_, err := th.GetWebDriver().FindElement(selenium.ByCSSSelector, ".profilePage")
		if err != nil {
			return false, nil
		}
		return true, nil
	}, selenium.DefaultWaitTimeout, 1*time.Second)
	if err != nil {
		return err
	}

	err = th.GetWebDriver().WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		_, err := th.GetWebDriver().FindElement(selenium.ByCSSSelector, "#claim-preferred_username")
		if err != nil {
			return false, nil
		}
		return true, nil
	}, selenium.DefaultWaitTimeout, 1*time.Second)
	if err != nil {
		return err
	}

	outputDiv, err := th.GetWebDriver().FindElement(selenium.ByCSSSelector, "#claim-preferred_username")
	if err != nil {
		return err
	}

	output, err := outputDiv.Text()

	if output != "username" {
		return err
	}

	return nil
}

func (th *TestHarness) iSubmitTheLoginForm() error {
	btn, err := th.GetWebDriver().FindElement(selenium.ByCSSSelector, "[data-type=\"save\"]")
	if err != nil {
		return err
	}

	if err := btn.Click(); err != nil {
		return err
	}
	return nil
}
