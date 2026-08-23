package domain

import (
	"fmt"
	"strings"
)

type OfflineIntake struct {
	PointID, OperatorID, AuthorPhone, Title, Body, Scope string
	Consent                                              bool
}

func (i OfflineIntake) Validate() error {
	for k, v := range map[string]string{"point": i.PointID, "operator": i.OperatorID, "phone": i.AuthorPhone, "title": i.Title, "body": i.Body, "scope": i.Scope} {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("%w: %s", ErrInvalid, k)
		}
	}
	if !i.Consent {
		return fmt.Errorf("%w: consent", ErrInvalid)
	}
	return nil
}
func (i OfflineIntake) Normalized() OfflineIntake {
	i.Title = strings.TrimSpace(i.Title)
	i.Body = strings.TrimSpace(i.Body)
	i.AuthorPhone = strings.ReplaceAll(i.AuthorPhone, " ", "")
	return i
}
