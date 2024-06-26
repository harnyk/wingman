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
openai_model = gpt-4o
```

`openai_model` can be omitted if you want to use the default model, which is `gpt-3.5-turbo`.

### Usage

```bash
wingman <your query>
# For example:
wingman display weather forecast in Kyiv
```

[![asciicast](https://asciinema.org/a/570008.svg)](https://asciinema.org/a/570008)
