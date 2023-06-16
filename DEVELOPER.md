# Development Setup

If you want to contribute to the Terraform provider for SAP BTP, be aware of the [contribution guidelines](CONTRIBUTING.md) available in this repository.

First, you need to setup your development environment. The following sections describe the options that you have.

## GitHub Codespaces

**Step 1:** Open the repository in GitHub Codespaces via the button:

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=618531988)

**Step 2:** There is no step 2 😎.

## Dev Container

> **Note** - In order to use dev containers you must have a container runtime up and running on the machine. For details, we refer to the official documentation about [Developing inside a Container](https://code.visualstudio.com/docs/devcontainers/containers)

First, you must clone the repository:

```bash
git clone https://github.com/SAP/terraform-provider-btp.git
```

Then open the cloned repository in [Visual Studio Code](https://code.visualstudio.com/). Within Visual Studio Code, press the "Open a remote Window" button in the lower left corner:

![screenshot of Visual Studio Code - Open a Remote Window](assets/VSCode_Show_Open_Remote_Window.png)

 Visual Studio Code will open the command palette. Choose the option "Reopen in Container":

![screenshot of Visual Studio Code - Open a Remote Window](assets/VSCode_Command_Palette_Reopen.png)

This will trigger the start of the dev container. You can choose to open a devcontainer with two configurations:

* without considering a `devcontainer.env` file using [.devcontainer/default/devcontainer.json](.devcontainer/default/devcontainer.json). Use this if you don't need to debug in the container.
* loading a `.env` file using [.devcontainer/withenvfile/devcontainer.json](.devcontainer/withenvfile/devcontainer.json). This configuration expects a file called `devcontainer.env` in the folder `.devcontainer`, which is needed for debugging.

> **Note** - `.env` files are excluded from git via `.gitignore`. You can use the file to store the environment variables `BTP_USERNAME` and `BTP_PASSWORD` that are needed when developing tests.

> **Note** - In the first run, the download of the container might take a while, so maybe time to grab a cup of coffee ☕.

## Local Setup

Ensure you have the following tools installed on your local machine.

* [git](https://git-scm.com/)
* [go](https://go.dev/)
* [golangci-lint](https://github.com/golangci/golangci-lint)
* [make](https://www.gnu.org/software/make/)
* [terraform](https://www.terraform.io/)

### MacOS (Homebrew)

If you run on MacOS, you can use [homebrew](https://brew.sh/) to speed up the installation process:

```bash
brew install git golang golangci-lint make terraform
```

### Windows (Chocolatey)

Windows users can let [chocolatey](https://chocolatey.org/) take over the installation for them:

```bash
choco install git golang golangci-lint make terraform
```

### Configuration of the Terraform CLI

Next you need to setup local development overrides in the Terraform CLI according to [this documentation](https://developer.hashicorp.com/terraform/plugin/debugging#terraform-cli-development-overrides). Once in place, Terraform will only consider local development builds for this provider.

Keep in mind that the [configuration file](https://developer.hashicorp.com/terraform/cli/config/config-file) location depends on your operating system:

* Mac/Linux/WSL: `~/.terraformrc`
* Windows: `%APPDATA%/terraform.rc`

The configuration should look similar to this:

```hcl
provider_installation {
  dev_overrides {
    "sap/btp" = "/path/to/go/bin" # the GOBIN directory can be found in the folder which `go env GOPATH` returns
  }

  direct {}
}
```

### Cloning of the Repository

The last step is then to clone the repository on your machine via:

```bash
git clone https://github.com/SAP/terraform-provider-btp.git
```

Navigate into the directory of the cloned repository.

## Install the Terraform Provider for SAP BTP Locally

Run the following command to build and install the provider:

```bash
make install
```

## Verify the Setup

Next, we verify that Terraform has access to your local build of the provider. Please navigate to one of the workspaces in the `examples` directory, e.g.:

```bash
cd examples/subaccount/
```

If you are now able to validate the scripts, everything is correctly set up:

```bash
terraform validate
```

In case of errors, please check first that you executed the previous steps correctly. If you are still stuck, feel free to ask for support by raising a [question](https://github.com/SAP/terraform-provider-btp/discussions/categories/q-a) in the [GitHub Discussions](https://github.com/SAP/terraform-provider-btp/discussions/categories/q-a) of this repository.
