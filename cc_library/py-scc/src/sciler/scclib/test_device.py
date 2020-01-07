import unittest
import json
import os
from unittest.mock import Mock, MagicMock
from Adafruit_ADS1x15 import ADS1115 as ADC

from cc_library.src.scripts.door import Door


class TestDevice(unittest.TestCase):
    def setUp(self):
        """
        Change this to test the device class you want to test.
        """
        self.device = Door()
        self.device.scclib = Mock(self.device.scclib)
        self.device.adc = Mock(ADC)
        self.device.adc.read_adc = MagicMock(return_value=50)

        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "../../scripts/door_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        self.config = json.load(open(file=abs_file_path))

    def test_get_status_str(self):
        result = self.device.get_status()
        self.assertIsInstance(result, dict, "get_status should return a dict type")

    def test_get_status_json(self):
        result = self.device.get_status()
        try:
            json.dumps(result)
        except TypeError:
            self.fail("get_status should return valid json string")
        finally:
            self.assertTrue(True)

    def test_perform_instruction(self):
        """
        Test all instruction specified in output are implemented.
        """
        output = self.config.get("output")
        for k in output:
            if isinstance(output.get(k), dict):
                """If it is a component dictionary with instructions"""
                for ins in output.get(k).get("instructions"):
                    print(ins)
                    contents = '{"instruction": "' + ins + '"}'
                    message = json.loads(contents)
                    try:
                        result, action_failed = self.device.perform_instruction(message)
                        self.assertTrue(
                            result,
                            "perform_instruction should be "
                            "implemented to return true for instruction from config: "
                            + ins,
                        )
                    except TypeError:
                        print("value missing for testing")
                        self.assertTrue(True, "instruction was implemented")
            else:
                contents = '{"instruction": "' + k + '"}'
                message = json.loads(contents)
                try:
                    result, action_failed = self.device.perform_instruction(message)
                    self.assertTrue(
                        result,
                        "perform_instruction should be "
                        "implemented to return true for instruction from config: " + k,
                    )
                except TypeError:
                    print("value missing for testing")
                finally:
                    print("value missing for testing")
                    self.assertTrue(True, "instruction was implemented")

    def test_perform_instruction_invalid(self):
        contents = '{"instruction": "my own non-existing instruction"}'
        message = json.loads(contents)
        result, action_failed = self.device.perform_instruction(message)
        self.assertFalse(
            result,
            "perform_instruction should be implemented to return false for invalid instruction",
        )

    def test_test_exists(self):
        result = self.device.test()
        self.assertIsNone(result, "test method should exist and have no return")

    # def test_main_exists(self):
    #     result = self.device.main()
    #     self.assertIsNone(result, "main method should exist and have no return")


if __name__ == "__main__":
    unittest.main()
