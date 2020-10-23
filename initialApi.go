package restapi

func (c *InitAPI) initAPI() {}

func createAPI() *InitAPI {
	c := InitAPI{}
	c.initAPI()

	return &c
}
