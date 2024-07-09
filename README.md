# lurl

little url app

There are lots of public URL shortners. This is my own.

# Functionality

This app has no mechanism (from the web) to add URLs. This is intentional, as
it prevents the need to write a login or authentication system. The URLs are
stored in a plain text file which is space-separated records one per line and
read on each request.  Each line has two keys: the short tag and the
destination URL.

# Configuration

This is run as a resident app.

There are two environment variables:

* `PORT` the port on which to listen. Defaults to 8080 if not set.
* `LURLS` the file from which to read the URLs. Defaults to `lurls.txt` in the same dir as the app if not set.

# License
Written by Matt Peperell (matt@peperell.com), licensed under MIT. The text for
this license should be distributed with this software but it can also be found
at https://en.wikipedia.org/wiki/MIT_License
