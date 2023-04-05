package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"

	goGit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var logger zerolog.Logger

func main() {
	fmt.Println("Started")
	path := "C:\\projects\\CCM2"
	logger = zerolog.New(os.Stdout).With().Logger()

	err := os.RemoveAll(path)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Could not create folder: %v", path)
	}

	os.Mkdir(path, os.ModePerm)
	logger.Info().Msgf("Created git folder at: %v", path)

	// Init git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = path

	out, err := cmd.Output()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error occurred while trying to initialize git in folder: %v", path)
	}

	logger.Debug().Msgf("Initialized git repository at: %v", path)
	logger.Debug().Msgf("Output was: %v", string(out))

	// Rename 'main' branch to 'trunk'
	cmd = exec.Command("git", "branch -m master trunk")
	cmd.Dir = path

	out, _ = cmd.Output()
	// Somehow git always write to StfErr (???)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msgf("Error occurred while trying to initialize git in folder: %v", path)
	// }

	logger.Debug().Msg("Changed named of 'main' branch to 'trunk'")
	logger.Debug().Msgf("Output was: %v", string(out))

	// goGit.Init()

	repository, err := goGit.PlainOpenWithOptions(path, &goGit.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		logger.Fatal().Err(err).Msgf("Could not get repository from folder: %v", path)
	}

	workTree, err := repository.Worktree()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Could not get worktree")
	}

	logger.Debug().Msgf("Repository initialized at folder: %v", filepath.Join(path, goGit.GitDirName))

	//------------------------------------------------------------------------------------------------------------------------------------

	file, _ := os.OpenFile(filepath.Join(path, "master.txt"), os.O_CREATE|os.O_TRUNC, os.ModePerm)
	io.WriteString(file, "On branch master")
	file.Close()

	status, err := workTree.Status()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error occurred while trying to receive a status of current branch")
	}

	for key, _ := range status {
		hash, err := workTree.Add(key)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not add file: %v", filepath.Join(path, key))
		}

		logger.Debug().Msgf("Added file: %v with hash: %v", filepath.Join(path, key), hash)
	}

	status, err = workTree.Status()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error occurred while trying to receive a status of current branch")
	}
	logger.Debug().Msg("2. Status")

	t, _ := time.Parse("Jan 2, 2006", "Jan 1, 1978")

	commit, err := workTree.Commit("This is commit on master", &goGit.CommitOptions{
		Author: &object.Signature{
			Name:  "Aleks",
			Email: "af@ncr.com",
			When:  t,
		},
	})

	if err != nil {
		logger.Fatal().Err(err).Msgf("Could not commit")
	}

	logger.Debug().Msgf("Commit is: %v", commit)

	h, err := repository.Head()
	if err != nil {
		logger.Fatal().Err(err).Msg("Can not get head !")
	}

	logger.Debug().Msgf("Head name is: %v and hash is: %v", h.Name(), h.Hash())

	//------------------------------------------------------------------------------------------------------------------------------------

	workTree.Checkout(&goGit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("trunk"),
		Create: true,
	})

	file, _ = os.OpenFile(filepath.Join(path, "trunk.txt"), os.O_CREATE|os.O_TRUNC, os.ModePerm)
	io.WriteString(file, "On branch trunk")
	file.Close()

	status, err = workTree.Status()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error occurred while trying to receive a status of current branch")
	}

	for key, _ := range status {
		hash, err := workTree.Add(key)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not add file: %v", filepath.Join(path, key))
		}

		logger.Debug().Msgf("Added file: %v with hash: %v", filepath.Join(path, key), hash)
	}

	t1, _ := time.Parse("Jan 2, 2006", "Jan 1, 1976")

	commit, err = workTree.Commit("This is commit on Trunk", &goGit.CommitOptions{
		Author: &object.Signature{
			Name:  "Aleks",
			Email: "af@ncr.com",
			When:  t1,
		},
	})

	if err != nil {
		logger.Fatal().Err(err).Msgf("Could not commit")
	}

	logger.Debug().Msgf("Commit is: %v", commit)

	//------------------------------------------------------------------------------------------------------------------------------------

	workTree.Checkout(&goGit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("v1.7.0"),
		Create: true,
	})

	file, _ = os.OpenFile(filepath.Join(path, "file-1.7.0.txt"), os.O_CREATE|os.O_TRUNC, os.ModePerm)
	io.WriteString(file, "On branch 1.7.0")
	file.Close()

	status, err = workTree.Status()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error occurred while trying to receive a status of current branch")
	}

	for key, _ := range status {
		hash, err := workTree.Add(key)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not add file: %v", filepath.Join(path, key))
		}

		logger.Debug().Msgf("Added file: %v with hash: %v", filepath.Join(path, key), hash)
	}

	t1, _ = time.Parse("Jan 2, 2006", "Jan 1, 1976")

	commit, err = workTree.Commit("This is commit on 1.7.0", &goGit.CommitOptions{
		Author: &object.Signature{
			Name:  "Aleks",
			Email: "af@ncr.com",
			When:  t1,
		},
	})

	if err != nil {
		logger.Fatal().Err(err).Msgf("Could not commit")
	}

	//------------------------------------------------------------------------------------------------------------------------------------

	workTree.Checkout(&goGit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("trunk"),
		Create: true,
	})

	file, _ = os.OpenFile(filepath.Join(path, "trunk2.txt"), os.O_CREATE|os.O_TRUNC, os.ModePerm)
	io.WriteString(file, "On branch trunk")
	file.Close()

	status, err = workTree.Status()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Error occurred while trying to receive a status of current branch")
	}

	for key, _ := range status {
		hash, err := workTree.Add(key)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not add file: %v", filepath.Join(path, key))
		}

		logger.Debug().Msgf("Added file: %v with hash: %v", filepath.Join(path, key), hash)
	}

	t1, _ = time.Parse("Jan 2, 2006", "Jan 1, 1926")

	commit, err = workTree.Commit("This is commit on trunk 2", &goGit.CommitOptions{
		Author: &object.Signature{
			Name:  "Aleks",
			Email: "af@ncr.com",
			When:  t1,
		},
	})

	if err != nil {
		logger.Fatal().Err(err).Msgf("Could not commit")
	}

}
