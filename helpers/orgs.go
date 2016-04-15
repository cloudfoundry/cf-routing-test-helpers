package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/cf-routing-test-helpers/schema"
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
)

func GetOrgQuotaDefinitionUrl(orgGuid string, timeout time.Duration) (string, error) {
	orgGuid = strings.TrimSuffix(orgGuid, "\n")
	orgResponse := runner.NewCmdRunner(
		cf.Cf("curl", fmt.Sprintf("/v2/organizations/%s", string(orgGuid))),
		timeout).Run().Out.Contents()

	var orgEntity schema.OrgResource
	err := json.Unmarshal(orgResponse, &orgEntity)
	if err != nil {
		return "", err
	}

	return orgEntity.Entity.QuotaDefinitionUrl, nil
}
