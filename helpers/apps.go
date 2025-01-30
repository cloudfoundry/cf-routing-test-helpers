package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry/cf-test-helpers/v2/cf"
	"github.com/cloudfoundry/cf-test-helpers/v2/generator"
	. "github.com/onsi/gomega"
)

func GetAppGuid(appName string, timeout time.Duration) string {
	cfApp := cf.Cf("app", appName, "--guid")
	Eventually(cfApp, timeout).Should(Exit(0))

	appGuid := strings.TrimSpace(string(cfApp.Out.Contents()))
	Expect(appGuid).NotTo(Equal(""))
	return appGuid
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
	InstancesRunning(app, 1, timeout)
}

func InstancesRunning(appName string, instances int, timeout time.Duration) {
	Eventually(func() string {
		return string(cf.Cf("app", appName).Wait(timeout).Out.Contents())
	}, timeout*2, 2*time.Second).
		Should(MatchRegexp(fmt.Sprintf("instances:\\s+%d/%d", instances, instances)))
}

func PushApp(appName, asset, buildpackName, domain string, timeout time.Duration, memoryLimit string) {
	PushAppNoStart(appName, asset, buildpackName, domain, timeout, memoryLimit)
	StartApp(appName, timeout)
}

func GenerateAppName() string {
	return generator.PrefixedRandomName("RATS", "APP")
}

func PushAppNoStart(appName, asset, buildpackName, domain string, timeout time.Duration, memoryLimit string, args ...string) {
	flags := map[string]string{
		"-b": buildpackName,
		"-m": memoryLimit,
		"-p": asset,
		"-d": domain,
	}

	for flag, value := range flags {
		if value == "" {
			continue
		}
		args = append([]string{flag, value}, args...)
	}

	args = append([]string{"push", appName, "--no-start"}, args...)

	Expect(cf.Cf(args...).Wait(timeout)).To(Exit(0))
}

func ScaleAppInstances(appName string, instances int, timeout time.Duration) {
	Expect(cf.Cf("scale", appName, "-i", strconv.Itoa(instances)).Wait(timeout)).To(Exit(0))
	InstancesRunning(appName, instances, timeout)
}

func DeleteApp(appName string, timeout time.Duration) {
	Expect(cf.Cf("delete", appName, "-f", "-r").Wait(timeout)).To(Exit(0))
}
