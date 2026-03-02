# SECRETS for the `READY_PUSHER`

The [TakT](https://www.lakruzz.com/stories/takt/) workflow which allows you to run a smooth Pull-Request free flow needs access to merge stuff into your trunk (`main`).

GitHub flows comes with a built-in worker-bee token. In the flows you can access it in `${{ secrets.GITHUB_TOKEN }}`

> [!CAUTION]
> While the ${{ secrets.GITHUB_TOKEN }} can be lifted to `contents:write` permissions. GitHub actions treats
> commits created with this token as a special case which:
> **These commits will not trigger other workflows**

We need that! The `ready.yml` flow is designed to trigger the `stage.yml` flow.

## Solution

- In you developer profile settings Go and create a [**Fine-grained personal access token**](https://github.com/settings/personal-access-tokens)
- Settings (recommended)
  - **Token name**: `TAKT_CONTENT_WRITE`
  - **Description**: `Used to support the workflows devx-cafe/gh-tt workflow`
  - **Ressource owner**: `<YOUR-ORGANIZATION>` (recommended) or `<USER>` if you don't have an organization
  - **Expiration**: `<A-LONG-TIME>`
  - **Repository access** `All` or `Only selected`
  - **Permissions**: `contents:write` _(nothing else is needed)_

When you are done you will be presented a `gho****` token (GitHub oAuth) \_this is the only time you'll see this, but possibly not the only time you'll need it! So store it in your favorite password wallet.

<!-- cspell:ignore Cemi Okxps  -->

> [!TIP]
> You can not store this in git!
> If you do, GitHub will see 👀 it in the security scans going on in the background and
> and it will revoke the token!
> IF you want to store it in git, encode it with `base64` first and decode it when you need to use it:
> Example:
> If your token is `gho_1ZCCemiYAvkChNLJ4zOkxpsBh6X7FUZn25`
>
> ```bash
> $ echo gho_1ZCCemiYAvkChNLJ4zOkxpsBh6X7FUZn25 | base64
> Z2hvXzFaQ0NlbWlZQXZrQ2hOTEo0ek9reHBzQmg2WDdGVVpuMjUK
> echo Z2hvXzFaQ0NlbWlZQXZrQ2hOTEo0ek9reHBzQmg2WDdGVVpuMjUK | base64 --decode
> gho_1ZCCemiYAvkChNLJ4zOkxpsBh6X7FUZn25
> ```

## Make it a Secret

### 1. Organization secret

Go to the organization **"Action secrets and variables"** page

`<https://github.com/organizations/<ORG>/settings/secrets/actions>`

Create a new organisation secret:

- **Name**: `READY_PUSHER`
- **Secret**: `<TAKT_CONTENT_WRITE-VALUE>` (starts with `gho***`)

### 2. Repo secret:

Go to the repo **"Action secrets and variables"** page:

`<https://github.com/<USER|ORG>/<REPO-NAME>/settings/secrets/actions>`

Create a new secret:

- **Name**: `READY_PUSHER`
- **Secret**: `<TAKT_CONTENT_WRITE-VALUE>` (starts with `gho***`)

## Use the secret

Regardless if you created the secret as a _repo_ or an _organizations_ wide secret you use it the same way.

When you need to manipulate a repo, push it back to origin and have possible flows trigger on the change, simply pass `${{ secrets.READY_PUSHER }}` at the token to the `checkout` action

```yaml
  ...

jobs:
  some-job:

    ...

    permissions:
      contents: write # Required if you plan to push back _ even with your own token

  steps:

  ...

      - uses: actions/checkout@v6
        with:
          fetch-depth: 0   # Required if you plan to push back - Fetch full history
          token: ${{ secrets.READY_PUSHER }} # Required if you want push backs to trigger other flows must have contents:write

```

See it live in the `ready.yml` flow
