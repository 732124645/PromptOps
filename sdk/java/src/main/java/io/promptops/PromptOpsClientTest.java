package io.promptops;

import java.util.Map;

/** Dependency-free test runner: exits non-zero on the first failed assertion. */
public final class PromptOpsClientTest {

    public static void main(String[] args) {
        check("hello Ada".equals(
                        PromptOpsClient.renderTemplate("hello {{name}}", Map.of("name", "Ada"))),
                "renderTemplate substitutes known variables");

        check("hello {{name}}".equals(
                        PromptOpsClient.renderTemplate("hello {{name}}", Map.of())),
                "renderTemplate leaves unknown variables untouched");

        check("X".equals(
                        PromptOpsClient.renderTemplate("{{ a.b }}", Map.of("a.b", "X"))),
                "renderTemplate handles whitespace and dotted keys");

        check("1-1".equals(
                        PromptOpsClient.renderTemplate("{{x}}-{{x}}", Map.of("x", 1))),
                "renderTemplate substitutes repeated occurrences");

        check("".equals(PromptOpsClient.renderTemplate(null, Map.of())),
                "renderTemplate handles null content");

        boolean threw = false;
        try {
            new PromptOpsClient("");
        } catch (IllegalArgumentException e) {
            threw = true;
        }
        check(threw, "constructor requires a server");

        Object parsed = Json.parse("{\"key\":\"k\",\"content\":\"line1\\nline2\",\"n\":3}");
        @SuppressWarnings("unchecked")
        Map<String, Object> map = (Map<String, Object>) parsed;
        check("line1\nline2".equals(map.get("content")), "Json parser decodes escapes");
        check("k".equals(map.get("key")), "Json parser reads string fields");

        System.out.println("All Java SDK tests passed");
    }

    private static void check(boolean condition, String name) {
        if (!condition) {
            System.err.println("FAIL: " + name);
            System.exit(1);
        }
        System.out.println("ok - " + name);
    }
}
