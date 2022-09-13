package model

type EventData struct {
	eventTitle       string
	eventDescription string
	eventGenreText   string
	orgName          string
	orgDescription   string
	snsTwitter       string
	snsFacebook      string
	snsInstagram     string
	snsWebsite       string
}

type ValidationErrors struct {
	err ValidationError
}
type ValidationError struct {
	Field   string
	Message string
}

type Status string

const (
	OK       Status = "OK"
	Changed  Status = "フォーマットの変更がされました"
	Selected Status = "複数の候補から選択しました"
	NG       Status = "不正な値です"
)

func NewEventData(builder EventDataBuilder) *EventData {
	var newData EventData
	newData.eventTitle = builder.EventTitle
	newData.eventDescription = builder.EventDescription
	newData.eventGenreText = builder.EventGenreText
	newData.orgName = builder.OrgName
	newData.orgDescription = builder.OrgDescription
	newData.snsTwitter = builder.SnsTwitter
	newData.snsFacebook = builder.SnsFacebook
	newData.snsInstagram = builder.SnsInstagram
	newData.snsWebsite = builder.SnsWebsite
	return &newData
}

func (e *EventData) Validate() {

}

func (e *EventData) validateEventTitle() ValidationError {
	return ValidationError{}
}

func (e *EventData) validateEventDescription() (string, Status) {
	return "", OK
}

func (e *EventData) validateEventGenreText() (string, Status) {
	return "", OK
}

func (e *EventData) validateOrgName() (string, Status) {
	return "", OK
}

func (e *EventData) validateOrgDescription() (string, Status) {
	return "", OK
}

func (e *EventData) validateSnsTwitter() (string, Status) {
	return "", OK
}

func (e *EventData) validateSnsFacebook() (string, Status) {
	return "", OK
}
func (e *EventData) validateSnsInstagram() (string, Status) {
	return "", OK
}
func (e *EventData) validateSnsWebsite() (string, Status) {
	return "", OK
}
