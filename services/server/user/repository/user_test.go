package repository

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

func TestFindUserByEmail(t *testing.T) {
	repo := NewUserRepository(repository.GetRepository())

	user := &model.User{
		Email: "testing",
		Hash: "test-hash",
	}

	err := repository.GetRepository().Insert(user)
	if err != nil {
		t.Error(err)
	}

	foundUser := repo.FindByEmail("testing")

	if foundUser == nil {
		t.Error("expected user not found")
	}

	if len(foundUser.ID) == 0 {
		t.Error("Failed to correctly parse the user!")
	}
}
