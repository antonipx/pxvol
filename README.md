pxvol
=====

Find who has mounted Portworx volume.

* dockerps - use docker api
* procps - use /proc filesystem

Usage
-----

* dockerps <volumeName>

* procps <volumeID>

Note procps uses numeric Portworx volume ID while dockerps uses name.
