source dev.env


go install
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi


echo "Doing some cleaning ..."
go clean
echo "Done."


echo "Running go format ..."
gofmt -w .
echo "Done."


echo "Running go build for test ..."
go build -race
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi
echo "Done."

echo "Running unit test ..."
go test -parallel 1 


echo "Running go build ..."
go build -race
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi
echo "Done."


if [ $? == 0 ]; then
    echo "Done."
	echo "## Starting service ##"
    ./elb-healthcheck-bulk-updater
fi