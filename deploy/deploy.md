
Deploy
======

    deploy -nfrvi [host:]/files/deplydata  relcy-r1:hduser:hduser:644:~/data
    deploy relcy-r1:hduser:hduser:644:~/data2/file1

    -n    Dry-run.
    -f    Force.
    -r    Reverse deploy direction.
    -v    Verbose.
    -i    Input manifest file.


    deploy-group {

        user hduser;
        group hduser;
        permissions 644;

        deploy file1  file2;
        deploy {
            user hduser;
            group hduser;
            permissions 666;

            file3, file4;
            }
            file3, file4;

        }
    hosts {
        relcy-r1, relcy-r2, relcy-r3
        include "./other-hosts"
        };


Flow
----
* Parse & validate input.
* Make sure each host is accessible.
* (Create a deployment lock on the host).
* Make sure that each user and group on each host exists.
* rsync each file to a deployment directory.
// * rsync each file to a move directory.
* Fix permissions and owners of move directory.
* Move files into place.
* Report errors & success.
* Write to a log file.
* (Remove deployment lock).

sshd reallybangat

    