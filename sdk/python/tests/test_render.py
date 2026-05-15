import os
import sys
import unittest

sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

from promptops import PromptOpsClient, render_template


class RenderTemplateTest(unittest.TestCase):
    def test_substitutes_known_variables(self):
        self.assertEqual(render_template("hello {{name}}", {"name": "Ada"}), "hello Ada")

    def test_leaves_unknown_variables(self):
        self.assertEqual(render_template("hello {{name}}", {}), "hello {{name}}")

    def test_whitespace_and_dotted_keys(self):
        self.assertEqual(render_template("{{ a.b }}", {"a.b": "X"}), "X")

    def test_repeated_occurrences(self):
        self.assertEqual(render_template("{{x}}-{{x}}", {"x": 1}), "1-1")

    def test_none_content(self):
        self.assertEqual(render_template(None, {}), "")


class ClientConstructionTest(unittest.TestCase):
    def test_requires_server(self):
        with self.assertRaises(ValueError):
            PromptOpsClient("")

    def test_strips_trailing_slash_and_defaults(self):
        client = PromptOpsClient("http://localhost:8080/")
        self.assertEqual(client.server, "http://localhost:8080")
        self.assertEqual(client.namespace, "prod")


if __name__ == "__main__":
    unittest.main()
