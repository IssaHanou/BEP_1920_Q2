import unittest

from cc_library.src.device import base_handler


class BaseHandlerTestCase(unittest.TestCase):
    def test_on_change(self):
        handler = base_handler.BaseHandler()
        handler.on_change("on")
        self.assertEqual(handler.status, "on")


if __name__ == "__main__":
    unittest.main()
