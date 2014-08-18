rm ./drive
cd src ; go build -o ../drive ; cd ..
echo "compiled, running..."
GOMAXPROCS=4 ./drive /Users/ryanstout/Sites/infinitydrive/go/drivemount

echo "unmounting..."
umount /Users/ryanstout/Sites/infinitydrive/go/drivemount