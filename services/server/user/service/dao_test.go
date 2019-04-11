package service

import (
	"testing"
	"services/server/test"
	"os"
	"services/server/core/repository"
	"services/server/user/model"
)

func TestMain(m *testing.M) {
	test.Setup("../../")

	// make sure to teardown on panic
	defer func() {
		if r := recover(); r != nil {
			test.Teardown()
		}
	}()

	code := m.Run()

	test.Teardown()

	os.Exit(code)
}

func TestFindUserByUsername(t *testing.T) {
	user := &model.User{
		Username: "testing",
		Hash: "test-hash",
	}

	err := repository.GetRepository().Insert(user)
	if err != nil {
		t.Error(err)
	}

	foundUser := FindUserByUsername("testing")

	if foundUser == nil {
		t.Error("expected user not found")
	}

	if len(foundUser.ID) == 0 {
		t.Error("Failed to correctly parse the user!")
	}
}
