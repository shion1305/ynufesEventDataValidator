package model

import (
	"errors"
	"github.com/gnue/go-disp_width"
	"net/http"
	"regexp"
)

type EventData struct {
	originOrg         string
	eventIdMD5        ID
	iconDataId        string
	eventTitle        string
	eventSummary      string
	eventDescription  string
	eventDescriptionP string
	eventGenre        EventGenre
	orgName           string
	orgDescription    string
	snsTwitter        verificationField
	snsFacebook       verificationField
	snsInstagram      verificationField
	snsWebsite        verificationField
	contactAddress    string
	originalBuilder   *EventDataBuilder
}

type verificationField struct {
	Value       string
	OriginValue string
	Verified    VerificationStatus
	Status      Status
}

type VerificationStatus string

const (
	Verified   VerificationStatus = "Verified"
	Unverified VerificationStatus = "Unverified"
	Error      VerificationStatus = "Error"
)

func (e *verificationField) setVerification(vStatus VerificationStatus) {
	e.Verified = vStatus
}

func (e *verificationField) setStatus(value string, status Status) {
	e.Status = status
	e.Value = value
}

func (e *verificationField) getSafeValue() string {
	if e.Status != NG && e.Verified == Verified {
		return e.Value
	}
	return ""
}

func (e *verificationField) getCheckString() string {
	if e.OriginValue == "" {
		return "(設定なし)"
	}
	if e.Status != NG && e.Verified == Verified {
		return e.Value + " (確認済み)"
	}
	if e.Status != NG && e.Verified == Unverified {
		return e.Value + " (未確認)"
	}
	return "(エラー・設定なし 無効な入力: " + e.OriginValue + ")"

}

type EventField string

const (
	EventTitle       EventField = "eventTitle"
	EventDescription EventField = "eventDescription"
	EventGenreF      EventField = "eventGenre"
	OrgName          EventField = "orgName"
	OrgDescription   EventField = "orgDescription"
	SnsTwitter       EventField = "snsTwitter"
	SnsFacebook      EventField = "snsFacebook"
	SnsInstagram     EventField = "snsInstagram"
	SnsWebsite       EventField = "snsWebsite"
	ContactAddress   EventField = "contactAddress"
)

type EventGenre string

const (
	Exhibition       EventGenre = "展示・体験・販売"
	Performance      EventGenre = "パフォーマンス"
	GameSports       EventGenre = "ゲーム・スポーツ"
	Dessert          EventGenre = "デザート"
	NoodleTeppanyaki EventGenre = "鉄板・麺類"
	FastFood         EventGenre = "ファストフード"
	Drink            EventGenre = "ドリンク"
	RiceDish         EventGenre = "ご飯もの"
)

func (genre EventGenre) getEventGenreId() int {
	switch genre {
	case Exhibition:
		return 1
	case Performance:
		return 2
	case GameSports:
		return 3
	case Dessert:
		return 4
	case NoodleTeppanyaki:
		return 5
	case FastFood:
		return 6
	case Drink:
		return 7
	case RiceDish:
		return 8
	default:
		return 0
	}
}

type Status string

const (
	OK      Status = "OK"
	Changed Status = "フォーマットの変更がされました"
	Warning Status = "確認が必要な変更がされました"
	NG      Status = "不正な値です"
)

func NewEventData(builder EventDataBuilder) *EventData {
	var newData EventData
	newData.originOrg = builder.OriginOrg
	newData.eventIdMD5 = genID(builder.OriginOrg)
	newData.iconDataId = getIconId(builder.IconDataId)
	newData.eventTitle = builder.EventTitle
	newData.eventSummary = builder.EventSummary
	newData.eventDescription = builder.EventDescription
	newData.eventDescriptionP = builder.EventDescriptionP
	newData.eventGenre = EventGenre(builder.EventGenreText)
	newData.orgName = builder.OrgName
	newData.orgDescription = builder.OrgDescription
	newData.snsTwitter = initVerificationField(builder.SnsTwitter)
	newData.snsFacebook = initVerificationField(builder.SnsFacebook)
	newData.snsInstagram = initVerificationField(builder.SnsInstagram)
	newData.snsWebsite = initVerificationField(builder.SnsWebsite)
	newData.contactAddress = builder.ContactAddress
	newData.originalBuilder = &builder
	newData.validate()
	return &newData
}

