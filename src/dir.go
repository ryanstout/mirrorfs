package main

import (
	"fmt"
	"log"
	"os"
  "io/ioutil"
  "path/filepath"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)


// Dir implements both Node and Handle for the root directory.
type Dir struct{
  Path string
}

func (dir *Dir) Attr() fuse.Attr {
  stats, err := os.Stat(dir.Path)
	if err != nil {
		log.Print(err)
		return fuse.Attr{}
	}

	return fuse.Attr{
    Size: uint64(stats.Size()),
    Mtime: stats.ModTime(),
    Mode: stats.Mode(),
    Uid: stats.Sys().(*syscall.Stat_t).Uid,
    Gid: stats.Sys().(*syscall.Stat_t).Gid,
  }

  //
  // return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0555}
}


/* Create a file */
func (dir *Dir) Create(req *fuse.CreateRequest, res *fuse.CreateResponse, intr fs.Intr) (fs.Node, fs.Handle, fuse.Error) {
  // fmt.Printf("Create: %s\n", req.Name)

  path := dir.Path + "/" + req.Name
  file := FileLookup(path)

  fsHandle, err := os.OpenFile(path, int(req.Flags), req.Mode)

  handle := &FileHandle{fsHandle, file}

  if err != nil {
    return nil, nil, err
  }

  return file, handle, nil
}

func (dir *Dir) Mkdir(req *fuse.MkdirRequest, intr fs.Intr) (fs.Node, fuse.Error) {
  path := dir.Path + "/" + req.Name
  fmt.Printf("MKDIR: %s, %s\n", req.Name, path)

  err := os.Mkdir(path, req.Mode)

  if err != nil {
    return nil, err
  }

  return &Dir{path}, nil
}


func (dir *Dir) Lookup(name string, intr fs.Intr) (fs fs.Node, error fuse.Error) {
	path := filepath.Join(dir.Path, name)
	stats, err := os.Stat(path)
	if err != nil {
    // log.Print(err)
		return nil, fuse.ENOENT
	}

  // fmt.Printf("Lookup: %s -- %s -- %s\n", path, dir.Path, name)

	switch {
	case stats.IsDir():
		fs = &Dir{path}
	case stats.Mode().IsRegular():
    // fmt.Printf("Regular: %s\n", path)
		fs = FileLookup(path)
	default:
    // fs = File{"missing"}
    fmt.Printf("Missing: %s\n", path)
    return nil, fuse.ENOENT
	}

	return
}

func (dir *Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	var out []fuse.Dirent
	files, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		log.Print(err)
		return nil, fuse.Errno(err.(syscall.Errno))
	}
	for _, node := range files {
		de := fuse.Dirent{Name: node.Name()}
		if node.IsDir() {
			de.Type = fuse.DT_Dir
		}
		if node.Mode().IsRegular() {
			de.Type = fuse.DT_File
		}
		out = append(out, de)
	}

	return out, nil
}

/* Remove a file or directory from the directory */
func (dir *Dir) Remove(req *fuse.RemoveRequest, intr fs.Intr) fuse.Error {
  path := dir.Path + "/" + req.Name
  fmt.Printf("Remove: %s\n", path)

  return os.Remove(path)
}


func Readlink(req *fuse.ReadlinkRequest, intr fs.Intr) (string, fuse.Error) {
  fmt.Printf("READLINK not supported yet\n")

	return "", fuse.EIO
}

func (dir *Dir) Link(req *fuse.LinkRequest, old fs.Node, intr fs.Intr) (fs.Node, fuse.Error) {
  fmt.Printf("LINK not supported yet\n")

	return nil, fuse.EIO
}

func (dir *Dir) Symlink(req *fuse.SymlinkRequest, intr fs.Intr) (fs.Node, fuse.Error) {
  fmt.Printf("Symlink not supported yet\n")

	return nil, fuse.EIO
}