package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StatusString_Undefined_ReturnsUndefined(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Undefined", Undefined.String())
}

func Test_StatusString_Upcoming_ReturnsUpcoming(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Upcoming", Upcoming.String())
}

func Test_StatusString_Live_ReturnsLive(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Live", Live.String())
}

func Test_StatusString_Archived_ReturnsArchived(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Archived", Archived.String())
}

func Test_StatusString_InvalidStatus_ReturnsStatusWithNumber(t *testing.T) {
	t.Parallel()
	invalidStatus := Status(99)
	assert.Equal(t, "Status(99)", invalidStatus.String())
}
