package io.promptops;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Minimal JSON parser (objects, arrays, strings, numbers, booleans, null).
 * Sufficient for PromptOps API responses; keeps the SDK dependency-free.
 */
final class Json {

    private final String src;
    private int pos;

    private Json(String src) {
        this.src = src;
    }

    static Object parse(String text) {
        Json json = new Json(text);
        json.skipWhitespace();
        Object value = json.readValue();
        json.skipWhitespace();
        return value;
    }

    private Object readValue() {
        char c = src.charAt(pos);
        switch (c) {
            case '{':
                return readObject();
            case '[':
                return readArray();
            case '"':
                return readString();
            case 't':
                pos += 4;
                return Boolean.TRUE;
            case 'f':
                pos += 5;
                return Boolean.FALSE;
            case 'n':
                pos += 4;
                return null;
            default:
                return readNumber();
        }
    }

    private Map<String, Object> readObject() {
        Map<String, Object> map = new LinkedHashMap<>();
        pos++; // consume '{'
        skipWhitespace();
        if (src.charAt(pos) == '}') {
            pos++;
            return map;
        }
        while (true) {
            skipWhitespace();
            String key = readString();
            skipWhitespace();
            pos++; // consume ':'
            skipWhitespace();
            map.put(key, readValue());
            skipWhitespace();
            if (src.charAt(pos++) == '}') {
                break;
            }
        }
        return map;
    }

    private List<Object> readArray() {
        List<Object> list = new ArrayList<>();
        pos++; // consume '['
        skipWhitespace();
        if (src.charAt(pos) == ']') {
            pos++;
            return list;
        }
        while (true) {
            skipWhitespace();
            list.add(readValue());
            skipWhitespace();
            if (src.charAt(pos++) == ']') {
                break;
            }
        }
        return list;
    }

    private String readString() {
        StringBuilder out = new StringBuilder();
        pos++; // consume opening '"'
        while (true) {
            char c = src.charAt(pos++);
            if (c == '"') {
                break;
            }
            if (c != '\\') {
                out.append(c);
                continue;
            }
            char esc = src.charAt(pos++);
            switch (esc) {
                case '"':
                    out.append('"');
                    break;
                case '\\':
                    out.append('\\');
                    break;
                case '/':
                    out.append('/');
                    break;
                case 'n':
                    out.append('\n');
                    break;
                case 't':
                    out.append('\t');
                    break;
                case 'r':
                    out.append('\r');
                    break;
                case 'b':
                    out.append('\b');
                    break;
                case 'f':
                    out.append('\f');
                    break;
                case 'u':
                    out.append((char) Integer.parseInt(src.substring(pos, pos + 4), 16));
                    pos += 4;
                    break;
                default:
                    out.append(esc);
            }
        }
        return out.toString();
    }

    private Object readNumber() {
        int start = pos;
        while (pos < src.length() && "-+.eE0123456789".indexOf(src.charAt(pos)) >= 0) {
            pos++;
        }
        String number = src.substring(start, pos);
        if (number.indexOf('.') >= 0 || number.indexOf('e') >= 0 || number.indexOf('E') >= 0) {
            return Double.parseDouble(number);
        }
        return Long.parseLong(number);
    }

    private void skipWhitespace() {
        while (pos < src.length() && Character.isWhitespace(src.charAt(pos))) {
            pos++;
        }
    }
}
