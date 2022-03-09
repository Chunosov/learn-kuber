# basic example for kustomize

https://medium.com/nerd-for-tech/getting-started-with-kustomize-bff87e820cde

Kustomize follows overlay mechanism. Unlike Helm it doesn't replace template placeholdes, it replaces parts of text in base files, like `sed`.

- It can generate resources from other sources
- It can customize the collection of resources

Kustomize reads base file and put overlays onto them. Base files stay unchanged.

`kustomization.yaml` informs Kustomize on how to render the resources:

```bash
cd base
kubectl apply -k . --dry-run=client -o yaml
```

We will create 2 more overlay folders to manage deployment environment wise. By this way we will enhance our base with some modification.

The overlay mechanism will allow us to append the information on base YAMLs without changing anything on them.

```bash
cd ../prod
kubectl apply -k . --dry-run=client -o yaml
```