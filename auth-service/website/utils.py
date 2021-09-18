import hashlib
from authlib.jose import jwt


def split_by_crlf(s):
    return [v for v in s.splitlines() if v]


def hash_password(password):
    salt = 'sssalt'
    pass_with_salt = password + salt
    return hashlib.md5(pass_with_salt.encode()).hexdigest()

def gen_token(client, grant_type, user, scope):
    header = {'alg': 'RS256'}
    payload = {
        'iss': 'http://127.0.0.1:5000',
        'sub': client.client_id,
        'aud': 'idk',
        'username': user.username,
        'id': user.id,
        'role': user.role
    }
    try:
        key = open('jwt-private.key', 'r').read()
        
        s = jwt.encode(header, payload, key)
        claims = jwt.decode(s, open('jwt-private.key.pub', 'r').read())
    except Exception as e:
        print('JWT exception', e)
    print("jwt encoded:{}\n decoded :{} \n header:{}".format(
        s, claims, claims.header))
    return s.decode()