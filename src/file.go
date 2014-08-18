package main

import (
	"log"
	"os"
	"syscall"
  "fmt"
  "bytes"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
  "bazil.org/fuse/syscallx"
)




// File implements both Node and Handle for the hello file.
type File struct{
  Path string
}

func FileLookup(path string) *File {
  return &File{path}
}

func (file *File) Attr() fuse.Attr {
  stats, err := os.Stat(file.Path)
	if err != nil {
		log.Print(err)
		return fuse.Attr{}
	}


  fmt.Printf("File Size: %d\n", uint64(stats.Size()))

	return fuse.Attr{
    Size: uint64(stats.Size()),
    Mtime: stats.ModTime(),
    Mode: stats.Mode(),
    Uid: stats.Sys().(*syscall.Stat_t).Uid,
    Gid: stats.Sys().(*syscall.Stat_t).Gid,
  }
}

func (file *File) uidAndGid() (int, int) {
  stats, err := os.Stat(file.Path)
	if err != nil {
		panic(err)
	}

  return int(stats.Sys().(*syscall.Stat_t).Uid), int(stats.Sys().(*syscall.Stat_t).Gid)
}



func (file *File) Open(req *fuse.OpenRequest, res *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
  fmt.Printf("Open with: %d - %d\n", req.Flags, int(req.Flags))

  // Allow kernel to use buffer cache
  res.Flags &^= fuse.OpenDirectIO

  fsHandle, err := os.OpenFile(file.Path, int(req.Flags), 0)

  handle := &FileHandle{fsHandle, file}

  if err != nil {
    return nil, err
  }

  return handle, nil
}

/* Called from the FileHandle when the file is closed */
func (file *File) CloseFromHandle() {

}

/* A special cast, since this number comes in actually as a int32, but bazil.org/fuse treats it
 * as an uint32 */
func castUint32ToInt(someInt uint32) int {
  if someInt > 1667855980 {
    return -1
  } else {
    return int(someInt)
  }
}

func (file *File) Setattr(req *fuse.SetattrRequest, resp *fuse.SetattrResponse, intr fs.Intr) fuse.Error {
  var err error

  // TODO: Check size for file truncate
  fmt.Printf("Setattr %s\n", file.Path)

  if req.Valid.Size() {
    fmt.Printf("CHANGE SIZE REQUESTED: %d\n", req.Size)

  }

  // don't change value
  if req.Valid.Mode() {
    // fmt.Printf("Change Mode on %s -- %d\n", file.Path, req.Mode)

    err = os.Chmod(file.Path, req.Mode)
    if err != nil {
      fmt.Printf("Error: %s\n", err)
      return err
    }
  }

  if req.Valid.Uid() {
    uid := castUint32ToInt(req.Uid)
    gid := castUint32ToInt(req.Gid)

    // fmt.Printf("Chown: %d vs %d --  %d vs %d\n", uid, req.Uid, gid, req.Gid)

    err = os.Chown(file.Path, uid, gid)
    if err != nil {
      fmt.Printf("Error: %s\n", err)
      return err
    }
  }

  if req.Valid.Chgtime() {
    err = os.Chtimes(file.Path, req.Atime, req.Mtime)
    if err != nil {
      fmt.Printf("Error: %s\n", err)
      return err
    }
  }
  return nil
}


func (file *File) Getxattr(req *fuse.GetxattrRequest, res *fuse.GetxattrResponse, intr fs.Intr) fuse.Error {
  // return fuse.EPERM
  // fmt.Printf("Get attr: %s\n", req.Name)


  buf := make([]byte, 8192)
  size, err := syscallx.Getxattr(file.Path, req.Name, buf)

  if err != nil {
    // fmt.Printf("Get xattr error: %s - %s: \n", file.Path, req.Name, err)
    // return err
    // On osx, we need to return NOATTR, but this isn't built into go or bazil.org/fuse, so we need to do this:
    return fuse.Errno(93)
  }

  res.Xattr = buf[:size]

  return nil
}

// TODO: This gets called twice, @Tv42 says this is just how the API works
func (file *File) Listxattr(req *fuse.ListxattrRequest, res *fuse.ListxattrResponse, intr fs.Intr) fuse.Error {
  // Get how large of a buffer is needed
  buf := make([]byte, 0)
  size, err := syscallx.Listxattr(file.Path, buf)
  if err != nil {
    fmt.Printf("Err listing attr: %s\n", err)
    return err
  }

  buf = make([]byte, size)
  size, err = syscallx.Listxattr(file.Path, buf)
  if err != nil {
    fmt.Printf("Err listing attr2: %s\n", err)
    return err
  }

  if size > 0 {
    attrNameBytes := bytes.Split(buf[:size-1], []byte{0})

    for _, name := range attrNameBytes {
      res.Append(string(name))
    }
  }
  return nil
}

func (file *File) Removexattr(req *fuse.RemovexattrRequest, intr fs.Intr) fuse.Error {
  err := syscallx.Removexattr(file.Path, req.Name)

  if err != nil {
    fmt.Printf("Err removing attr: %s\n", err)
    return err
  }

  return nil
}

func (file *File) Setxattr(req *fuse.SetxattrRequest, intr fs.Intr) fuse.Error {
  fmt.Printf("SetXattr: %s: %s\n", file.Path, req.Name)

  // TODO: Passing flags causes an exception
  err := syscallx.Setxattr(file.Path, req.Name, req.Xattr, 0)//int(req.Flags))
  if err != nil {
    fmt.Printf("SetXattr err: %s\n", err)
    return err
  }

  return nil

}
