package main

import (
	"os"
  "fmt"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)




type FileHandle struct{
  File *os.File
  FileRef *File
}

func (fh *FileHandle) Write(req *fuse.WriteRequest, resp *fuse.WriteResponse, intr fs.Intr) fuse.Error {
  fmt.Printf("Write %s - %d at %d\n", fh.FileRef.Path, len(req.Data), req.Offset)

  size, err := fh.File.WriteAt(req.Data, req.Offset)
  resp.Size = size

  return err
}

func (fh *FileHandle) Read(req *fuse.ReadRequest, resp *fuse.ReadResponse, intr fs.Intr) fuse.Error {
  fmt.Printf("Read %s - %d at %d\n", fh.FileRef.Path, req.Size, req.Offset)

  buf := make([]byte, req.Size)
  // buf := resp.Data

  // fmt.Printf("ReadAt: %d -- %d -- %d\n", req.Offset, req.Size, len(resp.Data))
  size, err := fh.File.ReadAt(buf, req.Offset)

  if err != nil {
    fmt.Printf("ERR FROM READ: %s\n", err)

    // return fuse.EIO
    // return err
  }
  fmt.Printf("DID READ %d\n", size, buf)//[:size])

  resp.Data = buf//[:size]

  return err
}

func (fh *FileHandle) Release(*fuse.ReleaseRequest, fs.Intr) fuse.Error {
  // path := (*fh.FileRef).Path

  // OpenFileReferenceAtPath(path)

  fmt.Printf("Close %s\n", fh.FileRef.Path)

	return fh.File.Close()
}

func (h *FileHandle) Flush(req *fuse.FlushRequest, intr fs.Intr) fuse.Error {
  fmt.Printf("FLUSH\n")
  return nil

}

// func (h *FileHandle) Flush(req *fuse.FlushRequest, intr fs.Intr) fuse.Error {
//   err := h.file.dir.fs.db.Update(func(tx *bolt.Tx) error {
//     b := h.file.dir.bucket(tx)
//     if b == nil {
//       return fuse.ESTALE
//     }
//     return b.Put(h.file.name, h.data)
//   })
//   if err != nil {
//     return err
//   }
//   return nil
// }


// func (f File) Read(req *fuse.ReadRequest, resp *fuse.ReadResponse, intr fs.Intr) fuse.Error {
//
// }