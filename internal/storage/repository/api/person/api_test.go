package person

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetNameInfo(t *testing.T) {
	repo := NewRepository()

	resp, err := repo.GetNameInfo("Helen")

	require.NoError(t, err)
	require.NotNil(t, resp)

	if assert.NotNil(t, resp.Agify) {
		assert.NotNil(t, resp.Agify.Count)
		assert.NotNil(t, resp.Agify.Name)
		assert.NotNil(t, resp.Agify.Age)

		assert.NotEqual(t, resp.Agify.StatusCode, 0)
	}

	if assert.NotNil(t, resp.Genderize) {
		assert.NotNil(t, resp.Genderize.Name)
		assert.NotNil(t, resp.Genderize.Count)
		assert.NotNil(t, resp.Genderize.Probability)
		assert.NotNil(t, resp.Genderize.Gender)

		assert.NotEqual(t, resp.Genderize.StatusCode, 0)
	}

	if assert.NotNil(t, resp.Nationalize) {
		assert.NotNil(t, resp.Nationalize.Name)
		assert.NotNil(t, resp.Nationalize.Count)
		assert.NotNil(t, resp.Nationalize.Country)

		assert.NotEqual(t, resp.Nationalize.StatusCode, 0)
	}
}

func TestGetAge(t *testing.T) {
	repo := NewRepository()

	resp, err := repo.GetAge("Helen")

	require.NoError(t, err)
	require.NotNil(t, resp)

	if assert.NotNil(t, resp) {
		assert.NotNil(t, resp.Count)
		assert.NotNil(t, resp.Name)
		assert.NotNil(t, resp.Age)

		assert.NotEqual(t, resp.StatusCode, 0)
	}
}

func TestGetGender(t *testing.T) {
	repo := NewRepository()

	resp, err := repo.GetGender("Helen")

	require.NoError(t, err)
	require.NotNil(t, resp)

	if assert.NotNil(t, resp) {
		assert.NotNil(t, resp.Name)
		assert.NotNil(t, resp.Count)
		assert.NotNil(t, resp.Probability)
		assert.NotNil(t, resp.Gender)

		assert.NotEqual(t, resp.StatusCode, 0)
	}
}

func TestGetNation(t *testing.T) {
	repo := NewRepository()

	resp, err := repo.GetNation("Helen")

	require.NoError(t, err)
	require.NotNil(t, resp)

	if assert.NotNil(t, resp) {
		assert.NotNil(t, resp.Name)
		assert.NotNil(t, resp.Count)
		assert.NotNil(t, resp.Country)

		assert.NotEqual(t, resp.StatusCode, 0)
	}
}
