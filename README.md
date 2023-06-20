# mach-composer-plugin-gcp

An unopinionated plugin for GCP


## Usage

```yaml
mach_composer:
  version: 1
  plugins:
    gcp:
      source: mach-composer/gcp
      version: 0.0.1

global:
  environment: test
  cloud: gcp
  terraform_config:
    gcs:
      bucket: "The name of the GCS bucket."
      prefix: "(optional) GCS prefix inside the bucket."
  gcp:
    project: "12345678910"
    region: "regiona"
    zone: "zonea"
    beta: True

sites:
  - identifier: my-site
    gcp:
      project: "siteproject"
      region: "siteregion"
      zone: "sitezone"
```
