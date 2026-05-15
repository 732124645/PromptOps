package io.promptops;

/** A prompt fetched from the PromptOps runtime API. */
public record Prompt(String key, String version, String env, String model, String content) {
}
