package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunCheckout(branch string, shouldCreate bool, repo *core.Repository) error {
	lastBranch := repo.GetCurrentBranch()

	if lastBranch == branch {
		return fmt.Errorf("you are already on branch %s", branch)
	}

	index, err := repo.LoadIndex()

	if err != nil {
		return err
	}

	// if len(index) > 0 {
	// 	return fmt.Errorf("staging area is not empty. commit changes before checkout")
	// }

	//Treba da cekira da li ima staged promjena, ako ima onda ne moze da se commita

	lastCommit := repo.GetBranchCommit(lastBranch)

	if lastCommit == "" {
		return fmt.Errorf("no commits yet, cannot create a branch")
	}

	commitObj, err := repo.LoadObject(lastCommit)
	if err != nil {
		return err
	}

	lastBranchFilesMap := map[string]core.IndexEntry{}

	if err := core.GetFilesFromTreeHash(commitObj.(*core.Commit).TreeHash, repo, "", lastBranchFilesMap); err != nil {
		return err
	}

	branchFile := filepath.Join(repo.RefsDir, branch)
	if _, err := os.Stat(branchFile); errors.Is(err, os.ErrNotExist) { //branch doesnt exist
		if shouldCreate {
			if lastCommit == "" {
				return fmt.Errorf("no commits yet, cannot create a branch")
			} else {
				if err := repo.SetBranchCommit(branch, lastCommit); err != nil { //if it's brand new, it should point to its parent last commit
					return err
				}
				fmt.Printf("created new branch %s", branch)

				//myb return here to skip deleting and recreating the same files
				//just add log from the end of this func
			}
		} else {
			return fmt.Errorf("branch %s not found\nuse checkout -b %s to create and switch to a new branch", branch, branch)
		}
	}

	if err := repo.SetCurrentBranch(branch); err != nil {
		return err
	}

	newBranchCommitHash := repo.GetBranchCommit(branch)
	newBranchCommit, err := repo.LoadObject(newBranchCommitHash)
	if err != nil {
		return err
	}

	newBranchFilesMap := map[string]core.IndexEntry{}

	if err := core.GetFilesFromTreeHash(newBranchCommit.(*core.Commit).TreeHash, repo, "", newBranchFilesMap); err != nil {
		return err
	}

	filesToRemove := []string{}
	for path, entry := range index {
		if val, ok := newBranchFilesMap[path]; ok {
			if entry.Hash != val.Hash {
				filesToRemove = append(filesToRemove, path)
			}
		} else {
			filesToRemove = append(filesToRemove, path)
		}
	}

	if err := RemoveOldFiles(filesToRemove, repo); err != nil {
		return err
	}

	if err := RestoreWorkingDirectoryFiles(newBranchCommit.(*core.Commit).TreeHash, "", repo); err != nil {
		return err
	}

	//save index with path:hash, mtime; where
	// 		path is core.GetFilesFromTreeHash(newBranchCommit.(*core.Commit).TreeHash, repo, "", &newBranchIndex)
	//		getfilesfromtreehash should be modified to access whole index
	//		mtime is modification time for each file in new index, os system call

	// if err := repo.SaveIndex(newBranchIndex); err != nil {
	// 	return err
	// }

	fmt.Printf("switched to branch %s", branch)
	return nil
}

func RemoveOldFiles(filesToRemove []string, repo *core.Repository) error {
	for _, file := range filesToRemove {
		fullPath := filepath.Join(repo.WorkTree, file)
		if err := os.Remove(fullPath); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
func RestoreWorkingDirectoryFiles(treeHash string, parentPath string, repo *core.Repository) error {
	obj, err := repo.LoadObject(treeHash)
	if err != nil {
		return err
	}
	tree := obj.(*core.Tree)

	for _, entry := range tree.Entries {
		fullPath := filepath.Join(repo.WorkTree, parentPath, entry.Name)
		switch entry.Mode {
		case "100644":
			obj, err := repo.LoadObject(entry.Hash)
			if err != nil {
				fmt.Println(err)
			}
			blob := obj.(*core.Blob)
			if err := os.WriteFile(fullPath, blob.Content, 0644); err != nil {
				fmt.Println(err)
			}

		case "040000":
			if err := os.Mkdir(fullPath, 0755); err != nil {
				//fmt.Println(err)
			}
			if err := RestoreWorkingDirectoryFiles(entry.Hash, filepath.Join(parentPath, entry.Name), repo); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
