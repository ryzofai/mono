from flask import Flask, jsonify, request, make_response, app, redirect, render_template, flash, redirect, url_for, session, request, logging
#import flask_login
import datetime
import jwt
from wtforms import Form, StringField, TextAreaField, PasswordField, validators
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
			return redirect('login')
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
	
@app.route('/login2')
def login2():
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

@app.route('/login', methods=['GET', 'POST'])
def login():
	if request.method == 'POST':
        # Get Form Fields
		username = request.form['username']
		password_candidate = request.form['password']

        # Create cursor
        #cur = mysql.connection.cursor()

        # Get user by username
        #result = cur.execute("SELECT * FROM users WHERE username = %s", [username])

		#if result > 0:
            # Get stored hash
            #data = cur.fetchone()
            #password = data['password']

            # Compare Passwords
        
		if (username == 'user' and password_candidate == 'secret'):
            # Passed
            #session['logged_in'] = True
            #session['username'] = username
			#auth = request.authorization
			token = jwt.encode({'user' : username, 'exp' : datetime.datetime.utcnow() + datetime.timedelta(seconds=10)}, app.config['SECRET_KEY'])
			#flash('You are now logged in', 'success')
			return redirect('protected'), {'Set-Cookie': 'token=' + token.decode('UTF-8') + '; Max-Age=10'}
		else:
			error = 'Invalid login'
			return render_template('login.html', error=error)
        # Close connection
		# cur.close()
		#else:
		#	error = 'Username not found'
		#	return render_template('login.html', error=error)

	return render_template('login.html')


if __name__ == '__main__':
	app.run(debug=True)
    #app.run(debug=True, ssl_context='adhoc')
