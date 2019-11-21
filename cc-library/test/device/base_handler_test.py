import unittest
from device import base_handler


class BaseHandlerTestCase(unittest.TestCase):
    def test_on_change(self):
        handler = base_handler.BaseHandler()
        handler.on_change("on")
        self.assertEqual(handler.status, "on")