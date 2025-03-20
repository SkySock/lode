package v1

import "testing"

func TestValidatePassword(t *testing.T) {
	type test struct {
		name     string
		password string
	}

	correctPasswords := []test{
		{
			name:     "Valid password with all requirements",
			password: "Passw0rd!",
		},
		{
			name:     "Valid long password with special chars",
			password: "LongSecurePass123!@#",
		},
		{
			name:     "Minimum length (8) with complex chars",
			password: "A1!aBc#d",
		},
	}

	incorrectPasswords := []test{
		{
			name:     "Too short (7 characters)",
			password: "Short1!",
		},
		{
			name:     "Missing uppercase",
			password: "alllower1!",
		},
		{
			name:     "Missing lowercase",
			password: "ALLUPPER1!",
		},
		{
			name:     "Missing digit",
			password: "NoDigitsHere!",
		},
		{
			name:     "Missing special char",
			password: "NoSpecial1",
		},
	}

	edgeCases := []test{
		{
			name:     "Only one of each required type",
			password: "A1!aaaaa",
		},
		{
			name:     "All requirements but wrong order",
			password: "1!aAaaaa",
		},
		{
			name:     "Special char first",
			password: "!Aa1aaaa",
		},
		{
			name:     "Unicode special characters",
			password: "Passw0rd#",
		},
	}

	for _, test := range correctPasswords {
		t.Run(test.name, func(t *testing.T) {
			data := SignUpRequest{
				Username: "testUser",
				Password: test.password,
				Email:    "test@example.com",
			}
			err := data.Validate()
			if err != nil {
				t.Errorf("For password '%s'\nExpected: not error", test.password)
			}
		})
	}

	for _, test := range incorrectPasswords {
		t.Run(test.name, func(t *testing.T) {
			data := SignUpRequest{
				Username: "testUser",
				Password: test.password,
				Email:    "test@example.com",
			}
			err := data.Validate()
			if err == nil {
				t.Errorf("For password '%s'\nExpected: error", test.password)
			}
		})
	}

	for _, test := range edgeCases {
		t.Run(test.name, func(t *testing.T) {
			data := SignUpRequest{
				Username: "testUser",
				Password: test.password,
				Email:    "test@example.com",
			}
			err := data.Validate()
			if err != nil {
				t.Errorf("For password '%s'\nExpected: not error", test.password)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	t.Run("Correct email", func(t *testing.T) {
		data := SignUpRequest{
			Username: "testUser",
			Password: "Passw0rd!",
			Email:    "test@example.com",
		}
		err := data.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Incorrect email", func(t *testing.T) {
		data := SignUpRequest{
			Username: "testUser",
			Password: "Passw0rd!",
			Email:    "testexample.com",
		}
		err := data.Validate()
		if err == nil {
			t.Error(err)
		}
	})
	t.Run("Incorrect email", func(t *testing.T) {
		data := SignUpRequest{
			Username: "testUser",
			Password: "Passw0rd!",
			Email:    "test@examplecom",
		}
		err := data.Validate()
		if err == nil {
			t.Error(err)
		}
	})
}