func initVerificationField(value string) verificationField {
	var resp verificationField
	resp.OriginValue = value
	resp.Value = value
	resp.Verified = Unverified
	return resp
}

func NewMultiEventData(builders []EventDataBuilder) []*EventData {
	var data []*EventData
	for _, builder := range builders {
		data = append(data, NewEventData(builder))
	}
	return data
}

func (e *EventData) UpdateField(field EventField, value string) error {
	switch field {
	case EventTitle:
		e.eventTitle = value
		break
	case EventDescription:
		e.eventDescription = value
		break
	case EventGenreF:
		e.eventGenre = EventGenre(value)
		break
	case OrgName:
		e.orgName = value
		break
	case OrgDescription:
		e.orgDescription = value
		break
	case SnsTwitter:
		e.snsTwitter.Value = value
		break
	case SnsFacebook:
		e.snsFacebook.Value = value
		break
	case SnsInstagram:
		e.snsInstagram.Value = value
		break
	case SnsWebsite:
		e.snsWebsite.Value = value
		break
	case ContactAddress:
		e.contactAddress = value
		break
	default:
		return errors.New("unknown Field")
	}
	e.validate()
	return nil
}

func (e *EventData) validate() {
	//_, s1 := e.validateEventTitle()
	//_, s1 := e.validateEventDescription()
	//_, s1 := e.validateEventGenreText()
	//_, s1 := e.validateOrgName()
	//_, s1 := e.validateOrgDescription()
	e.ValidateSnsTwitter()
	//_, s1 := e.validateSnsInstagram()
	//_, s1 := e.validateSnsFacebook()
	//_, s2 := e.validateSnsWebsite()
	//if s2 == NG {
	//	fmt.Println(e.snsWebsite)
	//}
}

func validAsID(s string) string {
	re := regexp.MustCompile(`^@?((\w){1,15}) *$`)
	id := re.FindStringSubmatch(s)
	if id == nil {
		return ""
	}
	return id[1]
}

//func (e *EventData) validateEventTitle() (string, Status) {
//
//	return "", OK
//}
//
//func (e *EventData) validateEventDescription() (string, Status) {
//	return "", OK
//}
//
//func (e *EventData) validateEventGenreText() (string, Status) {
//	return "", OK
//}
//
//func (e *EventData) validateOrgName() (string, Status) {
//	return "", OK
//}
//
//func (e *EventData) validateOrgDescription() {
//	return "", OK
//}

func (e *EventData) ValidateDescriptionP() bool {
	//limit EventDescriptionP visual text length to 60(with half size character)
	return disp_width.Measure(e.eventDescriptionP) <= 60
}

func (e *EventData) ValidateSnsTwitter() {
	if e.snsTwitter.Value == "" {
		e.snsTwitter.setStatus("", OK)
		return
	}
	if id := validAsID(e.snsTwitter.Value); id != "" {
		e.snsTwitter.setStatus(id, OK)
		return
	}
	re := regexp.MustCompile("^https://twitter.com/([A-Za-z0-9]*_?[A-Za-z0-9]*)")
	if id := re.FindStringSubmatch(e.snsTwitter.Value); id != nil {
		if name := validAsID(id[1]); name != "" {
			e.snsTwitter.setStatus(name, Changed)
			return
		}
	}
	e.snsTwitter.setStatus(e.snsTwitter.Value, NG)
	return
}

func (e *EventData) validateSnsFacebook() {
	if e.snsFacebook.Value == "" {
		e.snsFacebook.setStatus("", OK)
	}
	re := regexp.MustCompile("^@?([A-Z.a-z0-9]{5,})$")
	check := re.FindStringSubmatch(e.snsFacebook.Value)
	if check != nil && check[1] != "" {
		e.snsFacebook.setStatus(check[1], OK)
		return
	}
	e.snsFacebook.setStatus("", NG)
	return
}

func (e *EventData) validateSnsInstagram() {
	if e.snsInstagram.Value == "" {
		e.snsInstagram.setStatus("", OK)
		return
	}
	re := regexp.MustCompile("^@?([A-Z.a-z0-9_]+)$")
	check := re.FindStringSubmatch(e.snsInstagram.Value)
	if check != nil && check[1] != "" {
		e.snsInstagram.setStatus(check[1], OK)
		return
	}
	e.snsInstagram.setStatus("", NG)
	return
}

