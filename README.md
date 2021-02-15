# Readme

## Version bumping

To version bump:

    #update VERSION variable in Makefile (near top)
    
    #create version bump commit and tag commit (on `main` branch)
    git commit -am "Version bump" --allow-empty
    
    #push to github
    git push origin main  --tags
    
    #install new version locally
    make install

## Initial setup

By default, golang is configured to use public github repositories. In
order to setup access to private github repos (supplycom github repos in
this case):

<https://golang.cafe/blog/how-to-fix-go-mod-unknown-revision.html>

    # set GO111MODULES
    go env -w GO111MODULE=on

    # add org private repo to GOPRIVATE
    go env -w GOPRIVATE=github.com/supplycom/k8s_client_go
    
    # configure git to use ssh
    git config --global url."ssh://git@github.com/supplycom".insteadOf "https://github.com/supplycom"

This needs to be done in order to import this module into other golang
projects. It requires ssh keys to be setup for github.

Seems there is also a way to enable private repo access using personal
acccess tokens (you're on your own for this):

<https://medium.com/swlh/go-modules-with-private-git-repository-3940b6835727>