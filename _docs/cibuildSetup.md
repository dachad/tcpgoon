# CI BUILD SETUP

## Docker build assessment

### Option 1 - Travis driven

How?

* travis could do (cibuild script, using others):
    * running unit (go tests)
    * compiling (go build)
    * end to end / cmdline tests (to be developed)
    * docker build
    * end to end / cmdline tests using docker image (to be developed)
    * docker push
    * updating readme -help (to be developed)

Advantages?
* Just a single build/deploy logic in travis
* Easy to reproduce locally

Problems?
* We will need to maintain the description in dockerhub manually

What to improve?
* Binaries (and OS packages in the future) should be deployed to an artifact/ospackages repo
* Dockerbuild could pick it from there and/or still use local-generated binaries for testing
    * Maybe depending on the travis tag (merge to master means the usage of the artifact store)

### Option 2 - Docker hub automatic builds

How?

* **Option 2A - travis triggering the build**
    * The dockerfile and its README / dockerhub description could be in a different github repo
    * Travis should monitor dockerhub build, so once completed, it fetches the latest version and performs
end to end / cmdline tests using docker image (to be developed)
    * If they fail, we should remove the latest docker image from dockerhub, if that is possible from travis...

* **Option 2B - dockerhub watching for github pushes**
    * Requires the dockerfile and READMEs to be stored in the main repo, unless the dedicated repo for docker
has some kind of dependency with the main one and we need to manually update a "pined version" to trigger the build
    * we cannot run end to end tests in travis... only alternative could be the dockerfile introducing this as 
a build step. Not nice (execution environment will become dirtier)

Advantages?
* A README file in github will work as the description in dockerhub, being versioned and always updated

Problems?
* We need to commit binaries to github  or start using an artifact ASAP
    * If we don't introduce an artifact repo now (too soon?), we need to allow travis to push (-f) the binary with github keys 
* If we don't have a dedicated github repo for docker, public project README will be also the dockerhub description, so we need to 
present the info for docker users properly
* Automated build/testing will become more complex

### Decision

Start by option 1, so everything driven by travis (do dockerhub automated builds). 

Cleaner; easier to troubleshoot, reproduce and follow. At the cost of maintaining a public docker description
in parallel, in the web

