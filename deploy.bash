#!/bin/bash

APPNAME="power-monitor"
MAINPKG="github.com/jimmysawczuk/power-monitor/cmd/power-monitor"
VERSION=`git describe --tags --abbrev=0`

make release;

# Setup
mkdir -p deploy
echo "#!/bin/bash" > go-env.sh
go env >> go-env.sh;
source go-env.sh;

echo "Host OS:" $GOHOSTOS;
echo "Host Architecture:" $GOHOSTARCH;
echo ""

# Build
# for GOOS in windows darwin linux; do
for GOOS in linux; do
	# for GOARCH in amd64 386; do
	for GOARCH in amd64; do
		exe=""
		if [[ $GOOS == "windows" ]]; then
			exe=".exe"
		fi

		echo "$GOOS/$GOARCH" && echo "----------------------------";
		GOOS=$GOOS GOARCH=$GOARCH go build -v -o $APPNAME $MAINPKG

		# if this is the host OS/arch, the exe is put in the root of bin rather than a subdirectory
		if [[ $GOOS == $GOHOSTOS ]] && [[ $GOARCH == $GOHOSTARCH ]]; then
			mv $APPNAME deploy/$APPNAME-$VERSION-${GOOS}-${GOARCH}$exe
		else
			mv $APPNAME deploy/$APPNAME-$VERSION-${GOOS}-${GOARCH}$exe
		fi
		echo "";
	done
done

# Cleanup
rm go-env.sh
