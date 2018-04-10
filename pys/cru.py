from flask import Flask, jsonify, request, make_response, session, app, redirect
import flask_login
import datetime
import jwt
from functools import wraps

app = Flask(__name__)

app.config['SECRET_KEY'] = 'thisisthesecretkey'

#@app.before_request
#def make_session_permanent():
#    session.permanent = True
#    app.permanent_session_lifetime = datetime.timedelta(seconds=10)

def token_required(f):
	@wraps(f)
	def decorated(*args, **kwargs):
		req = request.get_json()
		#token = req['token']
		token =  request.cookies.get('token')
		if not token:
			return redirect("/login", code=401)
		try: 
			data = jwt.decode(token, app.config['SECRET_KEY'])
		except:
			return jsonify({'message' : 'Token is invalid!'}), 403

		return f(*args, **kwargs)

	return decorated

@app.route('/unprotected')
def unprotected():
	return jsonify({'message' : 'test!'}), 200, {'Set-Cookie': 'testur=tester; Max-Age=10'}
	#return jsonify({'message' : 'Anyone can view this!'})

@app.route('/getcookie')
def getcookie():
	cooki =  request.cookies.get('testur')
	return cooki
	
@app.route('/protected', methods=['GET', 'POST'])
@token_required
def protected():
	req = request.get_json()
	return jsonify({'message' : 'Valid token, you can view this page.'})

@app.route('/protected2', methods=['GET', 'POST'])
@token_required
def protected2():
	req = request.get_json()
	return jsonify({'message' : 'protected 2.'})
	
@app.route('/login')
def login():
	#req = request.get_json()
	#token_instance =  request.cookies.get('token')
	
	#if not token_instance:
	#	return make_response('Could not verify!', 401, {'WWW-Authenticate' : 'Basic realm="Login Required"'})
	#else:
	auth = request.authorization
	if auth and auth.password == 'secret':
		token = jwt.encode({'user' : auth.username, 'exp' : datetime.datetime.utcnow() + datetime.timedelta(seconds=10)}, app.config['SECRET_KEY'])
		return jsonify({'token' : token.decode('UTF-8')}), 200, {'Set-Cookie': 'token=' + token.decode('UTF-8') + '; Max-Age=10'}
	return make_response('Could not verify!', 401, {'WWW-Authenticate' : 'Basic realm="Login Required"'})

if __name__ == '__main__':
	app.run(debug=True)
    #app.run(debug=True, ssl_context='adhoc')
