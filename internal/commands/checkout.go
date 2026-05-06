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

				if err := repo.SetCurrentBranch(branch); err != nil {
					return err
				}
				fmt.Printf("switched to branch %s", branch)
				return nil
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

	// if err := RestoreWorkingDirectoryFiles(newBranchCommit.(*core.Commit).TreeHash, "", repo); err != nil {
	if err := RestoreWorkingDirectoryFiles(newBranchFilesMap, lastBranchFilesMap, "", repo); err != nil {
		return err
	}

	if err := repo.SaveIndex(newBranchFilesMap); err != nil {
		return err
	}

	fmt.Printf("switched to branch %s", branch)
	return nil
}

func RemoveOldFiles(filesToRemove []string, repo *core.Repository) error {
	for _, file := range filesToRemove {
		fullPath := filepath.Join(repo.WorkTree, file)
		if err := os.Remove(fullPath); err != nil {
			fmt.Println(err)
		}
		if err := os.Remove(filepath.Dir(fullPath)); err != nil {
			continue
		}
	}

	return nil
}

func RestoreWorkingDirectoryFiles(newBranchIndex, lastIndex map[string]core.IndexEntry, parentPath string, repo *core.Repository) error {
	for path, entry := range newBranchIndex {
		old, ok := lastIndex[path]
		if !ok || entry.Hash != old.Hash {
			fullPath := filepath.Join(repo.WorkTree, path)
			obj, err := repo.LoadObject(entry.Hash)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				fmt.Println(err)
			}

			blob := obj.(*core.Blob)

			_ = os.MkdirAll(filepath.Dir(fullPath), 0755)

			if err := os.WriteFile(fullPath, blob.Content, 0644); err != nil {
				fmt.Println(err)
			}

			info, _ := os.Stat(path)
			entry.MTime = info.ModTime().Unix()
			newBranchIndex[path] = entry
		} else {
			entry.MTime = old.MTime
			newBranchIndex[path] = entry
		}
	}

	return nil
}

// func RestoreWorkingDirectoryFilesv2(treeHash string, parentPath string, repo *core.Repository) error {
// 	obj, err := repo.LoadObject(treeHash)
// 	if err != nil {
// 		return err
// 	}
// 	tree := obj.(*core.Tree)

// 	for _, entry := range tree.Entries {
// 		fullPath := filepath.Join(repo.WorkTree, parentPath, entry.Name)
// 		switch entry.Mode {
// 		case "100644":
// 			obj, err := repo.LoadObject(entry.Hash)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			blob := obj.(*core.Blob)
// 			if err := os.WriteFile(fullPath, blob.Content, 0644); err != nil {
// 				fmt.Println(err)
// 			}

// 		case "040000":
// 			if err := os.Mkdir(fullPath, 0755); err != nil {
// 				//fmt.Println(err)
// 			}
// 			if err := RestoreWorkingDirectoryFilesv2(entry.Hash, filepath.Join(parentPath, entry.Name), repo); err != nil {
// 				fmt.Println(err)
// 			}
// 		}
// 	}
// 	return nil
// }
