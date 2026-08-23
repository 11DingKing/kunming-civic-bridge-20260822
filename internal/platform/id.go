package platform

import "github.com/google/uuid"

func ID() string { return uuid.NewString() }
