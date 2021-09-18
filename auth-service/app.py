from website.app import create_app
from website.utils import gen_token

import os 
os.environ['AUTHLIB_INSECURE_TRANSPORT'] = '1'



app = create_app({
    'SECRET_KEY': 'secret',
    'OAUTH2_REFRESH_TOKEN_GENERATOR': True,
    'OAUTH2_ACCESS_TOKEN_GENERATOR': gen_token,
    'SQLALCHEMY_TRACK_MODIFICATIONS': False,
    'SQLALCHEMY_DATABASE_URI': 'sqlite:///db.sqlite',
})