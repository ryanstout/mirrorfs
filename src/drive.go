// Hellofs implements a simple "hello world" file system.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)


var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

  MIRROR_FOLDER := "/Users/ryanstout/Sites/infinitydrive/go/drive"

	flag.Usage = Usage
	flag.Parse()



	if flag.NArg() != 1 {
		Usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	c, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, FS{MIRROR_FOLDER})
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

// FS implements the hello world file system.
type FS struct{
  Path string
}

func (fs FS) Statfs(req *fuse.StatfsRequest, res *fuse.StatfsResponse, intr fs.Intr) fuse.Error {
	// Make some stuff up, just to see if it makes "lsof" happy.
	res.Blocks = 1 << 35
	res.Bfree = 1 << 34
	res.Bavail = 1 << 34
	res.Files = 1 << 29
	res.Ffree = 1 << 28
	res.Namelen = 2048
  // res.Bsize = 1024
  // res.Bsize = 32768
  // res.Bsize = 49152
  res.Bsize = 4096 * 1024
	return nil
}

func (dir FS) Root() (fs.Node, fuse.Error) {
	return &Dir{dir.Path}, nil
}
