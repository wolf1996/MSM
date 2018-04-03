import requests
import unittest

prefix = 'http://localhost:8080'
class TestCase(unittest.TestCase):
    def setUp(self):
        self.app = requests.Session()
        res = self.app.post(prefix +'/user/sign_in', """{
	        "email": "ml@gmail.com",
	        "password": "123456"
        }""")
        print("login is \n")
        print(res.text)

    def test_user_info(self):
        resp = self.app.get(prefix+'/user/user_info')
        print('User Info:\n')
        print(resp.text)

    def test_get_user_controllers(self):
        resp = self.app.get(prefix+'/controller/get_user_controllers')
        print('User Controllers: \n')
        print(resp.text)
    
    def test_get_conroller_sensors(self):
        resp = self.app.get(prefix+'/controller/1/get_sensors')
        print('Controller Sensors: \n')
        print(resp.text)

    def test_get_controller_stats(self):
        resp = self.app.get(prefix+'/controller/1/get_controller_stats')
        print('Controller Stats \n')
        print(resp.text)

    def test_get_sensor_stats(self):
        resp = self.app.get(prefix+'/sensor/1/view_stats')
        print('Sensor Stats\n')
        print(resp.text)
    
    def test_get_sensor_data(self):
        resp = self.app.get(prefix+'/sensor/1/get_data')
        print('Sensor Data\n')
        print(resp.text)
    
    def test_get_user_object(self):
        resp = self.app.get(prefix + '/object/get_user_objects')
        print('Get User Objects\n')
        print(resp.text)

    def test_get_obj_controllers(self):
        resp = self.app.get(prefix + '/object/1/get_object_controllers')
        print('Get Object Controllers\n')
        print(resp.text)

    def test_get_obj_stats(self):
        resp = self.app.get(prefix + '/object/1/get_object_stats')
        print('Get Object Stats\n')
        print(resp.text)

if __name__ == '__main__':
    unittest.main()