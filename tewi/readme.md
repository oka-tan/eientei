# Tewi

Tewi is a half-baked read-only API for serving Eientei's basic data as JSON, but no searches.

## API Documentation

Swagger doesn't fucking work but I'll write it out as soon as I find something that does.

Until then:

- Endpoints are /:board (returns an array of thread OPs), /:board/posts (array of posts), /:board/thread/:resto (array of posts in the thread), /:board/post/:no (returns a single post).
- Empty arrays aren't returned, instead the body is empty and you get a 204 code.
- Posts and OPs follow the model in model.go.
- The parameter order can be passed to /:board and /:board/posts to change the order from the default desc to ascending (by passing in "asc", case-sensitive).
- The parameter keyset can be passed to /:board and /:board/posts to paginate through generic keyset pagination.

