package teamcity

import (
	"net/http"
	"strconv"

	"github.com/dghubble/sling"
)

// QueuedBuild represents a build that is on the Queue
// all available fields here: https://www.jetbrains.com/help/teamcity/rest/build.html#Properties
type QueuedBuild struct {
	ID          int64              `json:"id,omitempty"`
	BuildTypeID string             `json:"buildTypeId,omitempty"`
	Status      string             `json:"status,omitempty"`
	StatusText  string             `json:"statusText,omitempty"`
	State       string             `json:"state,omitempty"`
	Href        string             `json:"href,omitempty"`
	WebURL      string             `json:"webUrl,omitempty"`
	BuildType   BuildTypeReference `json:"buildType,omitempty"`
	WaitReason  string             `json:"waitReason,omitempty"`
	QueuedDate  string             `json:"queuedDate,omitempty"`
	Triggered   QueueTriggered     `json:"triggered,omitempty"`
}

// return build ID as a string
func (e QueuedBuild) GetID() string {
	return strconv.FormatInt(int64(e.ID), 10)
}

// TriggerBuildRequest represents parameters to put a build in queue
type TriggerBuildRequest struct {
	BuildTypeID string      `json:"buildTypeId"`
	BranchName  string      `json:"branchName,omitempty"`
	Properties  *Properties `json:"properties,omitempty"`
}

// NewTriggerBuildRequest returns a new request for triggering a build
func NewTriggerBuildRequest(buildTypeID string, props *Properties) *TriggerBuildRequest {
	return &TriggerBuildRequest{
		BuildTypeID: buildTypeID,
		Properties:  props,
	}
}

// QueueTriggered contains information about the trigger that created a queued build
type QueueTriggered struct {
	Type string `json:"type"`
	Date string `json:"date"`
}

// QueueService has operations for querying and interacting with server's build queue
type QueueService struct {
	restHelper *restHelper
}

func newQueueService(base *sling.Sling, httpClient *http.Client) *QueueService {
	sling := base.Path("buildQueue/")
	return &QueueService{
		restHelper: newRestHelperWithSling(httpClient, sling),
	}
}

// TriggerBuild will put a build in the queue with the given parameters
func (s *QueueService) TriggerBuild(req *TriggerBuildRequest) (*QueuedBuild, error) {
	var created QueuedBuild

	err := s.restHelper.post("", req, &created, "Trigger Build")

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetQueuedBuild will get a build status after queueing a build
func (s *QueueService) GetQueuedBuild(buildId string) (*QueuedBuild, error) {
	var created QueuedBuild

	// locator := LocatorID(strconv.FormatInt(int64(build.ID), 10)).String()
	locator := LocatorID(buildId).String()

	err := s.restHelper.get(locator, &created, "Get Queued Build")
	if err != nil {
		return nil, err
	}

	return &created, nil
}
