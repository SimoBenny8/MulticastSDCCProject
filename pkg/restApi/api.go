package restApi

import (
	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

func NewRoutes() routes {
	r := routes{
		router: gin.Default(),
	}

	//v1 := r.router.Group("/v1")

	//r.addPing(v1)
	//r.addUsers(v1)

	return r
}

func (r routes) Run(addr ...string) error {
	return r.router.Run()
}

func main() {
	router := gin.Default()
	router.GET("/groups/{id}/:messages", getMessages) //get messages of a group
	//router.GET("/groups", getInfoGroup)//get info group
	//router.GET("/albums/:id", getAlbumByID)
	//router.POST("/groups", addGroup)//add a new group
	router.POST("/groups/{id}/:messages", sendMessage) //send a message to a group
	//router.DELETE("/groups/{id}",closeGroup)//close group
	//router.PUT("/groups/{id}",startGroup)//start a new group

	err := router.Run("8080")
	if err != nil {
		return
	}
}

func getMembersGroup(c *gin.Context) {

}

func getMessages(c *gin.Context) {

}

func getDeliverQueue(c *gin.Context) {

}

func getInfoGroup(c *gin.Context) {

	groups := make([]*MulticastInfo, 0)

	for _, group := range MulticastGroups {
		group.groupMu.RLock()
		groups = append(groups, group.Group)
		group.groupMu.RUnlock()
	}

}

func addGroup(c *gin.Context) {

}

func sendMessage(c *gin.Context) {

}

func closeGroup(c *gin.Context) {

}

func startGroup(c *gin.Context) {

}

// getAlbums responds with the list of all albums as JSON.
/*func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}*/
