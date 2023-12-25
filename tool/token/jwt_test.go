package token

import (
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	a := NewJwtMgr("1", "2")
	d := Data{
		User:       "test",
		Permission: []string{"read", "write"},
		ValidTime:  time.Now().Add(time.Hour * 24 * 7).Unix(),
		token:      "",
	}
	a.Signature(&d)

	tests := []struct {
		name       string
		time       int64
		permission string
		want       bool
	}{
		{
			name:       "right permission, in time",
			time:       time.Now().Add(time.Hour * 24 * 6).Unix(),
			permission: "read",
			want:       true,
		},
		{
			name:       "right permission, in time",
			time:       time.Now().Add(time.Hour * 24 * 6).Unix(),
			permission: "write",
			want:       true,
		},
		{
			name:       "wrong permission, in time",
			time:       time.Now().Add(time.Hour * 24 * 6).Unix(),
			permission: "delete",
			want:       false,
		},
		{
			name:       "right permission, out time",
			time:       time.Now().Add(time.Hour * 24 * 8).Unix(),
			permission: "read",
			want:       false,
		},
		{
			name:       "right permission, out time",
			time:       time.Now().Add(time.Hour * 24 * 8).Unix(),
			permission: "write",
			want:       false,
		},
		{
			name:       "wrong permission, out time",
			time:       time.Now().Add(time.Hour * 24 * 8).Unix(),
			permission: "delete",
			want:       false,
		},
	}

	for _, test := range tests {
		t2 := time.Unix(test.time, 0)
		t.Run(test.name, func(ttt *testing.T) {
			if got := a.CheckSignature(&d, t2, test.permission); got != test.want {
				ttt.Errorf("Check(%v, %v) = %v, want %v\n", d.token, test.permission, got, test.want)
			} else {
				ttt.Logf("Check(%v, %v) = %v, want %v\n", d.token, test.permission, got, test.want)
			}
		})
	}

	d.Permission = append(d.Permission, "delete")
	t.Run("fake permission, in time", func(ttt *testing.T) {
		if got := a.CheckSignature(&d, time.Now().Add(time.Hour*24*6), "delete"); got != false {
			ttt.Errorf("Check(%v, %v) = %v, want %v\n", d.token, "fake", got, false)
		} else {
			ttt.Logf("Check(%v, %v) = %v, want %v\n", d.token, "fake", got, false)
		}
	})
}
