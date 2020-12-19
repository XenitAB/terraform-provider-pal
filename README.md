# Terrraform Provider PAL
Provider to configure [Partner Admin Link](https://docs.microsoft.com/en-us/azure/cost-management-billing/manage/link-partner-id) for indiviual Azure Service Principals.

## How To
There is a example in [examples/basic/main.tf](./examples/basic/main.tf) that creates a Azure Service Principal with a random generated password.
The Service Principal is then used to authenticate with the random password to create the Partner Admin Link.

## Contributing
Follow the [contribution guide](./CONTRIBUTING.md) to get started contributing.
