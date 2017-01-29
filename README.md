# Build tools for standard makefile

These tools add standard components (sub-makefiles and build scripts) as well as example starters for the Makefile and circle.yml.

## Add build-tools to a Makefile

```
git remote add -f build-tools git@github.com:drud/build-tools.git
git merge -s ours --allow-unrelated-histories build-tools/master
git read-tree --prefix=build-tools -u build-tools/master
git commit -m "Added build-tools for standard makefile as subtree"
```

## Update build-tools directory from this repository using subtree merge

```
# If there is not a build-tools remote, add it
git remote add -f build-tools git@github.com:drud/build-tools.git
# Pull current build-tools
git pull -s subtree build-tools master
```

## Set up a Makefile to begin with

* Copy the Makefile.example to "Makefile" in the root of your project
* Edit the sub-Makefiles included
* Update the variables at the top of the Makefile