//func (e *EventData) validateRequestInstagram() {
//	resp, err := http.Get("https://www.instagram.com/usernamafawefwaefaefawefawe")
//	fmt.Println(err)
//	fmt.Println(resp)
//}

func (e *EventData) validateSnsWebsite() {
	if e.snsWebsite.Value == "" {
		e.snsWebsite.setStatus("", OK)
		return
	}
	re := regexp.MustCompile("^(https?://[/.a-z0-9_-]+)$")
	if match := re.FindStringSubmatch(e.snsWebsite.Value); match != nil {
		if accessTest(match[0]) {
			e.snsWebsite.setStatus(match[0], OK)
			e.snsWebsite.setVerification(Verified)
		} else {
			e.snsWebsite.setStatus(match[0], OK)
			e.snsWebsite.setVerification(Error)
		}
		return
	}
	e.snsWebsite.setStatus(e.snsWebsite.Value, NG)
	return
}

func accessTest(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		print(err)
		return false
	}
	if resp.StatusCode > 300 {
		return false
	}
	return true
}

// WEB用のExporter
// 必要のないデータは含まない
type ExportEventData struct {
	EventIdMD5       string `json:"event_id"`
	EventTitle       string `json:"event_title"`
	EventSummary     string `json:"event_summary"`
	EventDescription string `json:"event_description"`
	EventGenreId     int    `json:"event_genre_id"`
	OrgName          string `json:"org_name"`
	OrgDescription   string `json:"org_description"`
	SnsTwitter       string `json:"sns_twitter"`
	SnsFacebook      string `json:"sns_facebook"`
	SnsInstagram     string `json:"sns_instagram"`
	SnsWebsite       string `json:"sns_website"`
}

func (e *EventData) Export() ExportEventData {
	return ExportEventData{
		EventIdMD5:       string(e.eventIdMD5),
		EventTitle:       e.eventTitle,
		EventDescription: e.eventDescription,
		EventGenreId:     e.eventGenre.getEventGenreId(),
		OrgName:          e.orgName,
		OrgDescription:   e.orgDescription,
		SnsTwitter:       e.snsTwitter.getSafeValue(),
		SnsFacebook:      e.snsFacebook.getSafeValue(),
		SnsInstagram:     e.snsInstagram.getSafeValue(),
		SnsWebsite:       e.snsWebsite.getSafeValue(),
	}
}

type CheckEventData struct {
	OriginOrg         string `csv:"OriginOrg"`
	ContactAddress    string `csv:"ContactAddress"`
	Url               string `csv:"Url"`
	EventTitle        string `csv:"eventTitle"`
	EventSummary      string `csv:"eventSummary"`
	EventDescription  string `csv:"eventDescription"`
	EventDescriptionP string `csv:"eventDescriptionP"`
	EventGenreText    string `csv:"eventGenreText"`
	OrgName           string `csv:"orgName"`
	OrgDescription    string `csv:"orgDescription"`
	SnsTwitter        string `csv:"snsTwitter"`
	SnsFacebook       string `csv:"snsFacebook"`
	SnsInstagram      string `csv:"snsInstagram"`
	SnsWebsite        string `csv:"snsWebsite"`
	ImageComment      string `csv:"imageComment"`
}

func (e *EventData) ExportCheck() *CheckEventData {
	return &CheckEventData{
		OriginOrg:         e.originOrg,
		ContactAddress:    e.contactAddress,
		Url:               "https://tokiwa22.ynu-fes.yokohama/preview/event-detail/" + string(e.eventIdMD5),
		EventTitle:        e.eventTitle,
		EventSummary:      e.eventSummary,
		EventDescription:  e.eventDescription,
		EventDescriptionP: e.eventDescriptionP,
		EventGenreText:    string(e.eventGenre),
		OrgName:           e.orgName,
		OrgDescription:    e.orgDescription,
		SnsTwitter:        e.snsTwitter.getCheckString(),
		SnsFacebook:       e.snsFacebook.getCheckString(),
		SnsInstagram:      e.snsInstagram.getCheckString(),
		SnsWebsite:        e.snsWebsite.getCheckString(),
	}
}
