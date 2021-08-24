package statics

type D2Statics struct {
	key string
}

func New(apiKey string) *D2Statics {
	return &D2Statics{apiKey}
}

func (d *D2Statics) MatchInfo(id string) {

}