package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary		SignIn LDAP
// @Tags			auth
// @Description	sign-in LDAP method
// @ID				login-employee-LDAP
// @Accept			json
// @Produce		json
// @Param			input	body		signInInput	true	"credentials"
// @Success		200		{object}	response.okResponse
// @Failure		400,404	{object}	response.errorResponse
// @Failure		500		{object}	response.errorResponse
// @Failure		default	{object}	response.errorResponse
// @Router			/auth/auth [post]
func (h *Handler) authenticate(c *gin.Context) {
	lg := sl.SetupLogger()
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid input body")
		return
	}

	mail := extractLoginFromEmailOrUseLogin(input.Email)

	ldapUser, err := h.Services.Authorization.VerifyCredentials(mail, input.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.Services.Authorization.SyncUserDataFromLDAP(*ldapUser)
	if err != nil {
		lg.Error("Error sync user from LDAP", sl.Err(err))
		return
	}

	token, err := h.Services.Authorization.GenerateTokenWithLDAP(mail, input.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "employee not found")
		return
	}

	lg.Info(mail + " sign-in")

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", token, 3600*12, "", "", true, false)
	c.JSON(http.StatusOK, gin.H{"message": "Sign-in successful"})
}
