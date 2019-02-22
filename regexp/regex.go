package regexp

import "regexp"

func AlphaNumericOnly(text, replaceby string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(text, replaceby), nil
}
