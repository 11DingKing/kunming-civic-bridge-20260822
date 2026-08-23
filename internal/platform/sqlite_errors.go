package platform

import (
	"errors"

	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	sqlite "modernc.org/sqlite"
)

// busyConflictCodes are SQLite result codes that surface when a writer races
// another active transaction on the same row (SQLITE_BUSY) or when the
// unlock-notify machinery detects an actual deadlock between concurrent
// writers (SQLITE_LOCKED). busy_timeout cannot recover from the deadlock
// variant, so these errors must be translated rather than retried.
var busyConflictCodes = map[int]bool{5: true, 6: true} // SQLITE_BUSY, SQLITE_LOCKED

// TranslateSQLiteBusy maps SQLite's concurrent-writer lock/deadlock errors onto
// domain.ErrConflict. It returns the original error unchanged for anything else
// so genuine database failures stay visible. Callers that take a row lease
// (claim workflows) wrap their UPDATE with this so a losing concurrent request
// receives an explicit conflict instead of a raw lock error.
func TranslateSQLiteBusy(e error) error {
	if e == nil {
		return nil
	}
	var se *sqlite.Error
	if errors.As(e, &se) && busyConflictCodes[se.Code()] {
		return domain.ErrConflict
	}
	return e
}
