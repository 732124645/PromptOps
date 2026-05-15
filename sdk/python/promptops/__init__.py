"""PromptOps Python SDK — fetch prompts with WebSocket hot-reload.

Standard library only, no third-party dependencies.
"""

from .client import PromptOpsClient, render_template

__all__ = ["PromptOpsClient", "render_template"]
__version__ = "0.1.0"
