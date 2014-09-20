package main

import (
	"io"
	"os"
	"fmt"
	"io/ioutil"
	)


func refreshData() AWSResultCode {

	//  refresh  -adl  [ <bundle-name> ]

 	writer, error := os.Create("./TestData/TestBackup.data")
    if error != nil {
    	log(AWSLogError, "Can't open temporary file: %v.", error);
        return AWSResultError;
    	}
	result := writeBundleStatusFile(writer, "./TestData/TestBackup.sparsebundle")
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


func writeAWSStatusFile(writer io.Writer, directory string) AWSResultCode {
	return AWSResultError
	}

	