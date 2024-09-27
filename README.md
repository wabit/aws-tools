# aws-tools

## Project Overview

**aws-tools** is a CLI tool to run tedious AWS commands

### Features

- **Log into AWS**: easily see which profiles you have configued, and with a quick selection lets you log in.
- **Log into EKS**: manage which clusters use what profile and easily update your kube config.
- **Cross-Platform**: Compatible with both macOS, and Linux.

### Getting Started

To get started with aws-tools, follow the installation instructions and explore the available commands to see how it can improve your workflow.

### Build

To build the project, run the following command:

```bash
go build -o aws-tools
```

### Usage

```bash
./aws-tools
```

### Configuration

On first run, the tool will promt you to enter the command you use to login to AWS. This will be saved in a configuration file in your home directory (~/.aws-tools) where you would normaly pass the profile to the command repalce it with `PROFILE_NAME`.

e.g. `aws-google-login -p PROFILE_NAME`

#### EKS Clusters

To configure the EKS clusters select `Add EKS Cluster` from the configuration menu. You will be prompted to enter the cluster name and then the Profile to use to connect to the cluster. This will be saved in the configuration file.

You will then be able to select the cluster from the `EKS Login` menu.

### Contributing

Contributions are welcome! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Contact

For any questions or feedback, please open an issue on GitHub
