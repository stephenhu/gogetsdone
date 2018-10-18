# icons

provide an icon storage service for getsdone

## requirements

* should limit the size of photos, ~100K # need to investigate reasonableness
* thumbnails should be generated
* typical flow should require upload from photos on iphone to target
* client side can handle thumbnailing, save on transfer and storage, hopefully just use iphone thumbnail since it already exists
* cache on client side, need to limit cache size


## persistent storage (p0)

lots of small files, need to find the average size, probably around 3K - 30K.
should the files be stored concatenated.

use glbs for storing on the filesystem, expose these links directly through
nginx.

## thumbnailing (p0)

should retrieve the thumbnail from storage

## api (p0)

expose an api to upload a single small file at a time

## caching (p1)

### server cache

### client cache
