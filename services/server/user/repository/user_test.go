package repository

import (
	"testing"
	"services/server/test"
	"os"
	"services/server/core/repository"
	"services/server/user/model"
)

func TestMain(m *testing.M) {
	test.Setup("../../", true)

	// make sure to teardown on panic
	defer func() {
		if r := recover(); r != nil {
			test.Teardown(true)
		}
	}()

	code := m.Run()

	test.Teardown(true)

	os.Exit(code)
}

func TestFindUserByID(t *testing.T) {
	repo := NewUserRepository(repository.GetRepository())

	user := &model.User{
		Hash: "test-hash",
	}

	err := repository.GetRepository().Insert(user)
	if err != nil {
		t.Error(err)
	}

	foundUser := repo.FindByID(user.ID)

	if foundUser == nil {
		t.Error("expected user not found")
	}

	if len(foundUser.ID) == 0 {
		t.Error("Failed to correctly parse the user!")
	}
}
