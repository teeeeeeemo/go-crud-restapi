package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	/* nickname */
	if strings.Contains(err, "nickname") {
		return errors.New("Nickname Already Taken")
	}

	/* email */
	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	/* title */
	if strings.Contains(err, "title") {
		return errors.New("Title Already Taken")
	}

	/* hashedPassword */
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}

	/* default */
	return errors.New("Incorrect Details")

}
