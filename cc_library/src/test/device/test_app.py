import unittest


class MyTestCase(unittest.TestCase):
    def test_something(self):
        self.assertEqual(True, False)


if __name__ == "__main__":
    unittest.main()


# import unittest

# from cc_library.src.sciler import base_handler


# class BaseHandlerTestCase(unittest.TestCase):
#    def test_on_change(self):
#        handler = base_handler.BaseHandler()
#        handler.on_change("on")
#        self.assertEqual(handler.status, "on")


# if __name__ == "__main__":
#    unittest.main()
