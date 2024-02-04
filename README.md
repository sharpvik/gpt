# `gpt` - ChatGPT 4 In Your Terminal

https://github.com/sharpvik/gpt/assets/23066595/8360b5c8-7e7d-4251-9cf2-9a770d807739

## Install

```bash
go install github.com/sharpvik/gpt
```

## Provide OpenAI API Key

```bash
gpt key <OPENAI_API_KEY>
```

## Ask a Quick Question

```bash
gpt "Tell me about football"
```

## Boot Up the REPL

```bash
gpt
```

```txt
👾
How are you?
^D
🤖
As an artificial intelligence, I don't have feelings, but I'm here and ready to assist you!
```

1. Use `CTRL+D` (EOF) to finish query input and send it to ChatGPT.
2. Use `CTRL+C` to leave the REPL environment.

## Copy Last Response

```bash
gpt copy
```
