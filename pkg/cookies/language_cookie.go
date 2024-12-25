package cookies

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

const LANGUAGE_COOKIE_NAME = "lang"

func SetLanguageCookie(c echo.Context, lang string) error {
	AVAILABLE_LANGUAGES := []string{"de", "en"}
	if !slices.Contains(AVAILABLE_LANGUAGES, lang) {
		return fmt.Errorf("Language not supported")
	}

	languageCookie := &http.Cookie{
		Name:  LANGUAGE_COOKIE_NAME,
		Value: lang,
		Path:  "/",
	}
	c.SetCookie(languageCookie)

	return nil
}

func GetLanguage(c echo.Context) (string, error) {
	languageCookie, err := c.Cookie(LANGUAGE_COOKIE_NAME)
	if err != nil {
		return "", nil
	}

	return languageCookie.Value, nil
}
