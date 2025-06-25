# Sync APIM artifacts between Tenant and GitHub with FlashPipe on GitHub Actions
The page describes the steps to set up _FlashPipe_ on [GitHub Actions](https://github.com/features/actions) to sync API Management artifacts between an Integration Suite tenant and a GitHub repository.

## A. Sync from Tenant to GitHub

### 1. Create GitHub repository
Create (or use) an existing repository on GitHub.

![New Repo](images/setup/git-sync/01_new_repo.png)
Ensure that the repository includes the following files at the root directory. The links provide samples for each file that can be used.

- [.gitignore](https://github.com/engswee/flashpipe-demo/blob/github-actions-sync/.gitignore) - ensures unwanted files are not included in commits
- [.gitattributes](https://github.com/engswee/flashpipe-demo/blob/github-actions-sync/.gitattributes) - ensures correct line endings for committed files

### 2. Create secrets in GitHub repository
Sensitive information can be stored securely on GitHub using [encrypted secrets](https://docs.github.com/en/actions/reference/encrypted-secrets). These can then be passed to the pipeline steps as environment variables. For _FlashPipe_, we will use these to securely store the details to access the tenant.

In the GitHub repository, go to `Settings` > `Secrets` to create new repository secrets as shown below.
![Secrets Setting](images/setup/github-actions/05a_secrets.png)

**OAuth Authentication**

Create the following repository secrets. Refer to [OAuth client setup page](oauth_client.md) for details on setting up the OAuth client for usage with _FlashPipe_. Note that API Management credentials differ from Cloud Integration, so ensure that the secret names do not have any conflicts.
1. `APIPORTAL_OAUTH_CLIENTID` - OAuth client ID
2. `APIPORTAL_OAUTH_CLIENTSECRET` - OAuth client secret
   ![OAuth Secrets](images/setup/git-sync-apim/02_oauth_secrets_apim.png)

### 3. Configure workflow permissions
In order for the workflows to be able to make changes to the repository, correct permissions need to be configured.

In the GitHub repository, go to `Settings` > `Actions` > `General`. Scroll down to the `Workflow permissions` section and select `Read and write permissions` and click `Save`.
![Workflow permissions](images/setup/git-sync/03c_workflow_permissions.png)

### 4. Create GitHub Actions workflow
In the GitHub repository, go to `Actions` to create new workflow.
![New Workflow](images/setup/git-sync/03a_new_workflow.png)

Skip the templates and choose `set up a workflow yourself`.

Provide a suitable name for the workflow file e.g. `sync-apim-to-git.yml` and replace the default content with the code sample below. Replace the tenant and authentication details accordingly - note that the host refers to the management node of the API Portal which differs from the tenant management node of Cloud Integration.

**NOTE** - FlashPipe comes with companion GitHub Action [engswee/flashpipe-action](https://github.com/engswee/flashpipe-action) that simplifies usage in a workflow. The following actions are used in the workflow:
- [engswee/flashpipe-action/sync/apiproxy@v1](https://github.com/engswee/flashpipe-action#sync-apiproxy)
- [engswee/flashpipe-action/sync/apiproduct@v1](https://github.com/engswee/flashpipe-action#sync-apiproduct)

![Sync Workflow](images/setup/git-sync-apim/04b_sync_apim_git_workflow.png)

<script src="https://gist.github.com/engswee/1a06264333060f8eb5c02690cb55957e.js"></script>

Save and commit the new workflow file.

Note: GitHub provides functionality to store unencrypted plain text as `repository variables`. Optionally, values like base URLs can be stored as repository variables instead of being hardcoded in the YAML configuration file, and can then be access using [the `vars` context](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/store-information-in-variables#using-the-vars-context-to-access-configuration-variable-values) in the configuration file.

### 5. Trigger workflow execution
This workflow has been configured with `on: workflow_dispatch` event triggering which allows it to be executed manually.

In the GitHub repository, go to `Actions`, select the workflow and click `Run workflow`.
![Execute Workflow](images/setup/git-sync-apim/05a_run_workflow_apim.png)

### 6. View execution results

During or upon completion of the workflow run, the logs can be viewed by clicking on the workflow run.
![Workflow Logs APIProxy](images/setup/git-sync-apim/06a_logs_apim.png)
![Workflow Logs APIProduct](images/setup/git-sync-apim/06a_logs_apiproduct.png)

The APIM artifact files have now been downloaded from the tenant and committed to the repository.
![IFlow Files](images/setup/git-sync-apim/06b_apim_files.png)

## B. Sync from GitHub to Tenant

### 1. Create GitHub Actions workflow
In the GitHub repository, go to `Actions` to create new workflow.
![New Workflow](images/setup/git-sync/03a_new_workflow.png)

Skip the templates and choose `set up a workflow yourself`.

Provide a suitable name for the workflow file e.g. `sync-apim-to-tenant.yml` and replace the default content with the code sample below. Replace the tenant and authentication details accordingly.

![Sync Workflow](images/setup/git-sync-apim/b_01_sync_apim_to_tenant_workflow.png)

<script src="https://gist.github.com/engswee/4271fbcef81db0ec719bc8cc46bc1d69.js"></script>

Save and commit the new workflow file.

### 2. Trigger workflow execution
This workflow has been configured with `on: workflow_dispatch` event triggering which allows it to be executed manually.

In the GitHub repository, go to `Actions`, select the workflow and click `Run workflow`.
![Execute Workflow](images/setup/git-sync-apim/b_02_run_workflow.png)

### 3. View execution results

During or upon completion of the workflow run, the logs can be viewed by clicking on the workflow run.
![Workflow Logs APIProxy](images/setup/git-sync-apim/b_03_sync_tenant_logs.png)
![Workflow Logs APIProduct](images/setup/git-sync-apim/b_03_sync_tenant_logs_apiproduct.png)

As seen from the logs, APIM artifacts were synced to the tenant.

_NOTE:_ Due to limitation of SAP's public API, updates of API Products are not synced to tenant, instead these will be skipped.