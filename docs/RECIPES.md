# Recipes

A recipe is a way of describing how to configure and run a built tool so that it can be executed in a uniform fashion.

## Default recipes

The default recipe list contains:
* [Ko](https://github.com/ko-build/ko) (`com.github.google.ko`) for Go
* [Jib](https://github.com/GoogleContainerTools/jib) (`com.google.cloud.tools.jib-maven-plugin`) for Java
* [BuildKit](https://github.com/moby/buildkit) (`com.github.moby.buildkit`) for Dockerfile
* [Nib](https://github.com/djcass44/nib) (`com.github.djcass44.nib`) for Static web applications
* [`all-your-base`](https://github.com/djcass44/all-your-base) (`com.github.djcass44.all-your-base`) for base images

## Custom recipes

Custom recipes can be provided using two flags:

1. `--recipe-template` - overrides the `recipes.tpl.yaml` file that is used. When this flag is provided, the default recipes will not be available.
2. `--extra-recipe-template` - appends and merges a custom `recipes.tpl.yaml` file with the default. When recipes with the same name as a default recipe are provided, they default recipe will be replaced.

## Writing a recipe

Writing a recipe is fairly simple.
A recipe can contain any build tool that results in the creation of an OCI artefact.

We will use the `docker build` command as an example.

```yaml
build:
  com.docker:
    # create a ~/.docker/config.json file so that
    # we can authenticate
    dockercfg: true
    # 'cd' into the directory that we are building.
    # Some tools such as docker expect the working directory and the
    # build directory to be the same
    cd: true
    command: docker
    args:
      - build
      - -f={{ .Dockerfile.File }}
      # add a '-t' for each tag
      {{- range $i, $e := .FQTags }}
      - -t={{ $e }}
      {{- end }}
      - {{ .Context }}
```

When executed, this will roughly equate to:

```shell
docker build -f=Dockerfile -t registry.example.org/foo/bar:v1.2.3 .
```

### Templating

The recipe file is templated using Go's [`text/template`](https://pkg.go.dev/html/template) package.
The available values for you to use can be found in the [`BuildContext`](../internal/api/v1/types.go) struct.
