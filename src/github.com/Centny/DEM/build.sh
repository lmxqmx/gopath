#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
##############################
######Install Dependence######
echo "Installing Dependence"
# go get github.com/go-sql-driver/mysql
##############################
#########Running Test#########
echo "Running Test"
pkgs="\
 github.com/Centny/DEM\
 github.com/Centny/DEM/example\
"
echo "mode: set" > a.out
for p in $pkgs;
do
 if [ "$1" = "-u" ];then
  go get -u $p
 fi
 go test -v --coverprofile=c.out $p
 cat c.out | grep -v "mode" >>a.out
done
gocov convert a.out > coverage.json

##############################
#####Create Coverage Report###
echo "Create Coverage Report"
cat coverage.json | gocov-xml -b $GOPATH/src > coverage.xml
cat coverage.json | gocov-html coverage.json > coverage.html
