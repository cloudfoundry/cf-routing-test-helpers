package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	. "github.com/onsi/gomega"
)

var DEFAULT_TIMEOUT = 30 * time.Second

func GetAppGuid(appName string) string {
	cfApp := cf.Cf("app", appName, "--guid")
	Eventually(cfApp, DEFAULT_TIMEOUT).Should(Exit(0))

	appGuid := strings.TrimSpace(string(cfApp.Out.Contents()))
	Expect(appGuid).NotTo(Equal(""))
	return appGuid
}

func SetBackend(appName string) {
	config := helpers.LoadConfig()
	if config.Backend == "diego" {
		EnableDiego(appName)
	} else if config.Backend == "dea" {
		DisableDiego(appName)
	}
}

func EnableDiego(appName string) {
	guid := GetAppGuid(appName)
	Eventually(cf.Cf("curl", "/v2/apps/"+guid, "-X", "PUT", "-d", `{"diego": true}`), DEFAULT_TIMEOUT).Should(Exit(0))
}

func DisableDiego(appName string) {
	guid := GetAppGuid(appName)
	Eventually(cf.Cf("curl", "/v2/apps/"+guid, "-X", "PUT", "-d", `{"diego": false}`), DEFAULT_TIMEOUT).Should(Exit(0))
}

func DisableDiegoAndCheckResponse(appName, expectedSubstring string) {
	guid := GetAppGuid(appName)
	Eventually(func() string {
		response := cf.Cf("curl", "/v2/apps/"+guid, "-X", "PUT", "-d", `{"diego":false}`)
		Expect(response.Wait(DEFAULT_TIMEOUT)).To(Exit(0))

		return string(response.Out.Contents())
	}, DEFAULT_TIMEOUT, "1s").Should(ContainSubstring(expectedSubstring))
}

func AppReport(appName string, timeout time.Duration) {
	Eventually(cf.Cf("app", appName, "--guid"), timeout).Should(Exit())
	Eventually(cf.Cf("logs", appName, "--recent"), timeout).Should(Exit())
}

func RestartApp(app string, timeout time.Duration) {
	Expect(cf.Cf("restart", app).Wait(timeout)).To(Exit(0))
}

func StartApp(app string, timeout time.Duration) {
	Expect(cf.Cf("start", app).Wait(timeout)).To(Exit(0))
}

func PushApp(appName, asset, buildpackName, domain string, timeout time.Duration) {
	PushAppNoStart(appName, asset, buildpackName, domain, timeout)
	SetBackend(appName)
	StartApp(appName, timeout)
}

func GenerateAppName() string {
	return generator.PrefixedRandomName("RATS-APP-")
}

func PushAppNoStart(appName, asset, buildpackName, domain string, timeout time.Duration, args ...string) {
	allArgs := []string{"push", appName,
		"-b", buildpackName,
		"--no-start",
		"-m", DEFAULT_MEMORY_LIMIT,
		"-p", asset,
		"-d", domain}
	for _, v := range args {
		allArgs = append(allArgs, v)
	}
	Expect(cf.Cf(allArgs...).Wait(timeout)).To(Exit(0))
}

func ScaleAppInstances(appName string, instances int, timeout time.Duration) {
	Expect(cf.Cf("scale", appName, "-i", strconv.Itoa(instances)).Wait(timeout)).To(Exit(0))
	Eventually(func() string {
		return string(cf.Cf("app", appName).Wait(timeout).Out.Contents())
	}, timeout*2, 2*time.Second).
		Should(ContainSubstring(fmt.Sprintf("instances: %d/%d", instances, instances)))
}

func DeleteApp(appName string, timeout time.Duration) {
	Expect(cf.Cf("delete", appName, "-f", "-r").Wait(timeout)).To(Exit(0))
}
