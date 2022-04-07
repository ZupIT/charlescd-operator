package client_test

import (
	"context"
	gitClient "github.com/ZupIT/charlescd-operator/internal/client"
	"github.com/ZupIT/charlescd-operator/internal/client/mocks"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Gitrepository", func() {
	var mockClient *mocks.Client
	var gitRepository *gitClient.GitRepository

	Context(" when  git repository exists on cluster", func() {
		It("should return it successfully", func() {
			mockClient = new(mocks.Client)
			gitRepository = gitClient.NewGitRepository(mockClient)
			key := client.ObjectKey{Namespace: "test", Name: "gitRepository"}
			expectedGit := new(sourcev1beta1.GitRepository)
			mockClient.On("Get", mock.Anything, key, mock.Anything).Return(nil)
			returnedGit, err := gitRepository.GetGitRepository(context.TODO(), key)
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedGit).To(Equal(expectedGit))
		})

		It("should return error when fails to get it", func() {
			statusError := errors.NewBadRequest("bad request, check your payload")
			mockClient = new(mocks.Client)
			gitRepository = gitClient.NewGitRepository(mockClient)
			key := client.ObjectKey{Namespace: "test", Name: "gitRepository"}
			mockClient.On("Get", mock.Anything, key, mock.Anything).Return(statusError)
			returnedGit, err := gitRepository.GetGitRepository(context.TODO(), key)
			Expect(err.Error()).To(ContainSubstring(statusError.Error()))
			Expect(returnedGit).To(BeNil())
		})
	})
	Context(" when  git repository does not exists on cluster", func() {
		It("should not return error", func() {
			statusError := errors.NewNotFound(schema.GroupResource{Group: "apps", Resource: "deployments"}, "test")
			mockClient = new(mocks.Client)
			gitRepository = gitClient.NewGitRepository(mockClient)
			key := client.ObjectKey{Namespace: "test", Name: "gitRepository"}
			mockClient.On("Get", mock.Anything, key, mock.Anything).Return(statusError)
			returnedGit, err := gitRepository.GetGitRepository(context.TODO(), key)
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedGit).To(BeNil())
		})
	})

})
