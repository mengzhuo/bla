package bla

type mulDocData struct {
	Hdl   *Handler
	Title string
	Docs  []*Doc
}

type singleData struct {
	Hdl   *Handler
	Title string
	Doc   *Doc
}

type tagData struct {
	Hdl     *Handler
	Title   string
	Docs    []*Doc
	TagName string
}
