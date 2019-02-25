package changelogutils

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/githubutils"

	"github.com/solo-io/solo-kit/pkg/utils/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("ChangelogTest", func() {

	expectGetProposedTag := func(latestTag, changelogDir, tag, err string) {
		actualTag, actualErr := GetProposedTagLocal(latestTag, changelogDir)
		Expect(actualTag).To(BeEquivalentTo(tag))
		if err == "" {
			Expect(actualErr).To(BeNil())
		} else {
			Expect(actualErr.Error()).To(BeEquivalentTo(err))
		}
	}

	It("works", func() {
		tmpDir := mustWriteTestDir()
		defer os.RemoveAll(tmpDir)
		changelogDir := filepath.Join(tmpDir, ChangelogDirectory)
		Expect(os.Mkdir(changelogDir, 0700)).To(BeNil())
		Expect(createSubdirs(changelogDir, "v0.0.1", "v0.0.2", "v0.0.3", "v0.0.4")).To(BeNil())
		expectGetProposedTag("v0.0.3", tmpDir, "v0.0.4", "")
	})

	It("works with git", func() {
		ctx := context.Background()
		client, err := githubutils.GetClient(ctx)
		Expect(err).NotTo(HaveOccurred())
		_, err = GetProposedTagFromGit(ctx, client, "solo-io", "testrepo", "v0.0.16")
		Expect(err).NotTo(HaveOccurred())
		_, err = ValidateProposedChangelogTag(ctx, client, "solo-io", "testrepo", "v0.0.16")
		Expect(err).NotTo(HaveOccurred())
	})
})

func createSubdirs(dir string, names ...string) error {
	for _, name := range names {
		subdir := filepath.Join(dir, name)
		err := os.Mkdir(subdir, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func mustWriteTestDir() string {
	tmpDir, err := ioutil.TempDir("", "changelog-test-")
	if err != nil {
		log.Fatalf("Failed to create test dir: %v", err)
	}
	return tmpDir
}