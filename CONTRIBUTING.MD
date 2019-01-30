# How to contribute #

We'd more than happy to accept all forms of contribution. Those may come in form of submitting a pull request, opening an issue, improving documentation or tests, or even helping the community that may ask questions about the project on issues. However, in order to organize, it's good to follow some guidelines.

## Reporting issues ##

### Problems and/or bugs ###

Be descriptive about your issue. Specify which TeamCity, Golang and OS version you are running and provide helpful information through logs and/or debug information.

### Feature Requests ###

Feature requests are welcome, still, try to be as specific as possible about what the intended feature should be. Then, after proposing the feature, be open to discussion and patient regarding it's implementation.

## Developing ##

If you want to contribute in form of development, you may fork this repo, clone it locally (preferably inside your `GOPATH`) then create a branch to work on. [This link](https://opensource.guide/how-to-contribute/#opening-a-pull-request) has helpful information to get started.

Seen an issue that you want to tackle? Make sure to comment indicating that intention, so others know that someone is already willing to work on it.

Try to maitain the code conventions, or spot where they are inconsistent across the codebase, whenever possible. When submitting your work, if you would like to propose an enhancement to existing conventions, do it in a separate issue for discussion, other than the issued worked on.

You'll need Go (1.10+) installed locally. For instructions on how to do that for your platform, you may refer to [Golang Installation](https://golang.org/doc/install).

### Running Locally ###

TeamCity can be run with the official docker images from JetBrains. We leverage that by implementing full integration tests with every feature.
You are expected to fully run these tests and make sure your changes do not break existing functionality.

In order to run the tests locally, run following command from the repository root folder:

```bash
make test
```

This will start a Docker container to run all tests against. If the container is already up and running, it will simply run the tests.

Once you are done with testing, you can destroy the test environment:

```bash
make clean-docker
```

If using an editor such as **Visual Studio Code** with `Go` integration, make sure to configure your workspace settings, such as:

```json
{
    "go.testEnvVars": {
        "TEAMCITY_ADDR":"http://localhost:8112"
    },
    "go.testFlags": ["-v"],
    "go.buildFlags": ["-v"]
}
```

If you have followed the steps but still need help to setup a development, help us by opening an issue so we can improve the documentation and the process for setting up the development environment. :+1:
