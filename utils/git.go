package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const (
	gitUrlSeparator = "#"
	gitNormalUrl    = "http"
	tagRegexp       = `\d(\..*)+`
)

const (
	repoColor   = lipgloss.Color("#f98b6c")
	urlColor    = lipgloss.Color("#4a26fd")
	tagColor    = lipgloss.Color("#ef4fa6")
	hashColor   = lipgloss.Color("#1e88e5")
	branchColor = lipgloss.Color("#03dac6")
)

func GitClone(
	repoName string,
	repoUrl string,
	repoDirPath string,
) {
	repoLog := ColorTextf("repo=%s", repoColor, repoName)
	urlLog := ColorTextf("url=%s", urlColor, repoUrl)
	log.Debugf("cloning %s %s", repoLog, urlLog)

	cleanModuleUrl, commitHash, branch, tag := parseGitUrl(repoUrl)
	reference := branch
	if tag != "" {
		reference = tag
	}

	options := &git.CloneOptions{
		URL:           cleanModuleUrl,
		ReferenceName: reference,
	}

	auth := getGitAuth(repoUrl)
	if auth != nil {
		options.Auth = auth
	}

	repo, err := git.PlainClone(repoDirPath, false, options)
	CheckErrorf(err, "error while clonning %s", repoLog)

	if commitHash != "" {
		GitCheckoutToCommit(repo, repoName, commitHash)
	}

	headName := GetHeadShortName(repo, commitHash != "", tag != "")
	headColor := getGitColor(commitHash != "", tag != "")

	headLog := ColorTextf("head=%s", headColor, headName)
	successLog := ColorSuccess("cloning successful")
	log.Debugf("%s %s %s", successLog, repoLog, headLog)
	fmt.Println()
}

func GitDirStatus(dirPath string) git.Status {
	repo, err := git.PlainOpen(dirPath)
	CheckErrorf(err, "error while trying to open directory in git (GitDirStatus): %s", dirPath)

	workTree, err := repo.Worktree()
	CheckErrorf(err, "error while getting workTree (GitDirStatus) for %s", dirPath)

	status, err := workTree.Status()
	CheckErrorf(err, "Error while trying to get git status for %s", dirPath)

	return status
}

func GitPull(dirPath string) {
	repo, err := git.PlainOpen(dirPath)
	if err != nil {
		LogReadableErrorf(err, "error while trying to open directory in git (GitPull): %s", dirPath)
		return
	}

	head, err := repo.Head()
	if err != nil {
		LogReadableErrorf(err, "error while getting repo head (GitPull) for %s", dirPath)
		return
	}

	if !head.Name().IsBranch() {
		LogWarningf("%s is not on any branch. skipping git pull", dirPath)
		return
	}

	workTree, err := repo.Worktree()
	if err != nil {
		LogReadableErrorf(err, "error while getting workTree (GitPull) for %s", dirPath)
		return
	}

	remotes, err := repo.Remotes()
	if err != nil {
		LogReadableErrorf(err, "error while getting remotes (GitPull) for %s", dirPath)
		return
	}

	remoteUrl := remotes[0].Config().URLs[0]
	branchName := head.Name().Short()

	authOptions := &git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
	}
	auth := getGitAuth(remoteUrl)
	if auth != nil {
		authOptions.Auth = auth
	}

	err = workTree.Pull(authOptions)
	isUptodateErr := errors.Is(err, git.NoErrAlreadyUpToDate)

	if err != nil && !isUptodateErr {
		LogReadableErrorf(err, "git pull error for %s", dirPath)
	} else {
		LogSuccessf("successfull git pull for %s", dirPath)
	}
}

func GitCheckoutToCommit(
	repo *git.Repository,
	repoName string,
	commitHash string,
) {
	if commitHash == "" {
		return
	}

	workTree, err := repo.Worktree()
	CheckErrorf(err, "error while getting repo %s worktree before checkout", repoName)

	err = workTree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})

	repoLog := "repo=" + repoName
	commitLog := "commitHash=" + commitHash
	CheckErrorf(err, "error while trying to checkout repo=%s commitHash=%s", repoLog, commitLog)

	repoColorLog := ColorText(repoLog, repoColor)
	commitColorLog := ColorText(commitLog, hashColor)
	log.Debugf("sucessful checkout for %s to %s", repoColorLog, commitColorLog)
}

