package nimble

type renderClient interface {
	Init(width, height int32) // Inform client of window size
	Render(pm PixMap)
}

var renderClientList []renderClient

func AddRenderClient(r renderClient) {
	renderClientList = append(renderClientList, r)
}

type WindowSpec interface {
	Size() (width, height int32)
	Title() string
}
