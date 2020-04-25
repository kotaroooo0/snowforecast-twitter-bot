package responses

type PostKanjiToHiraganaResponse struct {
	MaResult MaResult `xml:"ma_result" json:"ma_result"`
}

type MaResult struct {
	WordList []Word `xml:"word_list"`
}

type Word struct {
	Reading string `xml:"word"`
}

func NewPostKanjiToHiraganaResponse() PostKanjiToHiraganaResponse {
	return PostKanjiToHiraganaResponse{}
}