func IsGitUrl(url string) bool {
	return strings.Contains(url, "git")
}

func GetHeadShortName(repo *git.Repository, isCommit bool, isTag bool) string {
	head, err := repo.Head()
	CheckError(err, "error while getting repo head (GetHeadShortName)")

	if isCommit {
		return head.Hash().String()
	}

	if isTag {
		headTag := GetHeadTag(repo)
		if headTag == nil {
			ThrowErrorf("couldn't find tag for head %s (GetHeadShortName)", head.String())
		}

		return headTag.Name().Short()
	}

	return head.Name().Short()
}

func GetHeadTag(repo *git.Repository) *plumbing.Reference {
	head, err := repo.Head()
	CheckError(err, "error while getting repo head (GetHeadTag)")

	tags, err := repo.Tags()
	CheckError(err, "error while getting repo tags (GetHeadTag)")

	var headTag *plumbing.Reference

	tags.ForEach(func(tag *plumbing.Reference) error {
		tagHash, err := repo.ResolveRevision(plumbing.Revision(tag.Name()))
		CheckError(err, "error while resolving revision (GetHeadTag)")

		if *tagHash == head.Hash() {
			headTag = tag
		}

		return nil
	})

	return headTag
}

// Returns cleanUrl, commitHash, branch, tag
func parseGitUrl(gitUrl string) (
	string,
	string,
	plumbing.ReferenceName,
	plumbing.ReferenceName,
) {
	if !IsGitUrl(gitUrl) {
		return "", "", "", ""
	}

	splitUrl := strings.Split(gitUrl, gitUrlSeparator)
	cleanUrl := splitUrl[0]
	hasReference := len(splitUrl) > 1

	if !hasReference {
		return cleanUrl, "", "", ""
	}

	if (splitUrl[1] == "") || (len(splitUrl) > 2) {
		ThrowErrorf("cannot properly parse url %s. aborting", gitUrl)
	}

	baseReference := splitUrl[1]
	commitHash, branch, tag := prepareGitReference(baseReference)

	return cleanUrl, commitHash, branch, tag
}

// Returns commitHash, branch, tag
func prepareGitReference(baseReference string) (
	string,
	plumbing.ReferenceName,
	plumbing.ReferenceName,
) {
	var commitHash string
	var branch plumbing.ReferenceName
	var tag plumbing.ReferenceName

	tagRegexp, _ := regexp.Compile(tagRegexp)

	if tagRegexp.MatchString(baseReference) {
		tag = plumbing.NewTagReferenceName(baseReference)
	} else if plumbing.IsHash(baseReference) {
		commitHash = baseReference
	} else {
		branch = plumbing.NewBranchReferenceName(baseReference)
	}

	baseLog := "parsed from git url "

	if commitHash != "" {
		commitLog := ColorTextf("commitHash=%s", hashColor, commitHash)
		log.Debug(baseLog + commitLog)
	}

	if branch != "" {
		branchLog := ColorTextf("branch=%s", branchColor, branch.Short())
		log.Debug(baseLog + branchLog)
	}

	if tag != "" {
		tagLog := ColorTextf("tag=%s", tagColor, tag.Short())
		log.Debug(baseLog + tagLog + " ")
	}

	return commitHash, branch, tag
}

func getGitAuth(repoUrl string) *ssh.PublicKeys {
	if strings.Contains(repoUrl, gitNormalUrl) {
		return nil
	}

	sshKeyPath := GetSshKeyPath()
	sshKeyPassword := GetSshKeyPassword()

	auth, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, sshKeyPassword)
	CheckError(err, "error while creating ssh git auth")

	return auth
}

func getGitColor(isHash bool, isTag bool) lipgloss.Color {
	if isHash {
		return hashColor
	}

	if isTag {
		return tagColor
	}

	return branchColor
}
