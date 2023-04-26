import requests
from flask import Flask, redirect, request

app = Flask(__name__)

clientID = ""
clientSecret = ""
redirectURI = "http://localhost:8080/me"
scope = "account"
state = "12345"

@app.route("/", methods=['GET'])
def index():
    return redirect(f'https://oauth.vk.com/authorize?response_type=token&client_id={clientID}&redirect_uri={redirectURI}&scope={scope}&state={state}')

@app.route("/me", methods=['GET'])
def me():
    token = request.args.get('access_token')
    print(token)
    r = requests.get(f'https://api.vk.com/method/users.get?v=5.124&access_token={token}')
    return r.text


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=False)
