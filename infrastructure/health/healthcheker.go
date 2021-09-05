package health

type Health string

const Healthy Health = "Healthy"
const UnHealthy Health = "Unhealthy"

const defaultTimeStampFormat = "00:00:00.0000000"

var checker CheckerService

type Checker interface {
	Check() HealthResultDto
	Tag() string
}

type CheckerService interface {
	Add(checker Checker) CheckerService
	GetHealth() map[string]interface{}
}

type checkerService struct {
	checkers []Checker
}

func (s *checkerService) Add(checker Checker) CheckerService {
	s.checkers = append(s.checkers, checker)
	return s
}

type healthResult struct {
	tag    string
	health HealthResultDto
}

func newHealthResult(tag string, health HealthResultDto) healthResult {
	return healthResult{
		tag:    tag,
		health: health,
	}
}

func (s *checkerService) GetHealth() map[string]interface{} {
	entries := map[string]HealthResultDto{}
	resultObject := map[string]interface{}{
		"status":        Healthy,
		"totalDuration": defaultTimeStampFormat,
	}
	checkersCount := len(s.checkers)
	resultChan := make(chan healthResult, checkersCount)
	for _, checker := range s.checkers {
		go func(result chan healthResult, checkerUnit Checker) {
			health := checkerUnit.Check()
			result <- newHealthResult(checkerUnit.Tag(), health)
		}(resultChan, checker)
	}

	for i := 0; i < checkersCount; i++ {
		result := <-resultChan
		if result.health.Status == UnHealthy {
			resultObject["status"] = UnHealthy
		}
		entries[result.tag] = result.health
	}
	resultObject["entries"] = entries
	return resultObject
}

func NewCheckerService() CheckerService {
	checker = &checkerService{
		checkers: []Checker{},
	}
	return checker
}

func GetHealth() map[string]interface{} {
	return checker.GetHealth()
}

func Add(newChecker Checker) CheckerService {
	return checker.Add(newChecker)
}
