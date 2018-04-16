---
title: Add your project
weight: 10
---

This will guide you on how to add your own project to this site.

If you use [GoReleaser], all you need to do is to add an empty YAML
file to the [GoDownloader] `tree` folder, under the right path.

For example, if your project lives under `https://github.com/foo/bar`, you
may create the YAML file at `tree/github.com/foo/bar.yaml`.

After that, you just need to open a pull request, and everything will happen
magically.

If you do not use [GoReleaser], additional steps will be required.
<!-- TODO: document additional steps here -->

[GoReleaser]: https://goreleaser.com
[GoDownloader]: https://github.com/goreleaser/godownloader

You can see the list of all projects being served [here](/projects).
