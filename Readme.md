# Wingman

Your second pilot in the terminal.

## What is Wingman?

Wingman is a tool that helps you to generate shell commands using OpenAI GPT-3 model.

## How to use?

### Install

Just download the binary from the [releases](https://github.com/harnyk/wingman/releases) page and put it in your `$PATH`.

Alternatively, you can install it using [eget](https://github.com/zyedidia/eget):

```bash
eget harnyk/wingman
```

### Setup

You need to get an API key from [OpenAI](https://openai.com/).

Then, you need to create a file `~/.wingman.yaml` with the following content:

```yaml
openai_token = <your key>
```

### Usage

```bash
wingman <your query>
# For example:
wingman display weather forecast in Kyiv
```

[View Demo](https://asciinema.org/a/570008)
