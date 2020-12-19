# Contributing
This guide will help you get setup to make changes to terraform-provider-pal.

## How To
Install the required tools.
```
make tools
```

When developing you may want to test the provider locally. You can install the provider locally
so that it can be used by any terraform project on your computer. The default version will be `0.0.0-dev`.
```
make install
```

Run the provider unit tests.
```
make test
```

Run the provider acceptance tests.
```
make testacc
```

If you have made changes to the schema, including a description you need to generate the docs.
```
make docs
```
