package adminapi

func handleRealms() {
	// router.GET("/realms", getCurrentRealmHandler)
	// router.GET("/realms/:slug", getRealmHandler)
}

// func getCurrentRealmHandler(c echo.Context) error {
// 	claims := jwt.ExtractClaims(c)
// 	userID, err := models.FromUUID(claims["uuid"].(string))
// 	if err != nil {
// 		c.AbortWithStatus(http.StatusNotAcceptable)
// 		return
// 	}

// 	realm, err := models.GetRealmOfUser(userID)
// 	if err != nil {
// 		fmt.Println(err)
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}

// 	c.JSON(http.StatusOK, realm)
// }

// func getRealmHandler(c echo.Context) error {

// }
