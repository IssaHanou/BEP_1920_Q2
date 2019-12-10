import unittest
import json
from unittest.mock import Mock
from cc_library.src.sciler.scclib.app import SccLib
from cc_library.src.scripts.door import Door

class TestDevice(unittest.TestCase):

    def setUp(self):
        """
        Change this to test the device class you want to test.
        """
        self.device = Door()
        self.device.scclib = Mock(self.device.scclib)
        print("1 " + str(self.device.scclib))
        self.device.scclib.status_changed = Mock()
        print(self.device.scclib.status_changed())

    def test_get_status_str(self):
        result = self.device.get_status()
        self.assertIsInstance(result, str, "get_status should return a string type")

    def test_get_status_json(self):
        result = self.device.get_status()
        try:
            json.dumps(result)
        except TypeError:
            self.fail("get_status should return valid json string")
        finally:
            self.assertTrue(True)

    def test_perform_instruction_test(self):
        print("2:" + str(self.device.scclib))
        contents = '{"instruction": "test"}'
        message = json.loads(contents)
        result = self.device.perform_instruction(message)
        self.assertTrue(result, "perform_instruction should be implemented to return true for test instruction")

    def test_perform_instruction_invalid(self):
        contents = '{"instruction": "my own random instruction"}'
        message = json.loads(contents)
        result = self.device.perform_instruction(message)
        self.assertFalse(result, "perform_instruction should be implemented to return false for invalid instruction")

    def test_test_exists(self):
        result = self.device.test()
        self.assertIsNone(result, "test should have no return")
    #
    # def test_main_exists(self):
    #     result = self.device.main()
    #     self.assertIsNone(result, "main should have no return")


if __name__ == "__main__":
    unittest.main()
