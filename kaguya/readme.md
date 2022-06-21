# Kaguya Archiver

## Archival Flow

We query catalog.json and archive.json to start with. Old cached threads that are found in the archive are discarded, because no one really cares about archive post deletions, and old cached threads that are not found in either the catalog or the archive are assumed to have been deleted (note: do not use Kaguya for archiveless boards like /b/).

Old cached threads found in the catalog have LastModified values compared, and if the thread has been modified it's added to a queue. All threads in the queue are queried for, diffed against cached versions and the database is updated.

Meanwhile, new images are put in a channel belonging to the ImagesService. The ImagesService gradually downloads all images and puts them on S3.

## Project Structure

- main.go is the entrypoint file that sets up the appropriate archival loops.
- db wraps database interactions.
- api wraps api interactions.
- manager orchestrates the archival process and the caching.
- images downloads images and stores them in S3.
- utils is basically just random garbage that needs to be globally accessible.
- config just retrieves configuration.

## Notes:

- 4chan hands us file MD5s. MD5s have known collisions even for jpg files. We cannot skimp out on downloading images because we already have some file with the MD5, even though we want to. The archive search does use MD5, however, because it does work even though it isn't reliable (also very practical).
- Under regular conditions the list of first few posts that comes in /catalog.json is straight up useless because some moderator might have publically banned a post or someone might have deleted an image in it, so we'll need to query for the whole thread regardless.
- Everything that was deemed useless is missing from the api models and may be added in as necessary.
- We can reasonably assume that cached threads that have disappeared from both the catalog and the archive were deleted so long as we're taking less than three days between loops.
- I believe the Asagi/Torako strategy for identifying thread deletion is looking at the last page a thread was in before it vanished from the catalog. Which makes perfect sense in the archiveless case but is a bit of a pain otherwise.
- We don't do any sort of modification of post content. The plan is outputting everything raw and using CSP to prevent injection problems.
- Kaguya will crash if there are any enum problems. Say, 4chan introducing webp support will cause Kaguya to crash when it tries to insert a post with a webp file extension because .webp isn't in the enum. This problem won't necessitate any modification of the codebase, .webp just needs to be added to the postgres enum. Similar issues might happen if there are new board flags or capcodes.
- We could generate thumbnails in go code (go has native libraries for this. We can also use lilliput or whatever) and reduce request numbers but doing so is liable to break gimmicks like the image-looks-different-when-you-expand thing. Also generating thumbnails for pdfs probably fucking sucks.

