#!/bin/bash

# license header added for current workdir

for i in $(find  ./*  -not -path "./vendor/*"  -name "*.go"); 
do
  file $i  |grep -q "symbolic link"
  isLinkFile=$?
  if [ $isLinkFile -eq 0 ];
  then
     echo "$i is link file, skip..."
  fi
  grep -q "Copyright 2021 The AtomCI Group Authors." $i
  existLicenseHeader=$?
  if [ $existLicenseHeader -eq 0 ];
  then
    echo "$i already have license header, skip...";
    echo 
    continue
  fi
  sed -i '1i \
/*\
Copyright 2021 The AtomCI Group Authors.\
\
Licensed under the Apache License, Version 2.0 (the "License");\
you may not use this file except in compliance with the License.\
You may obtain a copy of the License at\
\
	http://www.apache.org/licenses/LICENSE-2.0\
\
Unless required by applicable law or agreed to in writing, software\
distributed under the License is distributed on an "AS IS" BASIS,\
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\
See the License for the specific language governing permissions and\
limitations under the License.\
*/\
' $i;
done
