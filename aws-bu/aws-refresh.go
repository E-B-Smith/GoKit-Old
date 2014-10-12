package main

import (
	"io"
	"os"
	"fmt"
	"io/ioutil"
	)


func refreshData() AWSResultCode {

	//  refresh  -adl  [ <bundle-name> ]

	//	Write local meta-data -- 

	log(AWSLogDebug, "Writing local file meta-data.")

 	writer, error := os.Create("./TestData/TestBackup.ldata")
    if error != nil {
    	log(AWSLogError, "Can't open temporary file: %v.", error)
        return AWSResultError;
    	}

	workingDirectory, _ := os.Getwd()
	os.Chdir("./TestData")
	result := writeBundleStatusFile(writer, "TestBackup.sparsebundle")
	os.Chdir(workingDirectory)
	writer.Close()

	//	Write AWS meta-data -- 

	log(AWSLogDebug, "Writing AWS file meta-data.")

 	writer, error = os.Create("./TestData/TestBackup.adata")
    if error != nil {
    	log(AWSLogError, "Can't open temporary file: %v.", error);
        return AWSResultError;
    	}
//	result = writeAWSStatusFile(writer, "TestBackup.sparsebundle")
	result = writeAWSStatusFile(writer, "Brennos.sparsebundle")
	writer.Close()

	return result
	}


func writeBundleStatusFile(writer io.Writer, directory string) AWSResultCode {

	filearray, error := ioutil.ReadDir(directory)
	if (error != nil) {
		log(AWSLogError, "Error reading %s: %v.", directory, error)
		return AWSResultError
		}

	for _, file := range filearray {
		path := directory + "/" + file.Name();
		if file.IsDir() {
			result := writeBundleStatusFile(writer, path)
			if result != AWSResultSuccess { return result }
		} else {
			fmt.Fprintf(writer, "%s\t%v\t%v\n", path, file.ModTime(), file.Size())
			}
		}

	return AWSResultSuccess
	}


func writeAWSStatusFile(writer io.Writer, path string) AWSResultCode {
	listAWSObjectsWithPrefixAndMarker(writer, path, "")
	return AWSResultSuccess;
	}

	