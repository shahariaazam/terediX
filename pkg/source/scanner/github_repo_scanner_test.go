package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"teredix/pkg/resource"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// GitHubClientMock is an autogenerated mock type for the GitHubClient type
type GitHubClientMock struct {
	mock.Mock
}

// ListRepositories provides a mock function with given fields: ctx, user
func (_m *GitHubClientMock) ListRepositories(ctx context.Context, user string) ([]*github.Repository, error) {
	ret := _m.Called(ctx, user)

	var r0 []*github.Repository
	if rf, ok := ret.Get(0).(func(context.Context, string) []*github.Repository); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*github.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestGitHubRepositoryScanner_Scan(t *testing.T) {
	testCases := []struct {
		name           string
		user           string
		ghRepositories []*github.Repository
		want           []resource.Resource
	}{
		{
			name: "returns resources",
			user: "testuser",
			ghRepositories: []*github.Repository{
				{
					ID:              github.Int64(123),
					Name:            github.String("testrepo"),
					FullName:        github.String("testuser/testrepo"),
					Language:        github.String("Go"),
					StargazersCount: github.Int(42),
				},
			},
			want: []resource.Resource{
				{
					Kind:       "GitHubRepository",
					UUID:       "123",
					Name:       "testrepo",
					ExternalID: "testuser/testrepo",
					MetaData: []resource.MetaData{
						{Key: "Language", Value: "Go"},
						{Key: "Stars", Value: "42"},
					},
				},
			},
		},
		{
			name:           "returns empty resource list on error",
			user:           "testuser",
			ghRepositories: []*github.Repository{},
			want:           []resource.Resource{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(GitHubClientMock)
			mockClient.On("ListRepositories", mock.Anything, mock.Anything).Return(tc.ghRepositories, nil)

			s := NewGitHubRepositoryScanner("test", mockClient, tc.user)
			got := s.Scan()

			assert.Equal(t, len(tc.ghRepositories), len(got))
		})
	}
}

func TestNewGitHubRepositoryClient_ListRepositories_Return_Data(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a list of repositories
		repos := []*github.Repository{
			{Name: github.String("repo1")},
			{Name: github.String("repo2")},
			{Name: github.String("repo3")},
		}
		jsonBytes, _ := json.Marshal(repos)
		fmt.Fprintln(w, string(jsonBytes))
	}))
	defer ts.Close()

	client, _ := github.NewEnterpriseClient(ts.URL, "", ts.Client())
	gc := NewGitHubRepositoryClient(client)
	repositories, err := gc.ListRepositories(ctx, "HI")

	assert.NoError(t, err)
	assert.Equal(t, 3, len(repositories))
}

func TestNewGitHubRepositoryClient_ListRepositories_Bad_Response_Code(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	client, _ := github.NewEnterpriseClient(ts.URL, "", ts.Client())
	gc := NewGitHubRepositoryClient(client)
	_, err := gc.ListRepositories(ctx, "HI")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list repositories for user HI")
}