import requests

response = requests.post('http://localhost:8080/api/cmd', json={'cmd': 'ls'})
print(response.text)
