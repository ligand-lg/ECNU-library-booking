rm -rf ./dist
mkdir ./dist

mkdir ./dist/ecnu_booking_windows
cp conf.json mh.json zb.json ./dist/ecnu_booking_windows/
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./dist/ecnu_booking_windows/booking.exe


mkdir ./dist/ecnu_booking_linux
cp conf.json mh.json zb.json ./dist/ecnu_booking_linux/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o ./dist/ecnu_booking_linux/booking


mkdir ./dist/ecnu_booking_mac
cp conf.json mh.json zb.json ./dist/ecnu_booking_mac/
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./dist/ecnu_booking_mac/booking
