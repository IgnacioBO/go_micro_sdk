package mock

import (
	"testing"

	"github.com/IgnacioBO/go_micro_sdk/user"
)

func TestMock_User(t *testing.T) {
	t.Run("should implement user.Transport interface", func(t *testing.T) {
		var _ user.Transport = (*UserSdkMock)(nil)
	})

}
