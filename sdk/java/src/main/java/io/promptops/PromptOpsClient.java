package io.promptops;

import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.http.WebSocket;
import java.nio.charset.StandardCharsets;
import java.util.Map;
import java.util.concurrent.CompletionStage;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Client for the PromptOps runtime API. Fetches prompts by key, caches them,
 * and (via {@link #watch}) keeps the cache fresh over a WebSocket connection.
 *
 * <p>Depends only on the JDK ({@code java.net.http}).
 */
public final class PromptOpsClient {

    private static final String DEFAULT_TOKEN = "promptops-dev-token";
    private static final Pattern VAR = Pattern.compile("\\{\\{\\s*([\\w.]+)\\s*}}");

    private final String server;
    private final String namespace;
    private final String token;
    private final String appName;
    private final HttpClient http = HttpClient.newHttpClient();
    private final Map<String, Prompt> cache = new ConcurrentHashMap<>();
    private volatile WebSocket webSocket;

    public PromptOpsClient(String server) {
        this(server, "prod", DEFAULT_TOKEN);
    }

    public PromptOpsClient(String server, String namespace, String token) {
        this(server, namespace, token, "");
    }

    /**
     * @param appName optional app name reported to the server's
     *                connected-clients registry.
     */
    public PromptOpsClient(String server, String namespace, String token, String appName) {
        if (server == null || server.isEmpty()) {
            throw new IllegalArgumentException("PromptOps: \"server\" is required");
        }
        this.server = server.replaceAll("/+$", "");
        this.namespace = namespace;
        this.token = token;
        this.appName = appName == null ? "" : appName;
    }

    /** Replace {@code {{variable}}} placeholders; unknown variables stay as-is. */
    public static String renderTemplate(String content, Map<String, ?> variables) {
        Matcher matcher = VAR.matcher(content == null ? "" : content);
        StringBuilder out = new StringBuilder();
        while (matcher.find()) {
            String key = matcher.group(1);
            Object value = variables == null ? null : variables.get(key);
            String replacement = value != null ? value.toString() : matcher.group(0);
            matcher.appendReplacement(out, Matcher.quoteReplacement(replacement));
        }
        matcher.appendTail(out);
        return out.toString();
    }

    public Prompt getPrompt(String key) {
        return getPrompt(key, false);
    }

    public Prompt getPrompt(String key, boolean refresh) {
        if (!refresh) {
            Prompt cached = cache.get(key);
            if (cached != null) {
                return cached;
            }
        }
        String url = server + "/api/sdk/prompts/"
                + URLEncoder.encode(key, StandardCharsets.UTF_8)
                + "?env=" + URLEncoder.encode(namespace, StandardCharsets.UTF_8);
        try {
            HttpRequest request = HttpRequest.newBuilder(URI.create(url))
                    .header("Authorization", "Bearer " + token)
                    .GET()
                    .build();
            HttpResponse<String> response = http.send(request, HttpResponse.BodyHandlers.ofString());
            if (response.statusCode() != 200) {
                throw new RuntimeException(
                        "PromptOps: failed to fetch \"" + key + "\" (HTTP " + response.statusCode() + ")");
            }
            @SuppressWarnings("unchecked")
            Map<String, Object> body = (Map<String, Object>) Json.parse(response.body());
            Prompt prompt = new Prompt(
                    str(body, "key"), str(body, "version"), str(body, "env"),
                    str(body, "model"), str(body, "content"));
            cache.put(key, prompt);
            return prompt;
        } catch (RuntimeException e) {
            throw e;
        } catch (Exception e) {
            throw new RuntimeException("PromptOps: request failed for \"" + key + "\"", e);
        }
    }

    public String render(String key, Map<String, ?> variables) {
        return renderTemplate(getPrompt(key).content(), variables);
    }

    /**
     * Open a WebSocket to the server and refresh any cached prompt whose key
     * appears in a hot-reload event. {@code onUpdate} receives the prompt key.
     */
    public void watch(Consumer<String> onUpdate) {
        // Authenticate and identify this connection to the server.
        String query = "token=" + URLEncoder.encode(token, StandardCharsets.UTF_8)
                + "&client=java-sdk&namespace="
                + URLEncoder.encode(namespace, StandardCharsets.UTF_8);
        if (!appName.isEmpty()) {
            query += "&app=" + URLEncoder.encode(appName, StandardCharsets.UTF_8);
        }
        String wsUrl = server.replaceFirst("^http", "ws") + "/ws?" + query;
        this.webSocket = http.newWebSocketBuilder()
                .buildAsync(URI.create(wsUrl), new HotReloadListener(onUpdate))
                .join();
    }

    public void close() {
        WebSocket ws = this.webSocket;
        if (ws != null) {
            ws.sendClose(WebSocket.NORMAL_CLOSURE, "client closed");
            this.webSocket = null;
        }
    }

    private static String str(Map<String, Object> map, String key) {
        Object value = map.get(key);
        return value == null ? null : value.toString();
    }

    /** Buffers WebSocket text parts and refreshes the cache on each event. */
    private final class HotReloadListener implements WebSocket.Listener {

        private final Consumer<String> onUpdate;
        private final StringBuilder buffer = new StringBuilder();

        HotReloadListener(Consumer<String> onUpdate) {
            this.onUpdate = onUpdate;
        }

        @Override
        public void onOpen(WebSocket ws) {
            ws.request(1);
        }

        @Override
        public CompletionStage<?> onText(WebSocket ws, CharSequence data, boolean last) {
            buffer.append(data);
            if (last) {
                String message = buffer.toString();
                buffer.setLength(0);
                handle(message);
            }
            ws.request(1);
            return null;
        }

        private void handle(String message) {
            try {
                @SuppressWarnings("unchecked")
                Map<String, Object> event = (Map<String, Object>) Json.parse(message);
                String key = str(event, "key");
                if (key != null && cache.containsKey(key)) {
                    getPrompt(key, true);
                    if (onUpdate != null) {
                        onUpdate.accept(key);
                    }
                }
            } catch (RuntimeException ignored) {
                // Malformed events must not break the connection.
            }
        }
    }
}
