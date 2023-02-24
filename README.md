Init

To do:

Hash password on front end
Ensure login call is functioning
Implement auth middleware

Fix error checking, and handle
Tests


So.  We are going to cache the product database images

We are going to hydrate the page with the cached images, and the query to the product database

If we cache the product database HTML, then we don't need the product database
assuming we have access to the assets

I want the assets to be bundled with the container, and put in the container repo
Then we can fetch the proper container for the user from the repo
This is slow...
We update the container repo when the container gets rebuilt, we can keep a few historical containers
Each container could be prod/dev/release/feature/hotfix instead of historical as well.

So we host the repo with Harbor, and we vuln scan with Trivy




