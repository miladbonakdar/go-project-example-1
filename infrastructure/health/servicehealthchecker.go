package health

import (
	"giftcard-engine/infrastructure/logger"
	"net/http"
)

type serviceHealthChecker struct {
	client   *http.Client
	endpoint string
	tag      string
}

func (c *serviceHealthChecker) Check() HealthResultDto {
	re, err := c.client.Get(c.endpoint)
	if err != nil {
		logger.WithException(err).WithDevMessage("error while checking service endpoint : " + c.endpoint).
			Error("service is not responding properly")
		return HealthResultDto{
			Status:      UnHealthy,
			Duration:    defaultTimeStampFormat,
			Exception:   err.Error(),
			Description: err.Error(),
			Data:        map[string]string{},
		}
	}
	if re.StatusCode < 200 || re.StatusCode > 260 {
		logger.WithData(map[string]interface{}{
			"statusCode": re.StatusCode,
			"endpoint":   c.endpoint,
		}).Warn("calling service endpoint resulted with invalid status code")
		return HealthResultDto{
			Status:      UnHealthy,
			Duration:    defaultTimeStampFormat,
			Exception:   "calling service endpoint resulted with invalid status code",
			Description: "calling service endpoint resulted with invalid status code",
			Data:        map[string]string{},
		}
	}
	return HealthResultDto{
		Status:   Healthy,
		Duration: defaultTimeStampFormat,
		Data:     map[string]string{},
	}
}

func (c *serviceHealthChecker) Tag() string {
	return c.tag
}

func NewServiceHealthChecker(tag, endpoint string) Checker {
	return &serviceHealthChecker{
		client:   &http.Client{},
		endpoint: endpoint,
		tag:      tag,
	}
}
