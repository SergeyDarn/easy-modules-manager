package utils

import (
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
	repoLog := ColorText("repo="+repoName, repoColor)
	urlLog := ColorText("url="+repoUrl, urlColor)
	log.Debugf("Cloning %s %s", repoLog, urlLog)

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
	CheckError(err, "Error while clonning "+repoLog)

	if commitHash != "" {
		GitCheckoutToCommit(repo, repoName, commitHash)
	}

	headName := GetHeadShortName(repo, commitHash != "", tag != "")
	headColor := getGitColor(commitHash != "", tag != "")

	headLog := ColorText("head="+headName, headColor)
	successLog := ColorSuccess("Cloning successful")
	log.Debugf("%s %s %s", successLog, repoLog, headLog)
}

func GitDirStatus(dirPath string) git.Status {
	repo, err := git.PlainOpen(dirPath)
	CheckError(err, "Error while trying to open module directory "+dirPath+" in git for gitFolderStatus")

	workTree, err := repo.Worktree()
	CheckError(err, "Error while getting workTree for gitFolderStatus")

	status, err := workTree.Status()
	CheckError(err, "Error while trying to get git status")

	return status
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
	CheckError(err, fmt.Sprintf("Error while getting repo %s worktree before checkout", repoName))

	err = workTree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})

	repoLog := "repo=" + repoName
	commitLog := "commitHash=" + commitHash
	CheckError(err, fmt.Sprintf("Error while trying to checkout %s %s", repoLog, commitLog))

	repoColorLog := ColorText(repoLog, repoColor)
	commitColorLog := ColorText(commitLog, hashColor)
	log.Debugf("Sucessful checkout for %s to %s", repoColorLog, commitColorLog)
}

func IsGitUrl(url string) bool {
	return strings.Contains(url, "git")
}

func GetHeadShortName(repo *git.Repository, isCommit bool, isTag bool) string {
	head, err := repo.Head()
	CheckError(err, "Error while getting repo head (GetHeadShortName)")

	if isCommit {
		return head.Hash().String()
	}

	if isTag {
		headTag := GetHeadTag(repo)
		if headTag == nil {
			ThrowError(
				fmt.Sprintf("Couldn't find tag for head %s (GetHeadShortName)", head.String()),
			)
		}

		return headTag.Name().Short()
	}

	return head.Name().Short()
}

func GetHeadTag(repo *git.Repository) *plumbing.Reference {
	head, err := repo.Head()
	CheckError(err, "Error while getting repo head (GetHeadTag)")

	tags, err := repo.Tags()
	CheckError(err, "Error while getting repo tags (GetHeadTag)")

	var headTag *plumbing.Reference

	tags.ForEach(func(tag *plumbing.Reference) error {
		tagHash, err := repo.ResolveRevision(plumbing.Revision(tag.Name()))
		CheckError(err, "Error while resolving revision (GetHeadTag)")

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
		ThrowError("Cannot properly parse url " + gitUrl + " Aborting.")
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

	baseLog := "Parsed from git url "

	if commitHash != "" {
		commitLog := ColorText("commitHash="+commitHash, hashColor)
		log.Debug(baseLog + commitLog)
	}

	if branch != "" {
		branchLog := ColorText("branch="+branch.Short(), branchColor)
		log.Debug(baseLog + branchLog)
	}

	if tag != "" {
		tagLog := ColorText("tag="+tag.Short(), tagColor) + " "
		log.Debug(baseLog + tagLog)
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
	CheckError(err, "Error while creating git clone auth")

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
