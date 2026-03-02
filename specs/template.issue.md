---
title: # *required* (text)
assign: [] # _optional_ (list of text) Assign people by their login. Use "@me" to self-assign.
labels: # _optional_ (list of tuples) Add labels by name
  - name: # *required* (text) Label name
    color: # _optinoal_ (text) Color of the label
    desc: # _optional_ (text) Description of the label
milestone: # _optional_ (text) Add the issue to a milestone by name
projects: # _optional_ (list of text) Add the issue to projects by title
---

## This is a sample issue instance template

It consists of two parts:

- Front Matter
- MarkDown body (content)

It supports the basic features of the `gh issue create` command.

It's designed to have a dedicated format `*.issue.md` as exemplified in the Front Matter above.

When a file in this format is passed to `mkissue` it will create an issue in the repo where it's executed, based on the Front Matter and markDown content.

## `assign`

Logins are typed without the `@` prefix, and exception to this rule is `@me` which is used as an abstraction for the user who executes the command.

Valid:

```yaml
assign: ["lakruzz", "@me"]
```

Valid:

```yaml
assign:
  - lakruzz
  - "@me"
```

Invalid:

```yaml
assign: ["@lakruzz", "@me"]
```

Valid:

```yaml
assign:
  - @lakruzz
  - @me
```

## `labels``

Is a list of YAML. Each item _must_ at least define `name` define the rest are optional.

Valid:

```yaml
labels:
  - name: "Help Wanted"
```

When only `name` is given, it's implied that the label _must exist_ already ...or the creation will fail.

The list YAML also supports `color` and `desc`. They are both _optional_ but if _any_ of them are given, it's implied, that the label should be created, if it doesn't exist.

<details>
<summary>Logic:</summary>

```shell
# If listing the label fails, then create it:
$LABEL_NAME=<label.mame>
$LABEL_DESC=<label.desc>
$LABEL_COLOR=<label.color>

gh label list --json name --jq '.[].name' | grep "$LABEL"
$? && gh label create "$LABEL_NAME" -c $LABEL_COLOR -d "$LABEL_DESC"
```

</details>

## `milestone``

The `milestone` is the name of an existing milestone. The setting is optional, but if it is given then the miles sone must exist already or the creation will fail.

Valid:

```yaml
milestone: "Some feature"
```

## `projects``

The `projects` setting is a list of project titles. The setting is optional but if it's given then all projects must exist already or the creation will fail.

Valid:

```yaml
projects: ["Kanban upstream", "kanban downstream"]
```

Valid:

```yaml
projects:
  - "Kanban upstream"
  - "kanban downstream"
```
