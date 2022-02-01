# Argo Workflows URL Finder
When Argo Workflows archives a workflow, the URL changes, and the URL cannot be pre-determined. If you're using tools such as [the ci-github-notifier](https://github.com/sendible-labs/ci-github-notifier) to annotate results onto Github Pull Requests, this change of URL quickly becomes an issue.

The Argo Workflows URL Finder simply locates your workflow by querying the Argo Workflows API and redirects the user to the correct URL.

![CI](https://github.com/sendible-labs/argo-workflows-url-finder/actions/workflows/ci.yaml/badge.svg) ![Code Quality](https://github.com/sendible-labs/argo-workflows-url-finder/actions/workflows/codeql-analysis.yaml/badge.svg) ![Release](https://github.com/sendible-labs/argo-workflows-url-finder/actions/workflows/release.yaml/badge.svg)


# Environment Variables
We pass key information to the container using environment variables.

| Environment Variable  | Type      | Description                                                                                                                                                                                                 |
|---------------------- |---------- |------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `ARGO_URL`            | string    | (Mandatory): The url to your argo-workflows installation. eg. https://argo-workflows.sendible.com                                                                                                         |
| `ACCESS_TOKEN`        | string    | (Optional if `TOKEN_FILE` is set): An Argo Workflows token with permission to query the API.                                                                                                                |
| `TOKEN_FILE`          | string    | (Optional if `ACCESS_TOKEN` is set): Path to a file within the container that contains the Argo Workflows access token. Useful for Vault secrets injection or similar. Takes precedence over `ACCESS_TOKEN` |


# Hosting
The Argo Workflows URL Finder can be deployed anywhere you can deploy a container. However, if you're running Argo Workflows already, you probably use Kubernetes. We have provided a [basic example deployment](example/k8s-deployment.yaml).

For testing purposes, a simple docker run command could be:
```
docker run \
-e ARGO_URL=https://argo-workflows.sendible.com/ \
-e ACCESS_TOKEN="xxxx" \
-p 8080:8080 \
ghcr.io/sendible-labs/argo-workflows-url-finder:stable
```

# Putting Argo Workflows URL Finder to work

Assuming we have performed the docker run command above and we have a workflow called `wonderful-whale` in the namespace `testing`. We would query for its location by going to `http://localhost:8080?workflowname=wonderful-whale&namespace=testing`.

The user will be redirected to the workflow regardless of whether it's archived or not.