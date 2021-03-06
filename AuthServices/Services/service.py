#Importing pre-defined modules
import jwt
import re

#Importing Services
from Services.mongo_config import MongoConfig
from Services.crypto import Crypto
from Services.jwt import Jwt
from Services.email_service import EmailService

#Importing models
from Models.User import User


class Service:

	def email_text(self, username):
		html = """
		<html>
		  <body>
		    <p>Hi """+username+""",<br><br>
		       We are happy that you have joined us. From now on, we will try to keep you happy<br><br>
		       Thanks,<br>
			   Utopia Team,<br>
			   Bloomington,<br>
			   IN, US.<br>
		    </p>
		  </body>
		</html>
		"""
		return html

	#Service Method to register a user
	def register(self, user):
		try:
			service = Service()
			if(service.is_valid_email(user['email'])):
				if(user['password']==user['confirmPassword']):
					mongo_config = MongoConfig()
					collection = mongo_config.db()
					if(collection):
						crypto = Crypto()
						user['password'] = crypto.encrypted_string(user['password'])
						if "confirmPassword" in user:
							del user['confirmPassword']
						saved_user = collection.insert_one(user)
						email_queue = EmailService()
						email_queue.send_email(user['email'], 'Welcome to Utopia', self.email_text(user['firstName']))
						user = User()
						user.user_id = str(saved_user.inserted_id)
						return user
					else:
						return "Unable to connect"
				else:
					return "Password did not match"
			else:
				return "Please enter proper email address"
		except Exception as e:
			return e

	# Checkthe validity of the email
	def is_valid_email(self, email):
		return bool(re.search(r"^[\w\.\+\-]+\@[\w]+\.[a-z]{2,3}$", email))

	# Service method to generate a JWT and provide to user
	def login(self, data):
		try:
			mongo_config = MongoConfig()
			collection = mongo_config.db()
			if(collection):
				user_obj = User()
				search_user = {'email': data['email']}
				user = collection.find(search_user)
				crypto = Crypto()
				if(user):
					if(crypto.verify_decrypted_string(data['password'], user[0]['password'])):
						user_obj.first_name = user[0]['firstName']
						#user_obj.last_name = user[0]['lastName']
						user_obj.user_id = str(user[0]['_id'])
						user_obj.token = Jwt.encode_auth_token(user_id=user[0]['_id']).decode()
						return user_obj
					else:
						return "Invalid Credentials"
				else:
					return "User not available"
			else:
				return "Unable to connect to database"
		except IndexError as IE:
			return "User not available"
		except Exception as e:
			raise e
