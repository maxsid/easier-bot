package contentType

type ContentType string

const (
	Document          ContentType = "Document"
	Audio             ContentType = "Audio"
	Animation         ContentType = "Animation"
	Game              ContentType = "Game"
	Photo             ContentType = "Photo"
	Sticker           ContentType = "Sticker"
	Video             ContentType = "Video"
	VideoNote         ContentType = "VideoNote"
	Voice             ContentType = "Voice"
	Contact           ContentType = "Contact"
	Location          ContentType = "Location"
	Venue             ContentType = "Venue"
	PinnedMessage     ContentType = "PinnedMessage"
	Invoice           ContentType = "Invoice"
	SuccessfulPayment ContentType = "SuccessfulPayment"
	PassportData      ContentType = "PassportData"
	Any               ContentType = "Any"
	NoContent         ContentType = ""
)

var AllContentTypes = map[ContentType]struct{}{Document: {}, Audio: {}, Animation: {}, Game: {}, Photo: {},
	Sticker: {}, Video: {}, VideoNote: {}, Voice: {}, Contact: {}, Location: {}, Venue: {}, PinnedMessage: {},
	Invoice: {}, SuccessfulPayment: {}, PassportData: {}}
