package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"code.cloudfoundry.org/cf-routing-test-helpers/schema"
	"github.com/cloudfoundry/cf-test-helpers/v2/cf"
)

func GetOrgQuotaDefinitionUrl(orgGuid string, timeout time.Duration) (string, error) {
	orgGuid = strings.TrimSuffix(orgGuid, "\n")
	response := cf.Cf("curl", fmt.Sprintf("/v3/organizations?guids=%s", string(orgGuid)))
	Expect(response.Wait(timeout)).To(Exit(0))

	var orgEntity schema.OrgResource
	err := json.Unmarshal(response.Out.Contents(), &orgEntity)
	if err != nil {
		return "", err
	}

	return orgEntity.Resources[0].Links.Quota.Href, nil
}
